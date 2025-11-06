package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/metrics"
	"github.com/rodrigogrosa/Template-GoLang/pkg/validation"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	service ports.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler(service ports.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// CreateItemRequest represents the request to create an item
type CreateItemRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item with the provided details
// @Tags items
// @Accept json
// @Produce json
// @Param item body CreateItemRequest true "Item details"
// @Success 201 {object} domain.Item
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/items [post]
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	item := &domain.Item{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.CreateItem(r.Context(), item); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to create item")
		respondError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	metrics.ItemsCreated.Inc()
	respondJSON(w, http.StatusCreated, item)
}

// GetItem godoc
// @Summary Get an item by ID
// @Description Get an item by its ID
// @Tags items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} domain.Item
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/items/{id} [get]
func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	item, err := h.service.GetItem(r.Context(), id)
	if err != nil {
		if err == domain.ErrItemNotFound {
			respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		logger.Logger.Error().Err(err).Msg("Failed to get item")
		respondError(w, http.StatusInternalServerError, "Failed to get item")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

// ListItems godoc
// @Summary List all items
// @Description Get a list of all items with pagination
// @Tags items
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.Item
// @Failure 500 {object} ErrorResponse
// @Router /v1/items [get]
func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	items, err := h.service.ListItems(r.Context(), limit, offset)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to list items")
		respondError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

// UpdateItem godoc
// @Summary Update an item
// @Description Update an item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body UpdateItemRequest true "Item details"
// @Success 200 {object} domain.Item
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/items/{id} [put]
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	item := &domain.Item{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.UpdateItem(r.Context(), item); err != nil {
		if err == domain.ErrItemNotFound {
			respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		logger.Logger.Error().Err(err).Msg("Failed to update item")
		respondError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	metrics.ItemsUpdated.Inc()
	respondJSON(w, http.StatusOK, item)
}

// DeleteItem godoc
// @Summary Delete an item
// @Description Delete an item by its ID
// @Tags items
// @Param id path string true "Item ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/items/{id} [delete]
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteItem(r.Context(), id); err != nil {
		if err == domain.ErrItemNotFound {
			respondError(w, http.StatusNotFound, "Item not found")
			return
		}
		logger.Logger.Error().Err(err).Msg("Failed to delete item")
		respondError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	metrics.ItemsDeleted.Inc()
	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
