package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"gorm.io/gorm"
)

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