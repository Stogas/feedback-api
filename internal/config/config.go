package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
  Host string
  Port int
  User string
  Password string
  Name string
}

type Config struct {
  Database DBConfig
}

func New() *Config {
  return &Config{
    Database: DBConfig{
        Host: getEnvAsString("POSTGRES_HOST", "localhost"),
        Port: getEnvAsInt("POSTGRES_PORT", 5432),
        User: getEnvAsString("POSTGRES_USER", ""),
        Password: getEnvAsString("POSTGRES_PASSWORD", ""),
        Name: getEnvAsString("POSTGRES_DATABASE", ""),
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