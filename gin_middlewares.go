package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"gorm.io/gorm"
)

func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Set("db", db.WithContext(ctx))
		c.Next()
	}
}

func SatisfactionMiddleware(c *gin.Context) {
	var s feedbacktypes.Satisfaction

	if err := c.ShouldBindJSON(&s); err != nil {
		// If there's an error in parsing JSON, return an error response
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.Satisfied == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Field 'satisfied' not provided"})
		return
	}

	c.Set("satisfaction", s)

	c.Next()
}
