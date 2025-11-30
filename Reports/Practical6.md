# Student Café Microservices Platform - Project Report

## Executive Summary

This project implements a modern, scalable microservices architecture for managing a student café. The system demonstrates professional-grade distributed system design using Go, gRPC, and PostgreSQL, with an HTTP API Gateway serving as the single entry point for client requests.

---

## Project Overview

### Objectives

- Design and implement a microservices-based architecture for café operations
- Implement HTTP to gRPC translation layer for seamless client communication
- Ensure data isolation through service-specific databases
- Demonstrate inter-service communication patterns
- Provide comprehensive API documentation and deployment infrastructure

### Technology Stack

- **Language**: Go (Golang)
- **Communication Protocol**: gRPC with Protocol Buffers
- **API Layer**: HTTP with Chi router
- **Database**: PostgreSQL (3 independent instances)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing framework with integration and e2e test suites

---

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (HTTP)                        │
│                      :8080                                   │
└────────┬─────────────┬─────────────┬──────────────────────┘
         │             │             │
    ┌────▼────┐   ┌────▼─────┐  ┌───▼──────┐
    │  User   │   │   Menu   │  │  Order   │
    │ Service │   │ Service  │  │ Service  │
    │ (gRPC)  │   │  (gRPC)  │  │  (gRPC)  │
    │ :9091   │   │  :9092   │  │  :9093   │
    └────┬────┘   └────┬─────┘  └───┬──────┘
         │             │             │
    ┌────▼────┐   ┌────▼─────┐  ┌───▼──────┐
    │ User DB │   │  Menu DB │  │ Order DB │
    │ (PG)    │   │   (PG)   │  │  (PG)    │
    └─────────┘   └──────────┘  └──────────┘
```

### Architectural Principles

1. **Service Independence**: Each service operates autonomously with dedicated database connections
2. **API Gateway Pattern**: Single entry point translates HTTP REST calls to gRPC for internal communication
3. **Data Isolation**: Service-specific databases prevent tight coupling and enable independent scaling
4. **Protocol Buffer Definition**: Shared proto definitions ensure consistent API contracts across services
5. **Containerized Deployment**: Docker Compose orchestrates multi-container environment for easy deployment

---

## Services Description

### 1. API Gateway

**Purpose**: HTTP to gRPC translation layer and request routing

**Responsibilities**:

- Accept HTTP REST requests on port 8080
- Translate HTTP requests to gRPC calls
- Route requests to appropriate backend services
- Handle error responses and status code mapping
- Provide request logging and recovery middleware

**Technology Stack**:

- Chi router for HTTP routing
- gRPC client connections to backend services
- Middleware: Logger, Recoverer

**Endpoints**:

```
User Operations:
  POST   /api/users              - Create new user
  GET    /api/users/{id}         - Retrieve user by ID
  GET    /api/users              - List all users

Menu Operations:
  POST   /api/menu               - Create menu item
  GET    /api/menu/{id}          - Retrieve menu item by ID
  GET    /api/menu               - List all menu items

Order Operations:
  POST   /api/orders             - Create order
  GET    /api/orders/{id}        - Retrieve order by ID
  GET    /api/orders             - List all orders
```

### 2. User Service

**Purpose**: User and café owner account management

**Responsibilities**:

- User account creation and validation
- Email uniqueness enforcement
- Café owner privilege management
- User data persistence and retrieval

**gRPC Port**: 9091
**Database**: `user_db` (PostgreSQL on port 5434)

**Data Model**:

```go
type User struct {
    ID          uint
    Name        string  // User's display name
    Email       string  // Unique email identifier
    IsCafeOwner bool    // Privilege indicator
    CreatedAt   time.Time
    UpdateedAt  time.Time
}
```

**Key Operations**:

- CreateUser: Register new user with email validation
- GetUser: Retrieve user information by ID
- GetUsers: Fetch all registered users

### 3. Menu Service

**Purpose**: Restaurant menu and item management

**Responsibilities**:

- Menu organization and categorization
- Menu item definition and pricing
- Item availability and description management
- Menu data persistence

**gRPC Port**: 9092
**Database**: `menu_db` (PostgreSQL on port 5433)

**Data Model**:

```go
type Menu struct {
    ID          uint
    Name        string
    Description string
    MenuItems   []MenuItem
    CreatedAt   time.Time
}

type MenuItem struct {
    ID          uint
    Name        string
    Description string
    Price       float64
    CreatedAt   time.Time
}
```

**Key Operations**:

- GetMenuItem: Retrieve specific menu item details
- GetMenu: Fetch complete menu with all items
- CreateMenuItem: Add new item to menu

### 4. Order Service

**Purpose**: Order creation, processing, and fulfillment tracking

**Responsibilities**:

- Order creation and validation
- Order item management with quantity tracking
- Order status tracking
- Inter-service communication with User and Menu services
- Price snapshot preservation at order time

**gRPC Port**: 9093
**Database**: `order_db` (PostgreSQL on port 5435)

**Data Model**:

```go
type Order struct {
    ID         uint
    UserID     uint        // Reference to User Service
    Status     string      // "pending", "completed"
    OrderItems []OrderItem
    CreatedAt  time.Time
}

type OrderItem struct {
    ID         uint
    OrderID    uint
    MenuItemID uint        // Reference to Menu Service
    Quantity   int
    Price      float64     // Snapshot price
    CreatedAt  time.Time
}
```

**Key Operations**:

- CreateOrder: Create order with items and user validation
- GetOrder: Retrieve order details with items
- GetOrders: List all orders with filtering

---

## Protocol Buffer Definitions

### Location

`student-cafe-protos/proto/` directory contains service definitions:

- `user/v1/user.proto` - User service RPC definitions
- `menu/v1/menu.proto` - Menu service RPC definitions
- `order/v1/order.proto` - Order service RPC definitions

### Generated Code

Protocol buffers are compiled to Go code in `student-cafe-protos/gen/go/` with:

- Service stubs and clients
- Message serialization/deserialization
- gRPC service interfaces

### Benefits

- Language-agnostic service contracts
- Efficient binary serialization
- Strong typing and versioning support
- Code generation automation

---

## Database Schema

### User Database (`user_db`)

```sql
-- Users table with email uniqueness
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_cafe_owner BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Menu Database (`menu_db`)

```sql
-- Menus table
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Menu Items table
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Database (`order_db`)

```sql
-- Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order Items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    menu_item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Deployment Architecture

### Docker Compose Configuration

The `docker-compose.yml` orchestrates:

- 3 PostgreSQL database instances (separate databases)
- API Gateway container
- User Service container
- Menu Service container
- Order Service container

### Container Services

**Databases**:

```yaml
PostgreSQL Instances:
  - user-db:5434 (user_db database)
  - menu-db:5433 (menu_db database)
  - order-db:5435 (order_db database)
```

**Applications**:

```yaml
API Gateway:
  - Port: 8080 (HTTP)
  - Environment: Backend service addresses

User Service:
  - Port: 9091 (gRPC)
  - Database: user-db:5434

Menu Service:
  - Port: 9092 (gRPC)
  - Database: menu-db:5433

Order Service:
  - Port: 9093 (gRPC)
  - Services: User Service, Menu Service
  - Database: order-db:5435
```

### Deployment Commands

```bash
# Start all services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f [service_name]

# Rebuild services
docker-compose build
```

---

## Testing Strategy

### Test Coverage

#### 1. Unit Tests

- Service handler logic validation
- Error handling and edge cases
- Database operation correctness

#### 2. Integration Tests

- Inter-service communication validation
- gRPC client/server interaction
- Database integration with services

#### 3. End-to-End Tests

- Complete request flow from HTTP to database
- Multi-service workflows
- Error propagation and recovery

### Test Execution

```bash
# Run all tests
make test

# Run with coverage
make coverage

# View coverage report
open coverage.html
```

---

## Error Handling

### gRPC Error Mapping

The system maps gRPC status codes to HTTP status codes:

```
gRPC Code          →  HTTP Status
NotFound           →  404
InvalidArgument    →  400
AlreadyExists      →  409
PermissionDenied   →  403
Internal           →  500
```

### Error Response Format

```json
{
  "error": "descriptive error message",
  "code": "gRPC error code",
  "details": "additional context"
}
```

---

## Getting Started

### Prerequisites

- Docker & Docker Compose (v3.8+)
- Go 1.18+ (for local development)
- PostgreSQL 13+ (for standalone deployment)
- Protocol Buffers compiler (for proto modifications)

### Quick Start

#### Using Docker Compose

```bash
# Navigate to project root
cd /path/to/Web303_p5-ab_p6

# Start all services
docker-compose up -d

# Verify services are running
docker-compose ps

# Test API Gateway
curl http://localhost:8080/api/users

# Stop services
docker-compose down
```

#### Local Development

```bash
# Setup databases
psql -c "CREATE DATABASE user_db;"
psql -c "CREATE DATABASE menu_db;"
psql -c "CREATE DATABASE order_db;"

# Start services individually (in separate terminals)

# Terminal 1: User Service
cd user-service
go run main.go

# Terminal 2: Menu Service
cd menu-service
go run main.go

# Terminal 3: Order Service
cd order-service
go run main.go

# Terminal 4: API Gateway
cd api-gateway
go run main.go
```

---

## API Usage Examples

### User Management

```bash
# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "is_cafe_owner": false
  }'

# Get user
curl http://localhost:8080/api/users/1

# Get all users
curl http://localhost:8080/api/users
```

### Menu Management

```bash
# Create menu item
curl -X POST http://localhost:8080/api/menu \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Coffee",
    "description": "Espresso-based coffee",
    "price": 3.99
  }'

# Get menu item
curl http://localhost:8080/api/menu/1

# Get menu
curl http://localhost:8080/api/menu
```

### Order Management

```bash
# Create order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"menu_item_id": 1, "quantity": 2}
    ]
  }'

# Get order
curl http://localhost:8080/api/orders/1

# Get orders
curl http://localhost:8080/api/orders
```

---

## Project Structure

```
Web303_p5-ab_p6/
├── api-gateway/              # HTTP API Gateway
│   ├── main.go              # Entry point
│   ├── handlers/            # HTTP request handlers
│   ├── grpc/                # gRPC client code
│   └── Dockerfile           # Container configuration
├── user-service/            # User management service
│   ├── main.go
│   ├── database/            # Database connection
│   ├── grpc/                # gRPC server implementation
│   ├── handlers/            # Service handlers
│   ├── models/              # Data models
│   └── Dockerfile
├── menu-service/            # Menu management service
│   ├── main.go
│   ├── database/
│   ├── grpc/
│   ├── handlers/
│   ├── models/
│   └── Dockerfile
├── order-service/           # Order management service
│   ├── main.go
│   ├── database/
│   ├── grpc/
│   ├── handlers/
│   ├── models/
│   └── Dockerfile
├── student-cafe-protos/     # Protocol Buffer definitions
│   ├── proto/               # .proto files
│   ├── gen/go/              # Generated Go code
│   └── Makefile             # Proto generation
├── tests/                   # Test suites
│   ├── e2e/                 # End-to-end tests
│   └── integration/         # Integration tests
├── docker-compose.yml       # Multi-container orchestration
├── Makefile                 # Build automation
└── README.md                # Project documentation
```

---

## Performance Considerations

### Scalability Features

- Independent service databases enable horizontal scaling
- gRPC binary protocol reduces latency vs HTTP/JSON
- Connection pooling in database layers
- Stateless service design allows load balancing

### Optimization Strategies

- Proto compilation for efficient serialization
- Database indexes on frequently queried fields
- Chi router optimization for request routing
- Middleware caching opportunities

### Load Distribution

- Multiple instances of each service can run behind load balancer
- Database read replicas for query-heavy operations
- API Gateway acts as single point for request distribution

---

## Security Considerations

### Current Implementation

- Email uniqueness ensures user identification
- Service isolation through network separation
- Database access control through credentials
- Error message sanitization to prevent information leakage

### Recommended Enhancements

- Implement authentication/authorization (JWT tokens)
- Add TLS/SSL for inter-service communication
- Implement rate limiting on API Gateway
- Add request validation and sanitization
- Implement audit logging for sensitive operations
- Use environment variables for database credentials

---

## Future Enhancements

### Planned Features

1. **Authentication Layer**: JWT-based authentication with role-based access control
2. **Caching Strategy**: Redis caching for frequently accessed data
3. **Message Queue**: RabbitMQ/Kafka for asynchronous order processing
4. **Service Discovery**: Consul or Eureka for dynamic service registration
5. **API Documentation**: OpenAPI/Swagger integration
6. **Monitoring & Metrics**: Prometheus and Grafana integration
7. **Distributed Tracing**: Jaeger for request tracing
8. **Circuit Breaker**: Resilience4j-like pattern for fault tolerance

### Technical Debt

- Add comprehensive error handling documentation
- Expand test coverage to >80%
- Add request validation middleware
- Implement database migration management

---

## Maintenance & Monitoring

### Health Checks

```bash
# API Gateway health
curl http://localhost:8080/health

# Service logs
docker-compose logs api-gateway
docker-compose logs user-service
docker-compose logs menu-service
docker-compose logs order-service
```

### Database Maintenance

```bash
# Connect to user database
psql -h localhost -p 5434 -U postgres -d user_db

# View tables
\dt

# Backup database
pg_dump -h localhost -p 5434 -U postgres user_db > backup.sql
```

---

## Conclusion

The Student Café Microservices Platform demonstrates a professional-grade implementation of distributed systems architecture. The project effectively showcases:

- Microservices design patterns and principles
- Protocol Buffer usage for efficient inter-service communication
- HTTP to gRPC translation patterns
- Database isolation and service independence
- Docker containerization and orchestration
- Modern Go development practices

The architecture is scalable, maintainable, and ready for production deployment with recommended security enhancements.

---

## References & Resources

- **Go Documentation**: https://golang.org/doc/
- **gRPC Guide**: https://grpc.io/docs/
- **Protocol Buffers**: https://developers.google.com/protocol-buffers
- **Docker Documentation**: https://docs.docker.com/
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/
- **Chi Router**: https://github.com/go-chi/chi
- **GORM**: https://gorm.io/

---


