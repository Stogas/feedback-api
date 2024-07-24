package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	host := "localhost"
	user := "test"
	password := "test"
	dbname := "test"
	port := "5432"
	sslmode := "disable"
	timeZone := "UTC"

	postgresConfig := postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timeZone), // data source name, refer https://github.com/jackc/pgx
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

func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx := c.Request.Context()
			c.Set("db", db.WithContext(ctx))
			c.Next()
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func submitSatisfactionEndpoint(c *gin.Context) {
	var newSatisfaction feedbacktypes.Satisfaction

	if err := c.ShouldBindJSON(&newSatisfaction); err != nil {
    // If there's an error in parsing JSON, return an error response
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

	db := c.MustGet("db").(*gorm.DB)

	result := db.Create(&newSatisfaction)

	if result.Error != nil {
		slog.Error("Welp, got error writing into the database")
	}
}