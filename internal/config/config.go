package config

import (
	"os"
	"strconv"
	"strings"
)

type SqlConfig struct {
	Driver   string
	Host     string
	User     string
	Password string
	DbName   string
	SslMode  string
}

type TelegramBotConfig struct {
	Token string
}

type Config struct {
	Sql         SqlConfig
	TelegramBot TelegramBotConfig
}

// New returns a new Config struct
func New() Config {
	return Config{
		Sql: SqlConfig{
			Driver:   getEnv("SQL_DRIVER", ""),
			Host:     getEnv("SQL_HOST", ""),
			User:     getEnv("SQL_USER", ""),
			Password: getEnv("SQL_PASSWORD", ""),
			DbName:   getEnv("SQL_DB_NAME", ""),
			SslMode:  getEnv("SQL_SSL_MODE", ""),
		},
		TelegramBot: TelegramBotConfig{
			Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
