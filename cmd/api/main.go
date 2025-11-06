package main

import (
	"context"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/http"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/kafka"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/repository"
	"github.com/rodrigogrosa/Template-GoLang/internal/config"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/services"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/tracer"
	"github.com/rodrigogrosa/Template-GoLang/pkg/validation"

	_ "github.com/rodrigogrosa/Template-GoLang/docs" // Import for swagger docs
)

// @title Template GoLang Microservice API
// @version 1.0
// @description Production-ready Golang microservice template with hexagonal architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Log.Level)
	logger.Logger.Info().Msg("Starting application")

	// Initialize validator
	validation.Init()

	// Initialize tracer
	var shutdownTracer func(context.Context) error
	if cfg.Tracing.Enabled {
		var err error
		shutdownTracer, err = tracer.Init(cfg.Tracing.ServiceName, cfg.Tracing.JaegerEndpoint)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to initialize tracer")
		} else {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := shutdownTracer(ctx); err != nil {
					logger.Logger.Error().Err(err).Msg("Failed to shutdown tracer")
				}
			}()
		}
	}

	// Initialize repository
	itemRepo := repository.NewInMemoryItemRepository()

	// Initialize Kafka producer (optional)
	var messageProducer ports.MessageProducer
	if cfg.Kafka.Enabled {
		var err error
		messageProducer, err = kafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to initialize Kafka producer")
		} else {
			defer messageProducer.Close()
			logger.Logger.Info().Msg("Kafka producer initialized")
		}
	}

	// Initialize Kafka consumer (optional)
	var messageConsumer ports.MessageConsumer
	if cfg.Kafka.Enabled {
		var err error
		messageConsumer, err = kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)
		if err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to initialize Kafka consumer")
		} else {
			defer messageConsumer.Close()
			ctx := context.Background()
			if err := messageConsumer.Start(ctx); err != nil {
				logger.Logger.Error().Err(err).Msg("Failed to start Kafka consumer")
			} else {
				logger.Logger.Info().Msg("Kafka consumer started")
			}
		}
	}

	// Initialize service
	itemService := services.NewItemService(itemRepo, messageProducer)

	// Initialize HTTP server
	server := http.NewServer(cfg, itemService)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != stdhttp.ErrServerClosed {
			logger.Logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info().Msg("Shutdown signal received")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Logger.Info().Msg("Application stopped")
}
