package main

import (
	"net/http"

	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GetIssuesEndpoint(c *gin.Context) {
	logger := getLogger(c.Request.Context())

	db := c.MustGet("db").(*gorm.DB)

	var issues []feedbacktypes.Issue
	if err := db.Find(&issues).Error; err != nil {
		logger.Error("Failed to fetch issue types from DB")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database read error"})
	}
	c.JSON(http.StatusOK, issues)
}

func submitReportEndpoint(c *gin.Context) {
	newReport := c.MustGet("report").(feedbacktypes.Report)

	logger := getLogger(c.Request.Context())

	db := c.MustGet("db").(*gorm.DB)

	var existingReport feedbacktypes.Report
	existingRow := db.Where("uuid = ?", newReport.UUID).First(&existingReport)
	if existingRow.Error == nil {
		logger.Warn("A submission with this UUID already exists", "uuid", newReport.UUID, "method", c.Request.Method)
		c.JSON(http.StatusConflict, gin.H{"error": "A submission with this UUID already exists", "uuid": newReport.UUID, "created_at": existingReport.CreatedAt})
		return
	} else if existingRow.Error != gorm.ErrRecordNotFound {
		logger.Error("Error reading database", "error", existingRow.Error, "uuid", newReport.UUID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database read error"})
		return
	}

	result := db.Create(&newReport)

	if result.Error != nil {
		logger.Error("Database write error", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database write error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"created_at": newReport.CreatedAt})
}

func updateReportEndpoint(c *gin.Context) {
	newReport := c.MustGet("report").(feedbacktypes.Report)

	logger := getLogger(c.Request.Context())

	db := c.MustGet("db").(*gorm.DB)

	var existingReport feedbacktypes.Report
	existingRow := db.Where("uuid = ?", newReport.UUID).First(&existingReport)
	if existingRow.Error == gorm.ErrRecordNotFound {
		logger.Warn("A PATCH submission tried to modify a non-existing resource", "uuid", newReport.UUID, "method", c.Request.Method)
		c.JSON(http.StatusNotFound, gin.H{"error": "A submission with this UUID has not been found, submit via HTTP POST instead", "uuid": newReport.UUID})
		return
	} else if existingRow.Error != nil {
		logger.Error("Error reading database", "error", existingRow.Error, "uuid", newReport.UUID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database read error"})
		return
	}

	newReport.ID = existingReport.ID
	newReport.CreatedAt = existingReport.CreatedAt

	result := db.Save(&newReport)

	if result.Error != nil {
		logger.Error("Database write error", "error", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database write error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_at": newReport.CreatedAt, "updated_at": newReport.UpdatedAt, "uuid": newReport.UUID})
}
