package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Helper to create test context with logger
func createTestContextWithLogger() (*gin.Context, *httptest.ResponseRecorder) {
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

func TestSubmitTokenMiddleware(t *testing.T) {
	t.Run("allows request when token is empty (no token required)", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/submit", nil)

		middleware := submitTokenMiddleware("")
		middleware(c)

		// If not aborted, middleware passed
		assert.False(t, c.IsAborted())
		assert.Equal(t, http.StatusOK, w.Code) // Default status
	})

	t.Run("allows request when correct token is provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/submit", nil)
		c.Request.Header.Set("X-Feedback-Submit-Token", "valid-token")

		middleware := submitTokenMiddleware("valid-token")
		middleware(c)

		assert.False(t, c.IsAborted())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("blocks request when token is required but not provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/submit", nil)
		// No token header set

		middleware := submitTokenMiddleware("required-token")
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "X-Feedback-Submit-Token not provided or incorrect", response["error"])
	})

	t.Run("blocks request when incorrect token is provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/submit", nil)
		c.Request.Header.Set("X-Feedback-Submit-Token", "wrong-token")

		middleware := submitTokenMiddleware("correct-token")
		middleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "X-Feedback-Submit-Token not provided or incorrect", response["error"])
	})
}

func TestReportMiddleware(t *testing.T) {
	t.Run("rejects request with invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create request first
		invalidJSON := `{"uuid": "invalid-uuid", "satisfied": not-a-boolean}`
		c.Request = httptest.NewRequest("POST", "/reports", strings.NewReader(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		// Set up logger
		logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
		ctx := context.WithValue(c.Request.Context(), contextLogger, logger)
		c.Request = c.Request.WithContext(ctx)

		reportMiddleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Just check that there's an error, the exact message may vary
		assert.NotEmpty(t, response["error"])
	})

	t.Run("rejects request when satisfied field is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		testUUID := uuid.New()
		jsonBody := `{
			"uuid": "` + testUUID.String() + `",
			"comment": "Missing satisfied field"
		}`

		// Create request first
		c.Request = httptest.NewRequest("POST", "/reports", strings.NewReader(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// Set up logger
		logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
		ctx := context.WithValue(c.Request.Context(), contextLogger, logger)
		c.Request = c.Request.WithContext(ctx)

		reportMiddleware(c)

		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		// The validation error message may vary, just check there's an error about satisfied field
		assert.Contains(t, response["error"], "Satisfied")
	})
}

func TestCreateDBMiddleware(t *testing.T) {
	t.Run("creates middleware function correctly", func(t *testing.T) {
		// We can't easily test GORM DB operations in unit tests without complex setup
		// But we can test that the function creates a valid middleware
		mockDB := &gorm.DB{}
		middleware := createDBMiddleware(mockDB)

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(func(*gin.Context) {}), middleware)
	})
}

func TestMetricsMiddleware(t *testing.T) {
	t.Run("sets prometheus client in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		mockProm := &ginprom.Prometheus{}
		middleware := metricsMiddleware(mockProm)
		middleware(c)

		assert.False(t, c.IsAborted())

		// Verify prometheus client is set in context
		prom, exists := c.Get("prom")
		assert.True(t, exists)
		assert.Equal(t, mockProm, prom)
	})
}

func TestRegularLogMiddleware(t *testing.T) {
	t.Run("creates middleware function correctly", func(t *testing.T) {
		middleware := regularLogMiddleware()

		assert.NotNil(t, middleware)
		assert.IsType(t, gin.HandlerFunc(func(*gin.Context) {}), middleware)
	})

	t.Run("sets logger in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		middleware := regularLogMiddleware()
		middleware(c)

		// Verify logger is accessible through getLogger
		logger := getLogger(c.Request.Context())
		assert.NotNil(t, logger)
	})
}
