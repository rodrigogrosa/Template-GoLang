package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Kafka      KafkaConfig
	Logging    LoggingConfig
	Telemetry  TelemetryConfig
	Auth       AuthConfig
	TLS        TLSConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type string // in-memory, postgres, etc.
	DSN  string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Enabled bool
	Brokers []string
	Topic   string
	GroupID string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string // json or console
}

// TelemetryConfig holds observability configuration
type TelemetryConfig struct {
	MetricsEnabled bool
	MetricsPort    int
	TracingEnabled bool
	TracingURL     string
	ServiceName    string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret string
	JWTIssuer string
	Enabled   bool
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:     getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Type: getEnv("DB_TYPE", "in-memory"),
			DSN:  getEnv("DB_DSN", ""),
		},
		Kafka: KafkaConfig{
			Enabled: getEnvAsBool("KAFKA_ENABLED", false),
			Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnv("KAFKA_TOPIC", "item-events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "item-service"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Telemetry: TelemetryConfig{
			MetricsEnabled: getEnvAsBool("METRICS_ENABLED", true),
			MetricsPort:    getEnvAsInt("METRICS_PORT", 9090),
			TracingEnabled: getEnvAsBool("TRACING_ENABLED", true),
			TracingURL:     getEnv("TRACING_URL", "http://localhost:4318"),
			ServiceName:    getEnv("SERVICE_NAME", "golang-microservice"),
		},
		Auth: AuthConfig{
			Enabled:   getEnvAsBool("AUTH_ENABLED", false),
			JWTSecret: getEnv("JWT_SECRET", ""),
			JWTIssuer: getEnv("JWT_ISSUER", "golang-microservice"),
		},
		TLS: TLSConfig{
			Enabled:  getEnvAsBool("TLS_ENABLED", false),
			CertFile: getEnv("TLS_CERT_FILE", ""),
			KeyFile:  getEnv("TLS_KEY_FILE", ""),
		},
	}

	// Validate required fields
	if cfg.Auth.Enabled && cfg.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required when AUTH_ENABLED is true")
	}
	if cfg.TLS.Enabled && (cfg.TLS.CertFile == "" || cfg.TLS.KeyFile == "") {
		return nil, fmt.Errorf("TLS_CERT_FILE and TLS_KEY_FILE are required when TLS_ENABLED is true")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		return []string{value}
	}
	return defaultValue
}
