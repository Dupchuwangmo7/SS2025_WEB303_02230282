# Student Café Microservices Platform

A modern, scalable microservices architecture for managing a student café, built with Go, gRPC, and PostgreSQL. The system features an API Gateway layer that translates HTTP requests to gRPC calls, with specialized backend services handling user management, menu operations, and order processing. Each service operates independently with dedicated database connections for optimal performance and data isolation.

## 🏗️ Architecture Overview

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

### Services

1. **API Gateway** - HTTP to gRPC translation layer

   - Port: `8080` (HTTP)
   - Translates REST API calls to gRPC calls for backend services
   - Handles request routing and response translation

2. **User Service** - User management

   - Port: `9091` (gRPC)
   - Database: `user_db` on port `5434`
   - Operations: Create, Read users and café owners

3. **Menu Service** - Menu item management

   - Port: `9092` (gRPC)
   - Database: `menu_db` on port `5433`
   - Operations: Create, Read menu items

4. **Order Service** - Order management
   - Port: `9093` (gRPC)
   - Database: `order_db` on port `5435`
   - Operations: Create, Read orders and order items
   - Communicates with User and Menu services

## 📋 API Endpoints

### User Endpoints

```
POST   /api/users              - Create a new user
GET    /api/users/{id}         - Get user by ID
GET    /api/users              - Get all users
```

### Menu Endpoints

```
POST   /api/menu               - Create a new menu item
GET    /api/menu/{id}          - Get menu item by ID
GET    /api/menu               - Get all menu items
```

### Order Endpoints

```
POST   /api/orders             - Create a new order
GET    /api/orders/{id}        - Get order by ID
GET    /api/orders             - Get all orders
```

## 🚀 Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.18+ (for local development)
- PostgreSQL 13+ (if running without Docker)
- Protocol Buffers compiler (for proto file modifications)

### Quick Start with Docker Compose

```bash
# Start all services and databases
docker-compose up -d

# Services will be available at:
# - API Gateway: http://localhost:8080
# - User Service: localhost:9091 (gRPC)
# - Menu Service: localhost:9092 (gRPC)
# - Order Service: localhost:9093 (gRPC)

# Databases:
# - User DB: localhost:5434
# - Menu DB: localhost:5433
# - Order DB: localhost:5435
```

### Local Development Setup

1. **Start databases only**:

   ```bash
   docker-compose up -d user-db menu-db order-db
   ```

2. **Install dependencies**:

   ```bash
   # Generate gRPC code (if protobuf files changed)
   cd student-cafe-protos
   make generate

   # Install Go dependencies for each service
   cd ../api-gateway && go mod download
   cd ../user-service && go mod download
   cd ../menu-service && go mod download
   cd ../order-service && go mod download
   ```

3. **Run services** (in separate terminals):

   ```bash
   # Terminal 1 - User Service
   cd user-service
   go run main.go

   # Terminal 2 - Menu Service
   cd menu-service
   go run main.go

   # Terminal 3 - Order Service
   cd order-service
   go run main.go

   # Terminal 4 - API Gateway
   cd api-gateway
   go run main.go
   ```

## 🧪 Testing

### Run Integration Tests

```bash
cd tests/integration
go test -v
```

### Run E2E Tests

```bash
cd tests/e2e
go test -v
```

### Run Service Tests

```bash
# User Service
cd user-service
go test -v ./...

# Menu Service
cd menu-service
go test -v ./...

# Order Service
cd order-service
go test -v ./...
```

### View Coverage Reports

```bash
# Coverage reports are pre-generated in each service folder
# Open in browser:
open user-service/coverage.html
open menu-service/coverage.html
open order-service/coverage.html
```

## 🔧 Project Structure

```
project-root/
├── api-gateway/              # HTTP to gRPC gateway
│   ├── main.go              # Service entry point
│   ├── go.mod               # Go dependencies
│   ├── grpc/                # gRPC client connections
│   └── handlers/            # HTTP request handlers
│
├── user-service/            # User management service
│   ├── main.go              # Service entry point
│   ├── database/            # Database connection & queries
│   ├── grpc/                # gRPC server implementation
│   ├── handlers/            # Business logic
│   ├── models/              # Data models
│   └── go.mod               # Go dependencies
│
├── menu-service/            # Menu management service
│   ├── main.go              # Service entry point
│   ├── database/            # Database connection & queries
│   ├── grpc/                # gRPC server implementation
│   ├── handlers/            # Business logic
│   ├── models/              # Data models
│   └── go.mod               # Go dependencies
│
├── order-service/           # Order management service
│   ├── main.go              # Service entry point
│   ├── database/            # Database connection & queries
│   ├── grpc/                # gRPC server implementation
│   ├── handlers/            # Business logic
│   ├── models/              # Data models
│   └── go.mod               # Go dependencies
│
├── student-cafe-protos/     # Protocol Buffer definitions
│   ├── proto/               # .proto files
│   │   ├── menu/v1/
│   │   ├── order/v1/
│   │   └── user/v1/
│   ├── gen/                 # Generated gRPC code
│   └── Makefile             # Proto generation commands
│
├── tests/
│   ├── integration/         # Integration tests
│   └── e2e/                 # End-to-end tests
│
├── docker-compose.yml       # Service orchestration
├── Makefile                 # Build commands
└── README.md               # This file
```

## 🐳 Docker Compose Services

The `docker-compose.yml` defines:

- **3 PostgreSQL databases** (user_db, menu_db, order_db)
- **4 Go services** (api-gateway, user-service, menu-service, order-service)
- **Shared network** (cafe-network) for inter-service communication

### Environment Variables

Services use the following environment variables (defaults provided):

```bash
# User Service
DATABASE_URL=host=localhost user=postgres password=postgres dbname=user_db port=5432 sslmode=disable
GRPC_PORT=9091

# Menu Service
DATABASE_URL=host=localhost user=postgres password=postgres dbname=menu_db port=5432 sslmode=disable
GRPC_PORT=9092

# Order Service
DATABASE_URL=host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable
GRPC_PORT=9093
USER_SERVICE_GRPC_ADDR=user-service:9091
MENU_SERVICE_GRPC_ADDR=menu-service:9092

# API Gateway
USER_SERVICE_GRPC_ADDR=user-service:9091
MENU_SERVICE_GRPC_ADDR=menu-service:9092
ORDER_SERVICE_GRPC_ADDR=order-service:9093
```

## 📝 Example Requests

### Create a User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "is_cafe_owner": false
  }'
```

### Get All Users

```bash
curl http://localhost:8080/api/users
```

### Create a Menu Item

```bash
curl -X POST http://localhost:8080/api/menu \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Espresso",
    "description": "Strong Italian coffee",
    "price": 3.50
  }'
```

### Create an Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"menu_item_id": 1, "quantity": 2}
    ]
  }'
```

## 🔌 Protocol Buffers

Protocol Buffer definitions are located in `student-cafe-protos/proto/`:

- `menu/v1/menu.proto` - Menu service definitions
- `order/v1/order.proto` - Order service definitions
- `user/v1/user.proto` - User service definitions

### Regenerate gRPC Code

If you modify `.proto` files:

```bash
cd student-cafe-protos
make generate
```

This will regenerate the Go code in `gen/go/`.

## 🛠️ Build and Deployment

### Build Docker Images

```bash
# Build all services
make build

# Or manually:
docker build -f api-gateway/Dockerfile -t student-cafe-api-gateway .
docker build -f user-service/Dockerfile -t student-cafe-user-service .
docker build -f menu-service/Dockerfile -t student-cafe-menu-service .
docker build -f order-service/Dockerfile -t student-cafe-order-service .
```

### Deploy with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 📚 Technology Stack

- **Language**: Go 1.18+
- **API Framework**: Chi (HTTP router)
- **RPC Framework**: gRPC & Protocol Buffers
- **Database**: PostgreSQL 13
- **Containerization**: Docker & Docker Compose
- **Testing**: Go's built-in testing package

## 📦 Dependencies

### Core Dependencies

- `google.golang.org/grpc` - gRPC framework
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/douglasswm/student-cafe-protos` - gRPC definitions

See individual `go.mod` files for complete dependency lists.

## 🤝 Contributing

1. Create a new branch for your feature
2. Make your changes
3. Run tests: `make test`
4. Commit and push
5. Create a pull request

## 📝 License

This project is part of the coursework for Web303 (Practical 5 & 6).

## 🆘 Troubleshooting

### Services can't connect to each other

- Ensure all services are running
- Check the `cafe-network` is created: `docker network ls`
- Verify service addresses in Docker Compose

### Database connection errors

- Ensure database containers are running: `docker-compose ps`
- Check database credentials in environment variables
- Verify database ports are not in use

### gRPC connection refused

- Ensure services are listening on correct ports
- Check firewall rules
- Verify service startup logs: `docker-compose logs <service-name>`

### Port conflicts

- Change ports in `docker-compose.yml` and service `.env` files
- Or stop other services using those ports

## 📞 Support

For issues or questions, check the logs:

```bash
docker-compose logs -f <service-name>
```

---

**Last Updated**: November 2025
