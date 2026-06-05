package core

import (
	"os"
	"strconv"
)

type Config struct {
	// App
	AppName    string
	AppVersion string
	AppEnv     string // ← APP_ENV (missing)
	AppPort    string // ← APP_PORT (missing)
	AppURL     string

	// Database
	DBDriver string // ← DB_DRIVER (missing)
	DBDSN    string // ← DB_DSN   (missing — you hardcode DSN in config/db.go)

	// Logging
	LogLevel  string // ← LOG_LEVEL  (missing)
	LogFormat string // ← LOG_FORMAT (missing)

	// Security
	JWTSecret string // ← JWT_SECRET (missing)
	JWTExpiry string // ← JWT_EXPIRY (missing)

	// Modules
	ModuleAutoLoad bool // ← MODULE_AUTO_LOAD (missing)

	// CORS
	CORSAllowedOrigins string // ← CORS_ALLOWED_ORIGINS (you hardcode cfg.AppURL in middleware)
	CORSAllowedMethods string // ← CORS_ALLOWED_METHODS (hardcoded string in middleware)
	CORSAllowedHeaders string // ← CORS_ALLOWED_HEADERS (hardcoded string in middleware)

	// Dev
	Debug         bool // ← DEBUG (missing)
	AutoMigration bool // ← ENABLE_MIGRATION_AUTO (missing — you always auto-migrate)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return fallback
}

func NewConfig() *Config {
	return &Config{
		// App
		AppName:    getEnv("APP_NAME", "GoERP"),
		AppVersion: getEnv("APP_VERSION", "1.0"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8000"),
		AppURL:     getEnv("APP_URL", "http://localhost:8000"),

		// DB
		DBDriver: getEnv("DB_DRIVER", "sqlite"),
		DBDSN:    getEnv("DB_DSN", "file:goapp.db?_fk=1&_busy_timeout=5000"),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Security
		JWTSecret: getEnv("JWT_SECRET", "change_me"),
		JWTExpiry: getEnv("JWT_EXPIRY", "24h"),

		// Modules
		ModuleAutoLoad: getEnvBool("MODULE_AUTO_LOAD", true),

		// CORS
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

		// Dev
		Debug:         getEnvBool("DEBUG", true),
		AutoMigration: getEnv("ENABLE_MIGRATION_AUTO", "true") == "true",
	}
}
