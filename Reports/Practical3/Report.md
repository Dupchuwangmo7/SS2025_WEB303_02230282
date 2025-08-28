# Practical 3 - Microservices Architecture with gRPC and Consul

## Practical Overview

This practical demonstrates a complete microservices architecture implementation using Go, gRPC, Protocol Buffers, Consul for service discovery, and PostgreSQL databases. The system consists of multiple services that communicate through gRPC and are exposed via an HTTP API Gateway.

## Architecture

### System Components

1. **API Gateway** - HTTP REST API that acts as a single entry point
2. **Users Service** - Manages user data (gRPC service)
3. **Products Service** - Manages product data (gRPC service)
4. **Consul** - Service discovery and health checking
5. **PostgreSQL Databases** - Separate databases for users and products

### Technology Stack

- **Language**: Go 1.23.8
- **Communication**: gRPC with Protocol Buffers
- **Service Discovery**: HashiCorp Consul
- **Databases**: PostgreSQL 13 with GORM
- **Containerization**: Docker & Docker Compose
- **Code Generation**: Buf for Protocol Buffers

## Practical Structure

```
practical-3/
├── proto/                           # Protocol Buffer definitions
│   ├── users.proto                 # User service contract
│   ├── products.proto              # Product service contract
│   ├── gen/                        # Generated Go code
│   │   └── proto/
│   │       ├── users_grpc.pb.go
│   │       ├── users.pb.go
│   │       ├── products_grpc.pb.go
│   │       └── products.pb.go
├── services/
│   ├── users-service/              # User microservice
│   │   ├── main.go
│   │   ├── go.mod
│   │   └── Dockerfile
│   └── products-service/           # Product microservice
│       ├── main.go
│       ├── go.mod
│       └── Dockerfile
├── api-gateway/                    # HTTP to gRPC gateway
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml              # Service orchestration
├── test-api.sh                     # Automated API testing
├── buf.yaml                        # Buf configuration
├── buf.gen.yaml                    # Code generation config
└── Makefile                        # Build automation
```

## What I Implemented

### 1. Protocol Buffer Definitions

Created service contracts using Protocol Buffers in [`proto/users.proto`](proto/users.proto) and [`proto/products.proto`](proto/products.proto):

- **User Service**: `CreateUser` and `GetUser` operations
- **Product Service**: `CreateProduct` and `GetProduct` operations
- Generated Go code using Buf tool for type-safe gRPC communication

### 2. Microservices Implementation

#### Users Service ([`services/users-service/main.go`](services/users-service/main.go))
- gRPC server listening on port 50051
- PostgreSQL database integration with GORM
- Automatic service registration with Consul
- Database retry logic for reliable startup
- User model with name and email fields

#### Products Service ([`services/products-service/main.go`](services/products-service/main.go))
- gRPC server listening on port 50052
- Separate PostgreSQL database for products
- Consul service registration
- Product model with name and price fields

### 3. Service Discovery with Consul

Implemented automatic service registration:
- Services register themselves with Consul on startup
- Health checking and service discovery
- Dynamic service location for the API Gateway

### 4. Database Architecture

- **Separate databases** for each service (microservices best practice)
- **Users DB**: Port 5435, database: `users_db`
- **Products DB**: Port 5434, database: `products_db`
- **GORM integration** for ORM functionality
- **Auto-migration** for database schema management

### 5. Containerization

Created Docker configurations:
- Individual Dockerfiles for each service
- Multi-stage builds for optimized images
- [`docker-compose.yml`](docker-compose.yml) for orchestration
- Proper service dependencies and networking

### 6. API Testing

Developed comprehensive testing in [`test-api.sh`](test-api.sh):
- Automated API endpoint testing
- User creation and retrieval
- Product creation and retrieval
- Combined purchase data testing
- JSON formatting with `jq`

### 7. Development Workflow

Created [`Makefile`](Makefile) with commands:
- `make build` - Build and start all services
- `make test` - Run API tests
- `make proto` - Generate Protocol Buffer code
- `make logs` - View service logs
- `make clean` - Clean up containers and volumes

### 8. Protocol Buffer Code Generation

Configured Buf tool ([`buf.gen.yaml`](buf.gen.yaml)):
- Automatic Go code generation from `.proto` files
- gRPC service and client code generation
- Type-safe message structures

## Key Features Implemented

###  Microservices Architecture
- Independent, loosely coupled services
- Service-specific databases
- gRPC inter-service communication

###  Service Discovery
- Consul integration for dynamic service location
- Automatic service registration and health checking

###  Database Design
- Database per service pattern
- GORM for database operations
- Connection retry logic

###  API Gateway Pattern
- Single entry point for external clients
- HTTP to gRPC translation
- Centralized request routing

###  Containerization
- Docker containers for all services
- Docker Compose orchestration
- Production-ready multi-stage builds

### Testing & Automation
- Automated API testing script
- Build automation with Makefile
- Development workflow documentation

## How to Run

### Prerequisites
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Quick Start

1. **Clone and navigate to the project**:
   ```bash
   cd /home/dupchuwangmo/Desktop/SEM\ 5/WEB303/Practicals/practical-3
   ```

2. **Start all services**:
   ```bash
   make build
   # or
   docker-compose up --build
   ```

3. **Test the API**:
   ```bash
   make test
   # or
   ./test-api.sh
   ```

4. **View Consul UI**:
   Open http://localhost:8500

### API Endpoints

- **Users**:
  - `POST /api/users` - Create user
  - `GET /api/users/{id}` - Get user

![alt text](<images/Screenshot from 2025-08-28 14-05-29.png>)

- **Products**:
  - `POST /api/products` - Create product
  - `GET /api/products/{id}` - Get product

![alt text](<images/Screenshot from 2025-08-28 14-06-02.png>)

- **Combined**:
  - `GET /api/purchases/user/{userId}/product/{productId}` - Get purchase data

![alt text](<images/Screenshot from 2025-08-28 14-06-56.png>)

## Technical Achievements

1. **Implemented gRPC Communication**: Efficient binary protocol for inter-service communication
2. **Service Discovery**: Dynamic service location using Consul
3. **Database Isolation**: Each service has its own database
4. **Container Orchestration**: Multi-container application with proper dependencies
5. **Code Generation**: Automated Protocol Buffer code generation
6. **Testing Automation**: Comprehensive API testing script
7. **Development Workflow**: Complete development environment setup

## Learning Outcomes

Through this project, I gained hands-on experience with:
- Microservices architecture patterns
- gRPC and Protocol Buffers
- Service discovery with Consul
- Container orchestration with Docker Compose
- Database design for microservices
- API Gateway implementation
- Go programming for distributed systems
- Development workflow automation



