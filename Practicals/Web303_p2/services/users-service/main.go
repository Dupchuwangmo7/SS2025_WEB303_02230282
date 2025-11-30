// services/users-service/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	consulapi "github.com/hashicorp/consul/api"
)

const serviceName = "users-service"
const servicePort = 8081

func main() {
	// Register with service discovery
	if err := registerWithConsul(); err != nil {
		log.Fatalf("Registration error: %v", err)
	}

	// Setup router
	router := chi.NewRouter()
	router.Get("/health", handleHealthCheck)
	router.Get("/users/{id}", handleGetUser)

	addr := fmt.Sprintf(":%d", servicePort)
	log.Printf("Starting %s on %s", serviceName, addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// handleGetUser retrieves user information by ID.
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "User Data - Service: %s, ID: %s\n", serviceName, id)
}

// handleHealthCheck returns service health status.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
// registerWithConsul registers the service instance with Consul.
func registerWithConsul() error {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return fmt.Errorf("consul client init failed: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("hostname lookup failed: %w", err)
	}

	reg := &consulapi.AgentServiceRegistration{
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

	if err := client.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	log.Printf("Registered %s on %s:%d", serviceName, hostname, servicePort)
	return nil
}log.Printf("Successfully registered '%s' with Consul", serviceName)
	return nil
}