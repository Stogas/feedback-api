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

func startAPI(conf config.APIConfig, globalMiddlewares []gin.HandlerFunc, dbMiddleware gin.HandlerFunc) {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	for _, m := range globalMiddlewares {
		// slog.Debug("Gin: Adding middleware")
		r.Use(m)
	}

	r.GET("/ping", ping)

	rSubmit := r.Group("/submit")
	rSubmit.Use(
		submitTokenMiddleware(conf.SubmitToken),
		dbMiddleware,
		satisfactionMiddleware,
	)
	{
		rSubmit.POST("/satisfaction", submitSatisfactionEndpoint)
	}

	slog.Info("Starting API", "host", conf.Host, "port", conf.Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", conf.Host, conf.Port),
		Handler: r.Handler(),
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("API listener failed", "error", err)
		}
	}()

	apiGracefulShutdown(srv)
}

func apiGracefulShutdown(srv *http.Server) {
	// Graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server with timeout
	timeout := 5 * time.Second
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down API listener ...", "timeout", timeout)

	// timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error while shutting down API listener gracefully. Initiating force shutdown...", "error", err, "timeout", timeout)
	} else {
		slog.Info("API listener exited successfully")
	}
}