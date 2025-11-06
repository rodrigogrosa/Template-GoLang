//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpAdapter "github.com/rodrigogrosa/Template-GoLang/internal/adapters/http"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/kafka"
	"github.com/rodrigogrosa/Template-GoLang/internal/adapters/repository"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/services"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/middleware"
)

func setupTestServer(t *testing.T) *httptest.Server {
	log := logger.New("info", "json")
	repo := repository.NewInMemoryRepository()
	publisher := kafka.NewNoOpPublisher()
	service := services.NewItemService(repo, publisher)
	handler := httpAdapter.NewHandler(service, log)
	
	authConfig := middleware.AuthConfig{
		Enabled: false,
	}
	
	server := httpAdapter.NewServer("", handler, authConfig)
	return httptest.NewServer(server.(*httpAdapter.Server))
}

func TestIntegration_CreateAndGetItem(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Create item
	createReq := map[string]interface{}{
		"name":        "Integration Test Item",
		"description": "Test Description",
		"price":       25.99,
	}
	
	body, _ := json.Marshal(createReq)
	resp, err := http.Post(ts.URL+"/v1/items", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
	
	var createdItem domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&createdItem); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	// Get item
	resp, err = http.Get(ts.URL + "/v1/items/" + createdItem.ID)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var retrievedItem domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&retrievedItem); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if retrievedItem.Name != createReq["name"] {
		t.Errorf("Expected name %s, got %s", createReq["name"], retrievedItem.Name)
	}
}

func TestIntegration_Health(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()
	
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to get health: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var health map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if health["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %s", health["status"])
	}
}

func TestIntegration_ListItems(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()
	
	// Create multiple items
	for i := 0; i < 3; i++ {
		createReq := map[string]interface{}{
			"name":  "Test Item",
			"price": 10.00,
		}
		body, _ := json.Marshal(createReq)
		http.Post(ts.URL+"/v1/items", "application/json", bytes.NewBuffer(body))
		time.Sleep(10 * time.Millisecond)
	}
	
	// List items
	resp, err := http.Get(ts.URL + "/v1/items?limit=10&offset=0")
	if err != nil {
		t.Fatalf("Failed to list items: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	items := result["items"].([]interface{})
	if len(items) < 3 {
		t.Errorf("Expected at least 3 items, got %d", len(items))
	}
}
