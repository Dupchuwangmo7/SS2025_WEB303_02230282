# Web303 Practical 5: Microservices Architecture Report

## Executive Summary

This report documents the implementation of a **Student Cafe Management System** using a microservices architecture with Go, PostgreSQL, and Docker containerization. The system demonstrates the transition from a monolithic architecture to a distributed, scalable microservices-based approach.

---

## 1. Project Overview

### Objective

Refactor a monolithic student cafe application into independent microservices that communicate through REST APIs, managed by an API Gateway.

### Technology Stack

- **Language:** Go (Golang)
- **Framework:** Chi (HTTP router)
- **Database:** PostgreSQL (separate databases per service)
- **Container Orchestration:** Docker & Docker Compose
- **Architecture Pattern:** Microservices with API Gateway

---

## 2. Architecture Design

### System Components

#### 2.1 API Gateway

- **Port:** 8080
- **Responsibility:** Routes incoming requests to appropriate microservices
- **Features:**
  - Path-based routing (`/api/users/*`, `/api/menu/*`, `/api/orders/*`)
  - Reverse proxy implementation using Go's `httputil` package
  - Request logging and monitoring
  - Single entry point for all client requests

#### 2.2 User Service

- **Port:** 8081
- **Database:** `user_db` (PostgreSQL)
- **Endpoints:**
  - `POST /users` - Create new user
  - `GET /users/{id}` - Retrieve specific user
  - `GET /users` - Retrieve all users
- **Responsibilities:**
  - User account management
  - User data persistence
  - User validation for order operations

#### 2.3 Menu Service

- **Port:** 8082
- **Database:** `menu_db` (PostgreSQL)
- **Endpoints:**
  - `GET /menu/{id}` - Get menu item by ID
  - `POST /menu` - Create new menu item
- **Responsibilities:**
  - Menu item management
  - Price management for menu items
  - Menu data retrieval

#### 2.4 Order Service

- **Port:** 8083
- **Database:** `order_db` (PostgreSQL)
- **Endpoints:**
  - `POST /orders` - Create new order
  - `GET /orders` - Retrieve all orders
- **Responsibilities:**
  - Order processing and management
  - Inter-service communication (validates users and menu items)
  - Order persistence with pricing snapshots

### Communication Flow

```
Client Request
    ↓
API Gateway (Port 8080)
    ↓
    ├─→ User Service (8081)
    ├─→ Menu Service (8082)
    └─→ Order Service (8083)
         ├─→ Validates with User Service
         └─→ Validates with Menu Service
```

---

## 3. Database Schema

### User Service Database (`user_db`)

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Menu Service Database (`menu_db`)

```sql
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    availability BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Service Database (`order_db`)

```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    total_price DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);
```

---

## 4. Implementation Details

### 4.1 API Gateway Implementation

**Key Features:**

- Uses Chi router for HTTP routing
- Implements reverse proxy pattern
- Middleware for request logging
- Error handling for malformed URLs
- URL path rewriting to strip `/api` prefix

**Code Structure:**

```go
- main.go: Entry point, router configuration, proxy functions
```

### 4.2 User Service Implementation

**Handlers:**

- `CreateUser()` - Validates input, creates user record
- `GetUser()` - Retrieves single user by ID
- `GetUsers()` - Retrieves all users

**Database Operations:**

- Uses GORM ORM for database interactions
- Implements error handling for database failures
- Uses WHERE clause for explicit query construction

### 4.3 Menu Service Implementation

**Handlers:**

- `GetMenu()` - Retrieves menu item by ID with explicit parameter extraction
- `CreateMenu()` - Validates and persists menu item data

**Database Operations:**

- GORM implementation with improved error messages
- Separate database connection for data isolation

### 4.4 Order Service Implementation

**Handlers:**

- `CreateOrder()` - Complex order creation with service-to-service communication
- `GetOrders()` - Retrieves orders with associated order items

**Inter-Service Communication:**

- Validates user existence via User Service
- Validates menu items via Menu Service
- Captures price snapshots at order creation time
- Handles service failures gracefully

---

## 5. Deployment & Containerization

### Docker Configuration

Each service has:

- **Dockerfile** - Multi-stage build for minimal image size
- **Docker Compose** - Orchestrates all services and databases
- **Network Isolation** - Services communicate via Docker network

### docker-compose.yml Structure

```yaml
Services:
  - api-gateway: Routes requests
  - user-service: User management
  - menu-service: Menu management
  - order-service: Order processing
  - PostgreSQL Databases: Separate DB instances
```

### Environment Configuration

- Connection strings via environment variables
- Default PostgreSQL credentials
- Port mappings for external access

---

## 6. Key Design Patterns

### 6.1 Microservices Pattern

- **Benefit:** Independent scaling, technology flexibility, isolated failures
- **Implementation:** Separate services with dedicated databases

### 6.2 API Gateway Pattern

- **Benefit:** Single entry point, simplified client logic, request routing
- **Implementation:** Reverse proxy with path-based routing

### 6.3 Database per Service

- **Benefit:** Data autonomy, independent scaling, loose coupling
- **Consideration:** Complex queries requiring joins must use APIs

### 6.4 Service-to-Service Communication

- **Benefit:** Decoupled systems, independent deployment
- **Implementation:** HTTP REST calls with error handling

---

## 7. Error Handling & Resilience

### Error Handling Strategy

- Database errors return HTTP 500 with descriptive messages
- Validation errors return HTTP 400
- Resource not found returns HTTP 404
- Service-to-service failures gracefully handled

### Resilience Considerations

- Order Service validates dependencies before processing
- Detailed error messages for debugging
- Explicit HTTP status codes for client handling

---

## 8. Code Quality Improvements

### Enhancements Made

1. **Error Handling** - Added explicit error handling with descriptive messages
2. **Variable Naming** - Improved clarity (e.g., `userData` vs `user`)
3. **Query Construction** - Used explicit WHERE clauses for maintainability
4. **HTTP Status Codes** - Explicit `WriteHeader()` calls for consistency
5. **Logging** - Enhanced logging for better traceability

---

## 9. Testing & Validation

### API Endpoints Testing

**User Service:**

```bash
# Create user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

# Get all users
curl http://localhost:8080/api/users

# Get specific user
curl http://localhost:8080/api/users/1
```

**Menu Service:**

```bash
# Create menu item
curl -X POST http://localhost:8080/api/menu \
  -H "Content-Type: application/json" \
  -d '{"item_name":"Burger","price":5.99}'

# Get menu item
curl http://localhost:8080/api/menu/1
```

**Order Service:**

```bash
# Create order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"items":[{"menu_item_id":1,"quantity":2}]}'

# Get all orders
curl http://localhost:8080/api/orders
```

---

## 10. Challenges & Solutions

| Challenge                | Solution                                                 |
| ------------------------ | -------------------------------------------------------- |
| Service Communication    | Implemented HTTP-based REST APIs with error handling     |
| Data Consistency         | Snapshot prices at order time to prevent inconsistencies |
| Distributed Failures     | Added validation before order creation                   |
| Database Isolation       | Separate databases per service for autonomy              |
| Configuration Management | Environment variables for flexible deployment            |

---

## 11. Scalability Considerations

### Horizontal Scaling

- Each service can be scaled independently
- Multiple instances behind load balancers
- Stateless design enables easy scaling

### Performance Optimization

- Database indexing on frequently queried fields (user_id, menu_item_id)
- Connection pooling through GORM
- Reverse proxy caching potential at gateway level

### Future Improvements

- Message queues for asynchronous communication
- Service mesh for advanced networking
- Distributed tracing and monitoring
- API rate limiting and throttling

---

## 12. Security Considerations

### Current Implementation

- Basic HTTP communication
- Database access controls via connection strings
- Input validation through request parsing

### Recommendations

- Implement HTTPS/TLS for service communication
- Add authentication/authorization (JWT tokens)
- Implement CORS policies
- Add API rate limiting
- Sanitize user input
- Use secrets management for credentials

---

## 13. Monitoring & Logging

### Implemented Features

- Request logging via Chi middleware
- Service startup logging
- Error logging with descriptive messages

### Future Enhancements

- Centralized logging (ELK stack, CloudWatch)
- Distributed tracing (Jaeger, Zipkin)
- Metrics collection (Prometheus)
- Health check endpoints
- Performance monitoring

---

## 14. Conclusion

This microservices implementation successfully demonstrates:

- Separation of concerns across independent services
- Scalable architecture with independent databases
- API Gateway pattern for unified client access
- Service-to-service communication with validation
- Docker containerization for consistent deployment
- Error handling and resilience patterns

The architecture provides a solid foundation for a student cafe management system with room for enhancement in security, monitoring, and advanced deployment patterns.

---

## 15. Appendix: Project Structure

```
Web303_p5/
├── api-gateway/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── menu-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── menu_handlers.go
│   └── models/
│       └── menu.go
├── order-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── order_handlers.go
│   └── models/
│       └── order.go
├── user-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── user_handlers.go
│   └── models/
│       └── user.go
├── docker-compose.yml
├── README.md
└── REPORT.md
```
