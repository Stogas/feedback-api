package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(conf config.DBConfig, tracing bool) *gorm.DB {
	// create connection config
	postgresConfig := postgres.New(postgres.Config{
		DSN: fmt.Sprintf( // data source name, refer https://github.com/jackc/pgx
			"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=UTC",
			conf.Host,
			conf.User,
			conf.Password,
			conf.Name,
			conf.Port,
		),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	})

	// set up logging
	gormLogger := slogGorm.New()

	// connect to database
	db, err := gorm.Open(postgresConfig, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		slog.Error("Failed to connect to database", "host", conf.Host, "port", conf.Port, "user", conf.User, "database", conf.Name)
		panic("failed to connect database")
	}

	// set up tracing
	if tracing {
		err := db.Use(otelgorm.NewPlugin())
		if err != nil {
			slog.Error("Failed to initialize GORM OTLP instrumentation", "error", err)
			panic("failed to initialize GORM OTLP instrumentation")
		}
	}

	// apply migrations within a trace context
	ctx, span := otel.Tracer("GORM-auto-migrations").Start(context.Background(), "Run DB migrations")
	logger := slog.With("traceId", span.SpanContext().TraceID(), "spanId", span.SpanContext().SpanID())
	logger.Info("Running DB migrations ...")
	mErr := db.WithContext(ctx).AutoMigrate(
		&feedbacktypes.Satisfaction{},
	)
	if mErr != nil {
		span.RecordError(mErr)
		logger.Error("DB Migrations failed", "error", mErr)
		span.End()
		panic("DB migrations failed")
	}
	logger.Info("DB Migrations succeeded!")
	span.End()

	return db
}
