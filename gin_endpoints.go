package main

import (
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
	newSatisfaction := c.MustGet("satisfaction").(feedbacktypes.Satisfaction)

	db := c.MustGet("db").(*gorm.DB)

	result := db.Create(&newSatisfaction)

	if result.Error != nil {
		logger := getLogger(c.Request.Context())
		logger.Error("Welp, got error writing into the database", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database write error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"created_at": newSatisfaction.CreatedAt})
}
