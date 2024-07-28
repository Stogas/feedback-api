package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type APIConfig struct {
	Host        string
	Port        int
	Debug       bool
	SubmitToken string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type TraceConfig struct {
	Enabled bool
	Host    string
	Port    int
}

type LogsConfig struct {
	JSON   bool
	Debug  bool
	Source bool
}

type Config struct {
	API        APIConfig
	Database   DBConfig
	Tracing    TraceConfig
	Logs       LogsConfig
	Metrics    MetricsConfig
	IssueTypes []string
}

type MetricsConfig struct {
	Host string
	Port int
}

func New() *Config {
	return &Config{
		IssueTypes: getEnvAsStringSliceRequired("ISSUE_TYPES"),
		API: APIConfig{
			Host:        getEnvAsString("API_LISTEN_HOST", "0.0.0.0"),
			Port:        getEnvAsInt("API_LISTEN_PORT", 80),
			Debug:       getEnvAsBool("API_DEBUG_MODE", false),
			SubmitToken: getEnvAsString("API_SUBMIT_TOKEN", "test"),
		},
		Database: DBConfig{
			Host:     getEnvAsString("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnvAsString("POSTGRES_USER", ""),
			Password: getEnvAsString("POSTGRES_PASSWORD", ""),
			Name:     getEnvAsString("POSTGRES_DATABASE", ""),
		},
		Tracing: TraceConfig{
			Enabled: getEnvAsBool("OTLP_TRACING_ENABLED", false),
			Host:    getEnvAsString("OTLP_GRPC_HOST", "127.0.0.1"),
			Port:    getEnvAsInt("OTLP_GRPC_PORT", 4317),
		},
		Logs: LogsConfig{
			JSON:   getEnvAsBool("LOGS_JSON", true),
			Debug:  getEnvAsBool("LOGS_DEBUG", false),
			Source: getEnvAsBool("LOGS_SOURCE", false),
		},
		Metrics: MetricsConfig{
			Host: getEnvAsString("METRICS_HOST", "0.0.0.0"),
			Port: getEnvAsInt("METRICS_PORT", 2222),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnvAsString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnvAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnvAsString(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read a comma-separated environment variable into a slice of strings
func getEnvAsStringSliceRequired(name string) []string {
	valStr := getEnvAsString(name, "")
	if valStr == "" {
		slog.Error("Missing required environment variable", "envVar", name)
		panic("Missing required environment variable")
	}
	return strings.Split(valStr, ",")
}
