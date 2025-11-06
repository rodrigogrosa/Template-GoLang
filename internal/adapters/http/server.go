package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/rodrigogrosa/Template-GoLang/internal/config"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/middleware"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	config *config.ServerConfig
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, itemService ports.ItemService) *Server {
	router := mux.NewRouter()

	// Middlewares
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.TracingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Apply JWT middleware if configured
	if cfg.Security.JWTSecret != "" {
		router.Use(middleware.JWTMiddleware(cfg.Security.JWTSecret))
	}

	// Health endpoint
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API v1 routes
	itemHandler := NewItemHandler(itemService)
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	v1.HandleFunc("/items", itemHandler.ListItems).Methods("GET")
	v1.HandleFunc("/items/{id}", itemHandler.GetItem).Methods("GET")
	v1.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	v1.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return &Server{
		server: srv,
		config: &cfg.Server,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	logger.Logger.Info().Str("addr", s.server.Addr).Msg("Starting HTTP server")

	if s.config.TLSEnabled && s.config.TLSCertFile != "" && s.config.TLSKeyFile != "" {
		logger.Logger.Info().Msg("TLS enabled")
		return s.server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Logger.Info().Msg("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// healthHandler handles health check requests
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}
