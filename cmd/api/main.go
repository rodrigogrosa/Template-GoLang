package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	httpAdapter "github.com/rodrigogrosa/Template-GoLang/internal/adapters/http"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/kafka"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/repository"
	"github.com/rodrigogrosa/Template-GoLang/internal/config"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/ports"
	"github.com/rodrigogrosa/Template-GoLang/internal/services"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Starting application", "config", cfg)

	// Initialize tracing
	if cfg.Telemetry.TracingEnabled {
		shutdown, err := initTracing(cfg.Telemetry.ServiceName, cfg.Telemetry.TracingURL)
		if err != nil {
			log.Error("Failed to initialize tracing", err)
		} else {
			defer func() {
				if err := shutdown(context.Background()); err != nil {
					log.Error("Failed to shutdown tracing", err)
				}
			}()
			log.Info("Tracing initialized")
		}
	}

	// Initialize repository
	repo := repository.NewInMemoryRepository()
	log.Info("Repository initialized")

	// Initialize event publisher
	var publisher ports.EventPublisher
	if cfg.Kafka.Enabled {
		pub, err := kafka.NewPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic, log)
		if err != nil {
			log.Error("Failed to create Kafka publisher", err)
			log.Info("Using no-op publisher")
			publisher = kafka.NewNoOpPublisher()
		} else {
			publisher = pub
			defer func() {
				if err := publisher.Close(); err != nil {
					log.Error("Failed to close publisher", err)
				}
			}()
			log.Info("Kafka publisher initialized")
		}
	} else {
		publisher = kafka.NewNoOpPublisher()
		log.Info("Kafka disabled, using no-op publisher")
	}

	// Initialize service
	service := services.NewItemService(repo, publisher)
	log.Info("Service initialized")

	// Initialize HTTP handler
	handler := httpAdapter.NewHandler(service, log)

	// Create auth config
	authConfig := middleware.AuthConfig{
		Secret:  cfg.Auth.JWTSecret,
		Issuer:  cfg.Auth.JWTIssuer,
		Enabled: cfg.Auth.Enabled,
	}

	// Initialize HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := httpAdapter.NewServer(addr, handler, authConfig)
	log.Info("HTTP server initialized", "address", addr)

	// Start metrics server if enabled
	if cfg.Telemetry.MetricsEnabled {
		metricsAddr := fmt.Sprintf(":%d", cfg.Telemetry.MetricsPort)
		go func() {
			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.Handler())
			log.Info("Metrics server started", "address", metricsAddr)

			metricsServer := &http.Server{
				Addr:              metricsAddr,
				Handler:           mux,
				ReadTimeout:       5 * time.Second,
				WriteTimeout:      10 * time.Second,
				ReadHeaderTimeout: 5 * time.Second,
			}

			if err := metricsServer.ListenAndServe(); err != nil {
				log.Error("Metrics server failed", err)
			}
		}()
	}

	// Start Kafka consumer if enabled
	if cfg.Kafka.Enabled {
		consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, log)
		if err != nil {
			log.Error("Failed to create Kafka consumer", err)
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				err := consumer.Subscribe(ctx, func(event *domain.ItemEvent) error {
					log.Info("Received event", "event_type", event.EventType, "item_id", event.ItemID)
					return nil
				})
				if err != nil {
					log.Error("Consumer error", err)
				}
			}()
			defer func() {
				if err := consumer.Close(); err != nil {
					log.Error("Failed to close consumer", err)
				}
			}()
			log.Info("Kafka consumer started")
		}
	}

	// Start HTTP server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Starting HTTP server", "address", addr, "tls_enabled", cfg.TLS.Enabled)
		if cfg.TLS.Enabled {
			serverErrors <- server.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			serverErrors <- server.Start()
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal("Server error", err)
	case sig := <-quit:
		log.Info("Received shutdown signal", "signal", sig.String())

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Error("Server shutdown error", err)
			os.Exit(1)
		}

		log.Info("Server shutdown complete")
	}
}

// initTracing initializes OpenTelemetry tracing
func initTracing(serviceName, tracingURL string) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(tracingURL),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	resource, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
