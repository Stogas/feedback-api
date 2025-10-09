package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test gin context with basic setup
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a basic request first
	c.Request = httptest.NewRequest("GET", "/", nil)

	// Set up a logger in the context
	logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	ctx := context.WithValue(c.Request.Context(), contextLogger, logger)
	c.Request = c.Request.WithContext(ctx)

	return c, w
}

func TestPing(t *testing.T) {
	t.Run("returns pong message with 200 status", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("GET", "/ping", nil)

		ping(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "pong", response["message"])
	})

	t.Run("handles different request methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			c, w := createTestContext()
			c.Request = httptest.NewRequest(method, "/ping", nil)

			ping(c)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "pong", response["message"])
		}
	})

	t.Run("returns proper JSON content type", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("GET", "/ping", nil)

		ping(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}

// Test that endpoints handle missing database context gracefully
func TestEndpointErrorHandling(t *testing.T) {
	t.Run("GetIssuesEndpoint panics gracefully when no database context", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("GET", "/issues", nil)
		// Note: Not setting database context to test error handling

		// This should panic with MustGet, which is expected behavior
		// In a real app, middleware would ensure db is always present
		assert.Panics(t, func() {
			GetIssuesEndpoint(c)
		})

		// Response recorder should not have been written to
		assert.Equal(t, http.StatusOK, w.Code) // Default value, nothing written
	})

	t.Run("submitReportEndpoint panics gracefully when no report context", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("POST", "/reports", nil)
		// Note: Not setting report context to test error handling

		// This should panic with MustGet, which is expected behavior
		assert.Panics(t, func() {
			submitReportEndpoint(c)
		})

		// Response recorder should not have been written to
		assert.Equal(t, http.StatusOK, w.Code) // Default value, nothing written
	})

	t.Run("updateReportEndpoint panics gracefully when no report context", func(t *testing.T) {
		c, w := createTestContext()
		c.Request = httptest.NewRequest("PATCH", "/reports", nil)
		// Note: Not setting report context to test error handling

		// This should panic with MustGet, which is expected behavior
		assert.Panics(t, func() {
			updateReportEndpoint(c)
		})

		// Response recorder should not have been written to
		assert.Equal(t, http.StatusOK, w.Code) // Default value, nothing written
	})
}

// Test helper function behavior
func TestCreateTestContext(t *testing.T) {
	t.Run("creates valid gin context with logger", func(t *testing.T) {
		c, w := createTestContext()

		assert.NotNil(t, c)
		assert.NotNil(t, w)
		assert.NotNil(t, c.Request)
		assert.Equal(t, gin.TestMode, gin.Mode())

		// Verify logger is in context
		logger := getLogger(c.Request.Context())
		assert.NotNil(t, logger)
	})
}
