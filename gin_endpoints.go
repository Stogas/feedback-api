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

	logger := getLogger(c.Request.Context())

	db := c.MustGet("db").(*gorm.DB)

	var existingSatisfaction feedbacktypes.Satisfaction
	existingRow := db.Where("uuid = ?", newSatisfaction.UUID).First(&existingSatisfaction)
	if existingRow.Error == nil {
		logger.Warn("A submission with this UUID already exists", "uuid", newSatisfaction.UUID, "method", c.Request.Method)
		c.JSON(http.StatusConflict, gin.H{"error": "A submission with this UUID already exists", "uuid": newSatisfaction.UUID, "created_at": existingSatisfaction.CreatedAt})
		return
	} else if existingRow.Error != gorm.ErrRecordNotFound {
		logger.Error("Error reading database", "error", existingRow.Error, "uuid", newSatisfaction.UUID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database read error"})
		return
	}

	result := db.Create(&newSatisfaction)

	if result.Error != nil {
		logger.Error("Database write error", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database write error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"created_at": newSatisfaction.CreatedAt})
}
