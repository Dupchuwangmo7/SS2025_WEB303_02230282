// api-gateway/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

const gatewayPort = 8080

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", routeRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", gatewayPort),
		Handler: router,
	}

	log.Printf("API Gateway initializing on port %d...", gatewayPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Gateway startup failed: %v", err)
	}
}

// routeRequest forwards HTTP requests to appropriate microservices based on URL path.
func routeRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

	// Parse the path to extract service name: /api/{service}/{resource}
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 || pathParts[0] != "api" {
		http.Error(w, "Invalid path format", http.StatusBadRequest)
		return
	}
	serviceName := pathParts[1] + "-service"

	// Locate the service in Consul service registry
	targetURL, err := discoverService(serviceName)
	if err != nil {
		log.Printf("Service discovery failed for '%s': %v", serviceName, err)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Located service at: %s", targetURL)

	// Create reverse proxy and adjust the request path
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Remove /api/{service} prefix before forwarding
	r.URL.Path = "/" + strings.Join(pathParts[2:], "/")
	log.Printf("Proxying to: %s%s", targetURL, r.URL.Path)

	reverseProxy.ServeHTTP(w, r)
}

// discoverService queries Consul to find a healthy instance of a service.
func discoverService(name string) (*url.URL, error) {
// discoverService retrieves a service endpoint from Consul.
func discoverService(serviceName string) (*url.URL, error) {
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("consul client error: %w", err)
	}

	// Fetch healthy service entries from Consul
	healthyInstances, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("consul query failed for '%s': %w", serviceName, err)
	}

	if len(healthyInstances) == 0 {
		return nil, fmt.Errorf("no healthy instances available for '%s'", serviceName)
	}

	// Use first available instance
	instance := healthyInstances[0].Service
	endpoint := fmt.Sprintf("http://%s:%d", instance.Address, instance.Port)

	return url.Parse(endpoint)
}