package main

import (
	"fmt"
	"strconv"

	"log/slog"

	"github.com/Stogas/feedback-api/internal/config"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"

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

	postgresConfig := postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.Name, strconv.Itoa(conf.Database.Port)), // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	})

	db, err := gorm.Open(postgresConfig, &gorm.Config{})
	if err != nil {
    panic("failed to connect database")
  }
	db.AutoMigrate(&feedbacktypes.Satisfaction{})

	r := gin.Default()
	r.GET("/ping", ping)

	rSubmit := r.Group("/submit")
	rSubmit.Use(DBMiddleware(db))
	{
		rSubmit.POST("/satisfaction", submitSatisfactionEndpoint)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}