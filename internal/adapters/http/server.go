package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rodrigogrosa/Template-GoLang/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	router *mux.Router
}

// NewServer creates a new HTTP server
func NewServer(addr string, handler *Handler, authConfig middleware.AuthConfig) *Server {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.Recovery)
	router.Use(middleware.Logging)
	router.Use(middleware.Metrics)
	router.Use(middleware.CORS)

	// Health endpoint (no auth required)
	router.HandleFunc("/health", handler.Health).Methods("GET")

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes with optional JWT auth
	api := router.PathPrefix("/v1").Subrouter()
	api.Use(middleware.JWTAuth(authConfig))

	api.HandleFunc("/items", handler.CreateItem).Methods("POST")
	api.HandleFunc("/items", handler.ListItems).Methods("GET")
	api.HandleFunc("/items/{id}", handler.GetItem).Methods("GET")
	api.HandleFunc("/items/{id}", handler.UpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", handler.DeleteItem).Methods("DELETE")

	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		router: router,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// StartTLS starts the HTTP server with TLS
func (s *Server) StartTLS(certFile, keyFile string) error {
	return s.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
