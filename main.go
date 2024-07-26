package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	conf := config.New()

	initLogger(conf.Logs)
	gin.SetMode(gin.ReleaseMode)

	if conf.Tracing.Enabled {
		tracerClose := initTracer(conf.Tracing)
		defer tracerClose()
	}

	rMetrics, p := initMetrics()
	// start metrics listener in the background
	go startMetrics(rMetrics, conf.Metrics)
	// p.AddCustomCounter("satisfaction", "Counts how many good/bad satisfactions are received", []string{"satisfied"})

	db := initDB(conf.Database, conf.Tracing.Enabled)

	if conf.API.Debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if conf.Tracing.Enabled {
		r.Use(otelgin.Middleware("feedback-api"))
	}

	if conf.Logs.JSON {
		if conf.Tracing.Enabled {
			r.Use(traceLogMiddleware())
		} else {
			r.Use(regularLogMiddleware())
		}
	} else {
		r.Use(gin.Logger())
	}
	r.Use(p.Instrument())

	r.GET("/ping", ping)

	rSubmit := r.Group("/submit")
	rSubmit.Use(
		dbMiddleware(db),
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
