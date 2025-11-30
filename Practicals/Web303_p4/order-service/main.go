package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hashicorp/go-uuid"
)

type Order struct {
	ID        string   `json:"id"`
	ItemIDs   []string `json:"item_ids"`
	Status    string   `json:"status"`
	Timestamp string   `json:"timestamp"`
}

var orders = make(map[string]Order)

// Service registration with Consul (simplified - no consul dependency)
func registerServiceWithConsul() {
	log.Println("Service registration skipped (no consul dependency)")
}

// Discover other services using Consul (simplified - returns hardcoded values)
func findService(serviceName string) (string, error) {
	// In a real Kubernetes environment, you'd use service names
	switch serviceName {
	case "food-catalog-service":
		return "http://food-catalog-service:8080", nil
	default:
		return "", fmt.Errorf("service %s not found", serviceName)
	}
}


func main() {
	// Try to register with Consul, but don't fail if it's not available
	go registerServiceWithConsul()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		var newOrder Order
		if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

        // Example of inter-service communication
        // Here you would call the food-catalog-service to validate ItemIDs
        catalogAddr, err := findService("food-catalog-service")
        if err != nil {
            http.Error(w, "Food catalog service not available", http.StatusInternalServerError)
            log.Printf("Error finding catalog service: %v", err)
            return
        }
        log.Printf("Found food-catalog-service at: %s. Would validate items here.", catalogAddr)


		orderID, _ := uuid.GenerateUUID()
		newOrder.ID = orderID
		newOrder.Status = "received"
		orders[orderID] = newOrder

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)
	})

	log.Println("Order Service starting on port 8081...")
	http.ListenAndServe(":8081", r)
}
