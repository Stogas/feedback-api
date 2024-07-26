package main

import (
	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"

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
	gin.SetMode(gin.ReleaseMode)

	// logging
	initLogger(conf.Logs)
	var loggingMiddleware gin.HandlerFunc
	if conf.Logs.JSON {
		if conf.Tracing.Enabled {
			loggingMiddleware = traceLogMiddleware()
		} else {
			loggingMiddleware = regularLogMiddleware()
		}
	} else {
		loggingMiddleware = gin.Logger()
	}

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

	startAPI(conf, loggingMiddleware, db, p)
}
