package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/ports"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("http-handler")

// Handler handles HTTP requests for items
type Handler struct {
	service ports.ItemService
	logger  *logger.Logger
}

// NewHandler creates a new HTTP handler
func NewHandler(service ports.ItemService, logger *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateItemRequest represents the request to create an item
type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error  string                    `json:"error"`
	Errors []validator.ValidationError `json:"errors,omitempty"`
}

// CreateItem handles POST /v1/items
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "CreateItem")
	defer span.End()

	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	v := validator.New()
	v.Required("name", req.Name).
		MinLength("name", req.Name, 1).
		MaxLength("name", req.Name, 100).
		Min("price", req.Price, 0)

	if !v.Valid() {
		h.respondValidationError(w, v.Errors())
		return
	}

	item := &domain.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := h.service.CreateItem(ctx, item); err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to create item", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	span.SetAttributes(attribute.String("item.id", item.ID))
	h.respondJSON(w, http.StatusCreated, item)
}

// GetItem handles GET /v1/items/{id}
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "GetItem")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]
	span.SetAttributes(attribute.String("item.id", id))

	item, err := h.service.GetItem(ctx, id)
	if err != nil {
		span.RecordError(err)
		if err == domain.ErrItemNotFound {
			h.respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		h.logger.Error("Failed to get item", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}

	h.respondJSON(w, http.StatusOK, item)
}

// ListItems handles GET /v1/items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "ListItems")
	defer span.End()

	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	span.SetAttributes(
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	items, err := h.service.ListItems(ctx, limit, offset)
	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list items", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"items":  items,
		"limit":  limit,
		"offset": offset,
		"count":  len(items),
	})
}

// UpdateItem handles PUT /v1/items/{id}
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "UpdateItem")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]
	span.SetAttributes(attribute.String("item.id", id))

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	v := validator.New()
	v.Required("name", req.Name).
		MinLength("name", req.Name, 1).
		MaxLength("name", req.Name, 100).
		Min("price", req.Price, 0)

	if !v.Valid() {
		h.respondValidationError(w, v.Errors())
		return
	}

	item := &domain.Item{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := h.service.UpdateItem(ctx, item); err != nil {
		span.RecordError(err)
		if err == domain.ErrItemNotFound {
			h.respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		h.logger.Error("Failed to update item", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	h.respondJSON(w, http.StatusOK, item)
}

// DeleteItem handles DELETE /v1/items/{id}
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "DeleteItem")
	defer span.End()

	vars := mux.Vars(r)
	id := vars["id"]
	span.SetAttributes(attribute.String("item.id", id))

	if err := h.service.DeleteItem(ctx, id); err != nil {
		span.RecordError(err)
		if err == domain.ErrItemNotFound {
			h.respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		h.logger.Error("Failed to delete item", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// respondJSON sends a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, ErrorResponse{Error: message})
}

// respondValidationError sends a validation error response
func (h *Handler) respondValidationError(w http.ResponseWriter, errors []validator.ValidationError) {
	h.respondJSON(w, http.StatusBadRequest, ErrorResponse{
		Error:  "Validation failed",
		Errors: errors,
	})
}
