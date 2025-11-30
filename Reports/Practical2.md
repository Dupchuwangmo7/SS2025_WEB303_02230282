# Web303 Practical 2: Microservices Architecture Report

## Executive Summary

This report documents the implementation of a microservices-based architecture using Go, featuring an API Gateway, Users Service, and Products Service. The system utilizes Consul for service discovery and implements health checks for monitoring service availability.

---

## Project Overview

### Architecture Design

The project implements a distributed microservices pattern with the following components:

1. **API Gateway** - Central routing and request forwarding service
2. **Users Service** - Handles user-related operations
3. **Products Service** - Manages product information

### Technology Stack

- **Language**: Go (Golang)
- **Service Discovery**: HashiCorp Consul
- **HTTP Router**: Chi (chi/v5)
- **Communication**: HTTP REST with reverse proxying
- **Architecture Pattern**: Microservices with API Gateway

---

## System Components

### 1. API Gateway (`api-gateway/main.go`)

**Port**: 8080

**Responsibilities**:

- Receives all incoming HTTP requests
- Routes requests to appropriate microservices based on URL path
- Implements reverse proxy forwarding
- Service discovery via Consul

**Key Features**:

- Path-based routing: `/api/{service}/{resource}`
- Automatic service lookup from Consul
- Request path rewriting before forwarding
- Error handling for service unavailability

**Implementation Details**:

```
Request Flow:
/api/users/123 → Gateway → Consul Lookup → users-service:8081 → /123
/api/products/456 → Gateway → Consul Lookup → products-service:8082 → /456
```

### 2. Users Service (`services/users-service/main.go`)

**Port**: 8081

**Endpoints**:

- `GET /health` - Service health check (for Consul)
- `GET /users/{id}` - Retrieve user information by ID

**Functionality**:

- Registers itself with Consul on startup
- Provides health check endpoint for Consul monitoring
- Returns user data in response to requests
- Implements HTTP header management

**Registration Details**:

- Service Name: `users-service`
- Health Check: HTTP endpoint at `/health`
- Check Interval: 10 seconds
- Check Timeout: 1 second

### 3. Products Service (`services/products-service/main.go`)

**Port**: 8082

**Endpoints**:

- `GET /health` - Service health check (for Consul)
- `GET /products/{id}` - Retrieve product information by ID

**Functionality**:

- Identical structure to Users Service
- Registers with Consul for service discovery
- Handles product-related requests
- Implements consistent health monitoring

**Registration Details**:

- Service Name: `products-service`
- Health Check: HTTP endpoint at `/health`
- Check Interval: 10 seconds
- Check Timeout: 1 second

---

## Service Discovery Implementation

### Consul Integration

All services use HashiCorp Consul for dynamic service discovery:

**Key Benefits**:

1. **Dynamic Registration**: Services register on startup
2. **Health Monitoring**: Consul periodically checks service health
3. **Automatic Deregistration**: Failed services are removed from registry
4. **Load Balancing**: Gateway can select from multiple healthy instances

**Registration Process**:

1. Service connects to Consul (default: localhost:8500)
2. Registers service metadata (name, port, address)
3. Provides health check endpoint
4. Consul periodically validates health

**Discovery Process**:

1. API Gateway receives request
2. Queries Consul for healthy service instances
3. Selects first available instance
4. Forwards request to selected service

---

## Request Flow Diagram

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ GET /api/users/123
       ▼
┌──────────────────────┐
│   API Gateway        │
│  (Port 8080)         │
└──────┬───────────────┘
       │ 1. Parse path
       │ 2. Extract service: "users-service"
       ▼
┌──────────────────────┐
│  Consul Registry     │
│  (Service Discovery) │
└──────┬───────────────┘
       │ 3. Query healthy instances
       │ 4. Return users-service:8081
       ▼
┌──────────────────────┐
│  Users Service       │
│  (Port 8081)         │
│  GET /123            │
└──────┬───────────────┘
       │ Response
       ▼
┌─────────────┐
│   Client    │
└─────────────┘
```

---

## Code Structure and Organization

### File Organization

```
Web303_p2/
├── README.md                          # Project documentation
├── REPORT.md                          # This report
├── api-gateway/
│   ├── go.mod                         # Go module definition
│   └── main.go                        # Gateway implementation
├── services/
│   ├── products-service/
│   │   ├── go.mod                     # Module definition
│   │   └── main.go                    # Products service
│   └── users-service/
│       ├── go.mod                     # Module definition
│       └── main.go                    # Users service
```

### Refactoring and Code Quality

**Improvements Made**:

1. **Consistent Naming Conventions**: Standardized variable and function names across services
2. **Simplified Error Handling**: Descriptive error messages with wrapped errors
3. **Code Organization**: Logical function ordering and clear separation of concerns
4. **Response Formatting**: Proper HTTP headers and consistent response formats
5. **Logging**: Improved logging messages for debugging and monitoring

**Function Naming**:

- Service Registration: `registerWithConsul()`
- Handler Functions: `handle{Action}` pattern (e.g., `handleGetUser`, `handleHealthStatus`)
- Discovery: `discoverService()`

---

## API Endpoints Reference

### API Gateway

| Method | Path              | Description                |
| ------ | ----------------- | -------------------------- |
| ANY    | `/api/users/*`    | Routes to Users Service    |
| ANY    | `/api/products/*` | Routes to Products Service |

### Users Service

| Method | Path          | Description    | Response                                       |
| ------ | ------------- | -------------- | ---------------------------------------------- |
| GET    | `/health`     | Health check   | "OK"                                           |
| GET    | `/users/{id}` | Get user by ID | "User Data - Service: users-service, ID: {id}" |

### Products Service

| Method | Path             | Description       | Response                                                     |
| ------ | ---------------- | ----------------- | ------------------------------------------------------------ |
| GET    | `/health`        | Health check      | "Healthy"                                                    |
| GET    | `/products/{id}` | Get product by ID | "Product Info - Service: products-service, Product ID: {id}" |

---

## Setup and Deployment

### Prerequisites

- Go 1.16 or higher
- HashiCorp Consul (running on localhost:8500)
- Network connectivity between services

### Running the Services

**1. Start Consul** (if not already running):

```bash
consul agent -server -ui -bootstrap-expect=1 -data-dir=/tmp/consul
```

**2. Start API Gateway**:

```bash
cd api-gateway
go run main.go
```

**3. Start Users Service** (in another terminal):

```bash
cd services/users-service
go run main.go
```

**4. Start Products Service** (in another terminal):

```bash
cd services/products-service
go run main.go
```

### Testing the System

**Test Users Service**:

```bash
curl http://localhost:8080/api/users/123
# Response: User Data - Service: users-service, ID: 123
```

**Test Products Service**:

```bash
curl http://localhost:8080/api/products/456
# Response: Product Info - Service: products-service, Product ID: 456
```

**Check Health**:

```bash
curl http://localhost:8081/health
# Response: OK

curl http://localhost:8082/health
# Response: Healthy
```

---

## Error Handling and Resilience

### Error Scenarios

1. **Service Not Found**:

   - Gateway returns HTTP 400 (Bad Request) for invalid paths
   - Returns HTTP 503 (Service Unavailable) if service not found in Consul

2. **Unhealthy Service**:

   - Consul automatically removes unhealthy services from registry
   - Gateway will report service unavailable until service recovers

3. **Connection Errors**:
   - Wrapped error messages provide detailed diagnostics
   - Proper HTTP status codes indicate failure type

### Resilience Features

- **Health Checks**: Continuous monitoring via Consul
- **Service Registry**: Dynamic registration/deregistration
- **Error Propagation**: Clear error messages for debugging
- **Graceful Failures**: Services handle network issues appropriately

---

## Performance Considerations

### Request Processing

1. **Path Parsing**: O(1) - simple string split
2. **Service Discovery**: O(n) - linear search through Consul instances
3. **Reverse Proxy**: Efficient HTTP forwarding with minimal overhead
4. **Health Checks**: Asynchronous, non-blocking

### Scalability

- **Horizontal Scaling**: Services can be replicated behind load balancers
- **Multiple Instances**: Gateway can distribute across multiple service instances
- **Consul Clustering**: Support for multi-node Consul deployments

---

## Security Considerations

### Current Implementation

1. **No Authentication**: Services accept all requests (development mode)
2. **No HTTPS**: Plain HTTP communication (localhost development)
3. **No Rate Limiting**: No request throttling implemented

### Recommendations for Production

1. Implement OAuth2/JWT authentication
2. Enable TLS/HTTPS encryption
3. Add API key management
4. Implement rate limiting
5. Add request validation and sanitization
6. Use network policies/firewalls

---

## Monitoring and Debugging

### Consul UI

Access Consul dashboard at: `http://localhost:8500/ui/`

**Available Information**:

- Service registrations
- Service instances and status
- Health check status
- Key-value store

### Logging

Each service logs:

- Service initialization
- Incoming requests
- Service registration status
- Errors and warnings

### Typical Log Output

```
Registered users-service on hostname:8081
Starting users-service on :8081
Incoming request: GET /users/123
Located service at: http://hostname:8081
Proxying to: http://hostname:8081/123
```

---

## Testing Summary

### Manual Testing

- API Gateway routes requests correctly
- Service discovery works via Consul
- Health checks respond properly
- Path rewriting functions correctly
- Error handling for invalid paths
- Service registration and deregistration

### Test Cases

| Test Case              | Expected Result       | Status |
| ---------------------- | --------------------- | ------ |
| Valid user request     | Returns user data     | Pass   |
| Valid product request  | Returns product data  | Pass   |
| Invalid path           | HTTP 400 response     | Pass   |
| Service not registered | HTTP 503 response     | Pass   |
| Health check           | Service responds "OK" | Pass   |

---

## Conclusion

This microservices architecture demonstrates:

- Effective service-to-service communication
- Dynamic service discovery and registration
- API Gateway pattern implementation
- Health monitoring and resilience
- Clean code organization and maintainability

The system is well-structured for development and can be extended with additional services following the same patterns. For production deployment, security measures and performance optimizations should be implemented.

---

## Appendix: Dependencies

### Go Modules

**API Gateway**:

- `github.com/hashicorp/consul/api` - Consul client

**Users Service**:

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/hashicorp/consul/api` - Consul client

**Products Service**:

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/hashicorp/consul/api` - Consul client

### Installation

```bash
go get github.com/hashicorp/consul/api
go get github.com/go-chi/chi/v5
```
