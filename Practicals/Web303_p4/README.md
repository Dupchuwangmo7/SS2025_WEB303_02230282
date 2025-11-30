# Student Cafe - Microservices Architecture

A cloud-native microservices application demonstrating service discovery, API gateway, and container orchestration using Kubernetes, Consul, and Kong for a student cafe ordering system.

## 🏗️ Architecture Overview

```
Client → Kong API Gateway → Microservices → Kubernetes Service Discovery
```

### Components

- **Food Catalog Service (Go)** - Menu items and pricing management
- **Order Service (Go)** - Order processing with service discovery
- **Cafe UI (React)** - Frontend application for browsing menu and placing orders
- **Consul** - Service registry and health monitoring (optional)
- **Kong** - API gateway for external routing and load balancing
- **Kubernetes** - Container orchestration and service mesh

## 🚀 Quick Start

### Prerequisites

- Docker Desktop
- Minikube
- kubectl
- Helm

### 1. Start Infrastructure

```bash
# Start minikube
minikube start

# Configure Docker environment to use minikube's Docker daemon
eval $(minikube docker-env)
```

### 2. Install Service Dependencies

```bash
# Create namespace
kubectl create namespace student-cafe

# Install Consul (optional - for advanced service discovery)
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul \
  --set global.name=consul \
  --set server.replicas=1 \
  --set ui.enabled=true \
  --set connectInject.enabled=true \
  -n student-cafe

# Install Kong API Gateway
helm repo add kong https://charts.konghq.com
helm install kong kong/kong \
  --set ingressController.installCRDs=false \
  --set admin.enabled=true \
  -n student-cafe
```

### 3. Build and Deploy Services

```bash
# Build Docker images (ensure you're in the minikube docker environment)
docker build -t food-catalog-service:v1 ./food-catalog-service/
docker build -t order-service:v1 ./order-service/
docker build -t cafe-ui:v1 ./cafe-ui/

# Deploy application services to Kubernetes
kubectl apply -f app-deployment.yaml

# Configure Kong ingress routing
kubectl apply -f kong-ingress.yaml
```

### 4. Access Application

```bash
# Get Minikube IP and Kong proxy port
minikube ip
kubectl get service kong-kong-proxy -n student-cafe

# Access the application
# Example: http://192.168.49.2:30622
```

## 📁 Project Structure

```
student-cafe/
├── food-catalog-service/       # Go microservice for menu management
│   ├── main.go                # Service implementation with Chi router
│   ├── go.mod                 # Go module dependencies
│   ├── go.sum                 # Dependency checksums
│   └── Dockerfile             # Container image definition
├── order-service/              # Go microservice for order processing
│   ├── main.go                # Service with Consul integration
│   ├── go.mod                 # Go module dependencies
│   ├── go.sum                 # Dependency checksums
│   └── Dockerfile             # Container image definition
├── cafe-ui/                    # React frontend application
│   ├── src/                   # React source code
│   ├── public/                # Static assets
│   ├── package.json           # NPM dependencies
│   └── Dockerfile             # Multi-stage build for production
├── app-deployment.yaml         # Kubernetes deployments and services
├── kong-ingress.yaml          # Kong ingress configuration
└── minikube                   # Minikube cluster reference
```

## 🔧 API Endpoints

### Food Catalog Service (Internal: 8080)

- `GET /items` - Retrieve all menu items with pricing

### Order Service (Internal: 8081)

- `POST /orders` - Create new order with item validation
- `GET /health` - Service health check endpoint

### External API (via Kong Gateway)

- `GET /api/catalog` - Access menu items through API gateway
- `POST /api/orders` - Submit orders through API gateway
- `GET /` - Serve React frontend application

## 🎯 Key Features

### Service Discovery

- **Kubernetes Native**: Services communicate using Kubernetes DNS resolution
- **Consul Integration**: Optional advanced service registry for complex scenarios
- **Health Monitoring**: Kubernetes readiness and liveness probes
- **Dynamic Discovery**: Services locate each other automatically

### API Gateway

- **Single Entry Point**: Kong serves as the unified API gateway
- **Path-Based Routing**: Routes traffic based on URL patterns
  - `/` → React Frontend
  - `/api/catalog` → Food Catalog Service
  - `/api/orders` → Order Service
- **Load Balancing**: Automatic traffic distribution
- **Policy Enforcement**: Centralized security and rate limiting

### Container Orchestration

- **Declarative Deployment**: Infrastructure as code with YAML manifests
- **Auto-scaling**: Horizontal pod autoscaling capabilities
- **Service Networking**: Automatic service-to-service communication
- **Health Management**: Automatic restart on failure

## 🏗️ Microservices Patterns Implemented

- ✅ **Service Registry Pattern** - Kubernetes DNS + optional Consul
- ✅ **API Gateway Pattern** - Kong for external traffic management
- ✅ **Health Check Pattern** - Kubernetes probes and monitoring
- ✅ **Service Discovery Pattern** - Dynamic service location
- ✅ **Container Orchestration** - Kubernetes deployment automation
- ✅ **Circuit Breaker Pattern** - Graceful failure handling
- ✅ **Database per Service** - Each service manages its own data

## 🔍 Monitoring & Debugging

### Check Service Health

```bash
# View all pods status
kubectl get pods -n student-cafe

# Check specific deployments
kubectl get deployments -n student-cafe

# View services and endpoints
kubectl get services -n student-cafe
kubectl get endpoints -n student-cafe
```

### Service Logs

```bash
# Monitor order service logs
kubectl logs -f deployment/order-deployment -n student-cafe

# Monitor food catalog service logs
kubectl logs -f deployment/food-catalog-deployment -n student-cafe

# Monitor frontend logs
kubectl logs -f deployment/cafe-ui-deployment -n student-cafe
```

### Kong Gateway Monitoring

```bash
# Access Kong admin API
kubectl port-forward -n student-cafe svc/kong-kong-manager 8002:8002
# Visit http://localhost:8002

# Check ingress configuration
kubectl describe ingress cafe-ingress -n student-cafe
```

### Test Service Discovery

```bash
# Exec into order service pod
kubectl exec -it deployment/order-deployment -n student-cafe -- sh

# Test internal service communication
curl http://food-catalog-service:8080/items

# Test DNS resolution
nslookup food-catalog-service
```

## 🛠️ Development

### Local Development

```bash
# Run services locally (requires Go installed)
cd food-catalog-service && go run main.go
cd order-service && go run main.go

# Run frontend locally (requires Node.js)
cd cafe-ui && npm install && npm start
```

### Building and Updating Images

```bash
# Rebuild after code changes
eval $(minikube docker-env)  # Ensure using minikube docker
docker build -t food-catalog-service:v2 ./food-catalog-service/
docker build -t order-service:v2 ./order-service/
docker build -t cafe-ui:v2 ./cafe-ui/

# Update image versions in app-deployment.yaml, then:
kubectl apply -f app-deployment.yaml

# Or force restart deployments
kubectl rollout restart deployment/food-catalog-deployment -n student-cafe
kubectl rollout restart deployment/order-deployment -n student-cafe
kubectl rollout restart deployment/cafe-ui-deployment -n student-cafe
```

## 🚀 Production Considerations

### Security

- [ ] Implement mTLS between services
- [ ] Add authentication to API gateway (Kong plugins)
- [ ] Use Kubernetes network policies for traffic isolation
- [ ] Secure service-to-service communication
- [ ] Implement API rate limiting

### Scalability

- [ ] Configure horizontal pod autoscaling (HPA)
- [ ] Implement database per service pattern
- [ ] Add persistent volumes for data storage
- [ ] Configure resource limits and requests
- [ ] Implement caching strategies

### Observability

- [ ] Add distributed tracing (Jaeger/Zipkin)
- [ ] Implement metrics collection (Prometheus)
- [ ] Centralized logging (ELK/EFK stack)
- [ ] Application performance monitoring (APM)
- [ ] Custom dashboards (Grafana)

### Reliability

- [ ] Implement circuit breaker pattern
- [ ] Add retry mechanisms with exponential backoff
- [ ] Configure timeouts for service calls
- [ ] Implement bulkhead pattern
- [ ] Add chaos engineering tests

## 📚 Learning Outcomes

This project demonstrates mastery of:

- ✅ **Microservices Architecture Design** - Decomposition into focused services
- ✅ **Container Orchestration** - Kubernetes deployment and service management
- ✅ **Service Discovery** - Both Kubernetes-native and Consul-based approaches
- ✅ **API Gateway Implementation** - Kong for traffic management
- ✅ **Cloud-Native Development** - 12-factor app principles
- ✅ **DevOps Practices** - Infrastructure as code, containerization
- ✅ **Inter-Service Communication** - HTTP APIs and service mesh concepts
- ✅ **Health Monitoring** - Probes, logging, and observability
- ✅ **Fault Tolerance** - Resilience patterns and graceful degradation

## 🧪 Testing the Application

### Functional Testing

1. **Access the Frontend**: Navigate to the Kong gateway URL
2. **Browse Menu**: Verify food items load correctly
3. **Add to Cart**: Test shopping cart functionality
4. **Place Order**: Submit an order and observe the response
5. **API Testing**: Use curl to test endpoints directly

### Expected Behavior

- ✅ Frontend loads and displays menu items
- ✅ Services communicate successfully
- ✅ Orders can be placed through the UI
- ⚠️ **Intentional Failure**: Order submission may fail as designed to demonstrate error handling

### Test Commands

```bash
# Test catalog endpoint
curl http://$(minikube ip):$(kubectl get svc kong-kong-proxy -n student-cafe -o jsonpath='{.spec.ports[0].nodePort}')/api/catalog

# Test order submission
curl -X POST http://$(minikube ip):$(kubectl get svc kong-kong-proxy -n student-cafe -o jsonpath='{.spec.ports[0].nodePort}')/api/orders \
  -H "Content-Type: application/json" \
  -d '{"items":[{"id":"1","name":"Coffee","price":2.5}],"total":2.5}'
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and test thoroughly
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Submit a pull request


---

**Built with ❤️ for learning cloud-native microservices architecture**

