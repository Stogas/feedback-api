package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Stogas/feedback-api/internal/config"
)

func startAPI(conf *config.Config, globalMiddlewares []gin.HandlerFunc, dbMiddleware gin.HandlerFunc) {

	r := gin.New()

	r.Use(gin.Recovery())
	for _, m := range globalMiddlewares {
		// slog.Debug("Gin: Adding middleware")
		r.Use(m)
	}

	r.GET("/ping", ping)

	rSubmit := r.Group("/submit")
	rSubmit.Use(
		dbMiddleware,
		satisfactionMiddleware,
	)
	{
		rSubmit.POST("/satisfaction", submitSatisfactionEndpoint)
	}

	slog.Info("Starting API", "host", conf.API.Host, "port", conf.API.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", conf.API.Host, conf.API.Port),
		Handler: r.Handler(),
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("API listener failed", "error", err)
		}
	}()

	// Graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with timeout
	timeout := 5 * time.Second
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down API listener ...")

	// timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error while shutting down API listener gracefully. Initiating force shutdown...", "error", err, "timeout", timeout)
	} else {
		slog.Info("API listener exited successfully")
	}
}
