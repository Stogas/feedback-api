package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
			ctx := c.Request.Context()
			c.Set("db", db.WithContext(ctx))
			c.Next()
	}
}