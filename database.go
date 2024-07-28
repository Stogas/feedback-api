package main

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/Stogas/feedback-api/internal/config"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB(conf config.DBConfig, tracing bool, issues []string) *gorm.DB {
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

		// apply migrations within a trace context
		dbMigrateWithTracing(db)

		// prefill issue types from config within a trace context
		fillDBWithIssueTypesTracing(db, issues)
	} else {
		// apply migrations without tracing
		err := dbMigrate(db)
		if err != nil {
			slog.Error("DB Migrations failed", "error", err)
			panic("DB migrations failed")
		}

		// prefill issue types from config without tracing
		err = fillDBWithIssueTypes(db, issues)
		if err != nil {
			slog.Error("DB issue type prefill failed", "error", err)
			panic("DB issue type prefill failed")
		}
	}

	return db
}

// the db...Tracing() functions are definitely not ideal and
// result in some duplicate work in them and in initDB(), but
// let's deal with that later

func dbMigrateWithTracing(db *gorm.DB) {
	ctx, span := otel.Tracer("GORM-auto-migrations").Start(context.Background(), "Run DB migrations")
	logger := slog.With("traceId", span.SpanContext().TraceID(), "spanId", span.SpanContext().SpanID())
	logger.Info("Running DB migrations ...")
	mErr := dbMigrate(db.WithContext(ctx))
	if mErr != nil {
		span.RecordError(mErr)
		logger.Error("DB Migrations failed", "error", mErr)
		span.End()
		panic("DB migrations failed")
	}
	logger.Info("DB Migrations succeeded!")
	span.End()
}

func dbMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&feedbacktypes.Satisfaction{},
	)
}

func fillDBWithIssueTypesTracing(db *gorm.DB, issues []string) {
	ctx, span := otel.Tracer("GORM-issue-loading").Start(context.Background(), "Fill DB with issue types")
	logger := slog.With("traceId", span.SpanContext().TraceID(), "spanId", span.SpanContext().SpanID())
	logger.Info("Filling DB with provided issue types ...")
	mErr := fillDBWithIssueTypes(db.WithContext(ctx), issues)
	if mErr != nil {
		span.RecordError(mErr)
		logger.Error("DB issue type prefill failed", "error", mErr)
		span.End()
		panic("DB issue type prefill failed")
	}
	logger.Info("DB issue type prefill succeeded!")
	span.End()
}

func fillDBWithIssueTypes(db *gorm.DB, typesFromConfig []string) error {
	// get existing issues in DB
	var existingIssues []feedbacktypes.Issue
	// existingTypesMap := make(map[string]feedbacktypes.Issue)
	err := db.Find(&existingIssues).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		slog.Error("Failed to fetch existing issue types", "error", err)
		return err
	}

	// delete issue types not present in config from DB
	for _, existingIssue := range existingIssues {
		if !slices.Contains(typesFromConfig, existingIssue.Name) {
			if err := db.Delete(&existingIssue).Error; err != nil {
				slog.Error("Failed to mark existing issue not present in config as deleted", "issueName", existingIssue.Name)
				return err
			} else {
				slog.Warn("Found issue type in DB, but not in config. Marked it as deleted", "issueName", existingIssue.Name)
			}
		}
	}

	// create any new issue types from config
	var existingTypesSlice []string
	for _, issue := range existingIssues {
		existingTypesSlice = append(existingTypesSlice, issue.Name)
	}
	for _, typeName := range typesFromConfig {
		if !slices.Contains(existingTypesSlice, typeName) {
			newType := feedbacktypes.Issue{Name: typeName}
			if err := db.Create(&newType).Error; err != nil {
				slog.Error("Failed to create issue type", "issueName", newType.Name, "error", err)
				return err
			} else {
				slog.Info("Created new issue type", "issueName", newType.Name)
			}
		}
	}

	return nil
}
