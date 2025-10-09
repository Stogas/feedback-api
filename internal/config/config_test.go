package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsString(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal string
		expected   string
		setEnv     bool
	}{
		{
			name:       "returns default when env var not set",
			envKey:     "TEST_STRING_VAR",
			defaultVal: "default_value",
			expected:   "default_value",
			setEnv:     false,
		},
		{
			name:       "returns env value when set",
			envKey:     "TEST_STRING_VAR",
			envValue:   "env_value",
			defaultVal: "default_value",
			expected:   "env_value",
			setEnv:     true,
		},
		{
			name:       "returns empty string when env var is empty",
			envKey:     "TEST_STRING_VAR",
			envValue:   "",
			defaultVal: "default_value",
			expected:   "",
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvAsString(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal int
		expected   int
		setEnv     bool
	}{
		{
			name:       "returns default when env var not set",
			envKey:     "TEST_INT_VAR",
			defaultVal: 42,
			expected:   42,
			setEnv:     false,
		},
		{
			name:       "returns parsed int when env var is valid integer",
			envKey:     "TEST_INT_VAR",
			envValue:   "123",
			defaultVal: 42,
			expected:   123,
			setEnv:     true,
		},
		{
			name:       "returns default when env var is invalid integer",
			envKey:     "TEST_INT_VAR",
			envValue:   "not_a_number",
			defaultVal: 42,
			expected:   42,
			setEnv:     true,
		},
		{
			name:       "returns default when env var is empty string",
			envKey:     "TEST_INT_VAR",
			envValue:   "",
			defaultVal: 42,
			expected:   42,
			setEnv:     true,
		},
		{
			name:       "handles negative numbers",
			envKey:     "TEST_INT_VAR",
			envValue:   "-456",
			defaultVal: 42,
			expected:   -456,
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvAsInt(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsBool_BasicCases(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal bool
		expected   bool
		setEnv     bool
	}{
		{
			name:       "returns default when env var not set",
			envKey:     "TEST_BOOL_VAR",
			defaultVal: true,
			expected:   true,
			setEnv:     false,
		},
		{
			name:       "returns true when env var is 'true'",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "true",
			defaultVal: false,
			expected:   true,
			setEnv:     true,
		},
		{
			name:       "returns false when env var is 'false'",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "false",
			defaultVal: true,
			expected:   false,
			setEnv:     true,
		},
		{
			name:       "returns default when env var is invalid boolean",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "invalid_bool",
			defaultVal: true,
			expected:   true,
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvAsBool(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsBool_SpecialCases(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal bool
		expected   bool
		setEnv     bool
	}{
		{
			name:       "returns default when env var is empty string",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "",
			defaultVal: false,
			expected:   false,
			setEnv:     true,
		},
		{
			name:       "handles '1' as true",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "1",
			defaultVal: false,
			expected:   true,
			setEnv:     true,
		},
		{
			name:       "handles '0' as false",
			envKey:     "TEST_BOOL_VAR",
			envValue:   "0",
			defaultVal: true,
			expected:   false,
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvAsBool(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsStringSlice(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		defaultVal []string
		expected   []string
		setEnv     bool
	}{
		{
			name:       "returns default when env var not set",
			envKey:     "TEST_SLICE_VAR",
			defaultVal: []string{"default1", "default2"},
			expected:   []string{"default1", "default2"},
			setEnv:     false,
		},
		{
			name:       "returns single item when no comma",
			envKey:     "TEST_SLICE_VAR",
			envValue:   "single_item",
			defaultVal: []string{"default"},
			expected:   []string{"single_item"},
			setEnv:     true,
		},
		{
			name:       "returns multiple items when comma-separated",
			envKey:     "TEST_SLICE_VAR",
			envValue:   "item1,item2,item3",
			defaultVal: []string{"default"},
			expected:   []string{"item1", "item2", "item3"},
			setEnv:     true,
		},
		{
			name:       "returns default when env var is empty",
			envKey:     "TEST_SLICE_VAR",
			envValue:   "",
			defaultVal: []string{"default1", "default2"},
			expected:   []string{"default1", "default2"},
			setEnv:     true,
		},
		{
			name:       "handles spaces around commas",
			envKey:     "TEST_SLICE_VAR",
			envValue:   "item1, item2 ,item3",
			defaultVal: []string{"default"},
			expected:   []string{"item1", " item2 ", "item3"},
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvAsStringSlice(tt.envKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsStringSliceRequired(t *testing.T) {
	tests := []struct {
		name        string
		envKey      string
		envValue    string
		expected    []string
		shouldPanic bool
		setEnv      bool
	}{
		{
			name:        "panics when env var not set",
			envKey:      "TEST_REQUIRED_SLICE_VAR",
			shouldPanic: true,
			setEnv:      false,
		},
		{
			name:        "panics when env var is empty",
			envKey:      "TEST_REQUIRED_SLICE_VAR",
			envValue:    "",
			shouldPanic: true,
			setEnv:      true,
		},
		{
			name:     "returns single item when no comma",
			envKey:   "TEST_REQUIRED_SLICE_VAR",
			envValue: "required_item",
			expected: []string{"required_item"},
			setEnv:   true,
		},
		{
			name:     "returns multiple items when comma-separated",
			envKey:   "TEST_REQUIRED_SLICE_VAR",
			envValue: "item1,item2,item3",
			expected: []string{"item1", "item2", "item3"},
			setEnv:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv(tt.envKey, tt.envValue)
			}

			if tt.shouldPanic {
				assert.Panics(t, func() {
					getEnvAsStringSliceRequired(tt.envKey)
				})
			} else {
				result := getEnvAsStringSliceRequired(tt.envKey)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("creates config with default values", func(t *testing.T) {
		// Set the required ISSUE_TYPES to avoid panic
		t.Setenv("ISSUE_TYPES", "bug,feature,improvement")

		config := New()

		// Test some default values
		assert.Equal(t, "0.0.0.0", config.API.Host)
		assert.Equal(t, 80, config.API.Port)
		assert.Equal(t, false, config.API.Debug)
		assert.Equal(t, "localhost", config.Database.Host)
		assert.Equal(t, 5432, config.Database.Port)
		assert.Equal(t, []string{"bug", "feature", "improvement"}, config.IssueTypes)
	})

	t.Run("creates config with custom environment values", func(t *testing.T) {
		// Set custom environment variables
		t.Setenv("ISSUE_TYPES", "custom1,custom2")
		t.Setenv("API_LISTEN_HOST", "192.168.1.1")
		t.Setenv("API_LISTEN_PORT", "8080")
		t.Setenv("API_DEBUG_MODE", "true")
		t.Setenv("POSTGRES_HOST", "db.example.com")
		t.Setenv("POSTGRES_PORT", "5433")
		t.Setenv("LOGS_JSON", "false")

		config := New()

		// Verify custom values are used
		assert.Equal(t, "192.168.1.1", config.API.Host)
		assert.Equal(t, 8080, config.API.Port)
		assert.Equal(t, true, config.API.Debug)
		assert.Equal(t, "db.example.com", config.Database.Host)
		assert.Equal(t, 5433, config.Database.Port)
		assert.Equal(t, false, config.Logs.JSON)
		assert.Equal(t, []string{"custom1", "custom2"}, config.IssueTypes)
	})
}
