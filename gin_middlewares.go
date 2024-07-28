package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Depado/ginprom"
	feedbacktypes "github.com/Stogas/feedback-api/internal/types"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func createDBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		c.Set("db", db.WithContext(ctx))
		c.Next()
	}
}

func metricsMiddleware(p *ginprom.Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the prometheus client to allow using custom metrics
		c.Set("prom", p)
		c.Next()
	}
}

func submitTokenMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Feedback-Submit-Token") == token {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-Feedback-Submit-Token not provided or incorrect"})
			return
		}
	}
}

func satisfactionMiddleware(c *gin.Context) {
	var s feedbacktypes.Satisfaction

	if err := c.ShouldBindJSON(&s); err != nil {
		// If there's an error in parsing JSON, return an error response
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// special handling for booleans, as it's necessary to detect if it was not provided (default value for booleans is False)
	if s.Satisfied == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Field 'satisfied' not provided"})
		return
	}

	c.Set("satisfaction", s)

	c.Next()

	if c.Request.Method == "POST" {
		logger := getLogger(c.Request.Context())
		statusCode := c.Writer.Status()
		if statusCode >= 200 && statusCode < 300 {
			logger.Debug("Response will be a success, will increment metrics")
			p := c.MustGet("prom").(*ginprom.Prometheus)
			err := p.IncrementCounterValue("satisfaction", []string{strconv.FormatBool(*s.Satisfied)})
			if err != nil {
				logger.Error("Failed to increment metrics counter")
			}
		} else {
			logger.Debug("Response will not be a success, skipping metrics increment")
		}
	}
}

func regularLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Create a logger without trace/span IDs
		logger := slog.Default()

		// Store the default logger in the request context
		ctx := context.WithValue(c.Request.Context(), contextLogger, logger)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Log request details
		logger.Info("Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
		)
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
		logger := slog.With("traceId", traceID, "spanId", spanID)

		// Store the logger in the request context
		ctx := context.WithValue(c.Request.Context(), contextLogger, logger)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Log request details
		logger.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
		)
	}
}
