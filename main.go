package main

import (
	"context"
	"fmt"
	"strconv"

	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
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

	postgresConfig := postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, strconv.Itoa(conf.Database.Port)), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	})

	db, err := gorm.Open(postgresConfig, &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "host", conf.Database.Host, "port", conf.Database.Port, "user", conf.Database.User, "database", conf.Database.Name)
    panic("failed to connect database")
  }
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		slog.Error("Failed to initialize GORM OTLP instrumentation", "error", err)
		panic("failed to initialize GORM OTLP instrumentation")
	}
	db.AutoMigrate(&feedbacktypes.Satisfaction{})

	r := gin.Default()

	r.Use(otelgin.Middleware("feedback-api"))

	r.GET("/ping", ping)

	rSubmit := r.Group("/submit")
	rSubmit.Use(DBMiddleware(db))
	{
		rSubmit.POST("/satisfaction", submitSatisfactionEndpoint)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}