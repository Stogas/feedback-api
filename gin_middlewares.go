package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func dbMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Set("db", db.WithContext(ctx))
		c.Next()
	}
}

func satisfactionMiddleware(c *gin.Context) {
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

func regularLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the default logger in the request context
		ctx := context.WithValue(c.Request.Context(), "logger", slog.Default())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func traceLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extract the span from the context
		span := trace.SpanFromContext(c.Request.Context())

		// Get the trace and span IDs
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Create a logger with trace/span IDs
		logger := slog.With("traceID", traceID, "spanID", spanID)

		// Store the logger in the request context
		ctx := context.WithValue(c.Request.Context(), "logger", logger)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Log request details
		logger.Info("Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"traceID", traceID,
			"spanID", spanID,
		)
	}
}
