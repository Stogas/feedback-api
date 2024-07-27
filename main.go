package main

import (
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

	gin.SetMode(gin.ReleaseMode)
	var globalMiddlewares []gin.HandlerFunc

	// tracing
	if conf.Tracing.Enabled {
		tracerClose, b3 := initTracer(conf.Tracing)
		defer tracerClose()
		globalMiddlewares = append(globalMiddlewares, otelgin.Middleware("feedback-api", otelgin.WithPropagators(b3)))
	}

	// logging
	initLogger(conf.Logs)
	var l gin.HandlerFunc
	if conf.Logs.JSON {
		if conf.Tracing.Enabled {
			l = traceLogMiddleware()
		} else {
			l = regularLogMiddleware()
		}
	} else {
		l = gin.Logger()
	}
	globalMiddlewares = append(globalMiddlewares, l)

	// metrics
	rMetrics, p := initMetrics(globalMiddlewares)
	// start metrics listener in the background
	go startMetrics(rMetrics, conf.Metrics)
	globalMiddlewares = append(globalMiddlewares, p.Instrument(), metricsMiddleware(p))

	// database
	db := initDB(conf.Database, conf.Tracing.Enabled)
	dbMiddleware := createDBMiddleware(db)

	startAPI(conf.API, globalMiddlewares, dbMiddleware)
}
