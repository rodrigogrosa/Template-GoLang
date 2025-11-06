package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Log      LogConfig
	Tracing  TracingConfig
	Kafka    KafkaConfig
	Security SecurityConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
	TLSEnabled      bool
	TLSCertFile     string
	TLSKeyFile      string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled        bool
	JaegerEndpoint string
	ServiceName    string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Enabled bool
	Brokers []string
	Topic   string
	GroupID string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	JWTSecret string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Host:            getEnv("HOST", "0.0.0.0"),
			ReadTimeout:     getEnvAsInt("READ_TIMEOUT", 10),
			WriteTimeout:    getEnvAsInt("WRITE_TIMEOUT", 10),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 30),
			TLSEnabled:      getEnvAsBool("TLS_ENABLED", false),
			TLSCertFile:     getEnv("TLS_CERT_FILE", ""),
			TLSKeyFile:      getEnv("TLS_KEY_FILE", ""),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Tracing: TracingConfig{
			Enabled:        getEnvAsBool("TRACING_ENABLED", false),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", ""),
			ServiceName:    getEnv("SERVICE_NAME", "template-service"),
		},
		Kafka: KafkaConfig{
			Enabled: getEnvAsBool("KAFKA_ENABLED", false),
			Brokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnv("KAFKA_TOPIC", "items"),
			GroupID: getEnv("KAFKA_GROUP_ID", "template-service"),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", ""),
		},
	}
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

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Parse comma-separated values
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
