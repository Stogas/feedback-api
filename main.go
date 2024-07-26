package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"log/slog"

	"github.com/Depado/ginprom"
	"github.com/Stogas/feedback-api/internal/config"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
	}
}

func main() {
	conf := config.New()

	initLogger(conf.API.JSONlogging)

	gin.SetMode(gin.ReleaseMode)

	if conf.Tracing.Enabled {
		// Initialize tracer
		tp, err := initTracer(conf.Tracing)
		if err != nil {
			slog.Error("failed to initialize tracer", "error", err)
			panic("failed to initialize tracer")
		}
		// Clean up on shutdown
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				slog.Error("failed to shut down tracer", "error", err)
				panic("failed to shut down tracer")
			}
		}()
	}

	// metrics
	rMetrics := gin.Default()
	p := ginprom.New(
		ginprom.Engine(rMetrics),
		ginprom.Subsystem("feedbackapi"),
		ginprom.Path("/metrics"),
	)
	slog.Info("Starting Prometheus exporter", "host", conf.Metrics.Host, "port", conf.Metrics.Port)
	go func() {
		if err := rMetrics.Run(fmt.Sprintf("%s:%v", conf.Metrics.Host, conf.Metrics.Port)); err != nil {
			slog.Error("Failed to run metrics exporter", "error", err)
		}
	}()
	// p.AddCustomCounter("satisfaction", "Counts how many good/bad satisfactions are received", []string{"satisfied"})

	// Database
	postgresConfig := postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, strconv.Itoa(conf.Database.Port)), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                                                                                                                                                                                            // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	})

	gormLogger := slogGorm.New()
	db, err := gorm.Open(postgresConfig, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "host", conf.Database.Host, "port", conf.Database.Port, "user", conf.Database.User, "database", conf.Database.Name)
		panic("failed to connect database")
	}
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		slog.Error("Failed to initialize GORM OTLP instrumentation", "error", err)
		panic("failed to initialize GORM OTLP instrumentation")
	}
	db.AutoMigrate(&feedbacktypes.Satisfaction{})

	if conf.API.Debug {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if conf.Tracing.Enabled {
		r.Use(otelgin.Middleware("feedback-api"))
	}

	if conf.API.JSONlogging {
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
		slog.Info("API listener exiting")
	}
}
