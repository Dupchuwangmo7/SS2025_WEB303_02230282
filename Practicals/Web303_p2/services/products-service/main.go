// services/products-service/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	consulapi "github.com/hashicorp/consul/api"
)

const serviceName = "products-service"
const servicePort = 8082

func main() {
	if err := registerWithConsul(); err != nil {
		log.Fatalf("Service registration failed: %v", err)
	}

	mux := chi.NewRouter()
	mux.Get("/health", handleHealthStatus)
	mux.Get("/products/{id}", handleProductRequest)

	log.Printf("%s is starting on port %d", serviceName, servicePort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", servicePort), mux); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}

// handleProductRequest returns product information.
func handleProductRequest(w http.ResponseWriter, r *http.Request) {
	prodID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Product Info - Service: %s, Product ID: %s\n", serviceName, prodID)
}

// handleHealthStatus indicates the service is operational.
func handleHealthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Healthy")
}

func registerWithConsul() error {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to get hostname: %w", err)
	}

	serviceReg := &consulapi.AgentServiceRegistration{
		ID:      serviceName + "-" + hostname,
		Name:    serviceName,
		Port:    servicePort,
		Address: hostname,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, servicePort),
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	if err := client.Agent().ServiceRegister(serviceReg); err != nil {
		return fmt.Errorf("service registration error: %w", err)
	}

	log.Printf("Service %s registered successfully", serviceName)
	return nil
}