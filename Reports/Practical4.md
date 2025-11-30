# Web303 Practical 4: Campus Cafe Express - Microservices Architecture Report

## Executive Summary

This practical demonstrates the implementation of a cloud-native microservices architecture for a student cafe ordering system. The application showcases modern cloud computing principles including containerization, service discovery, API gateway routing, and container orchestration using Kubernetes. The system is designed to handle distributed service communication with built-in health checks and scalability.

---

## 1. Project Overview

### 1.1 Application Purpose

Campus Cafe Express is a microservices-based ordering system that allows students to:

- Browse available menu items
- Add items to their shopping cart
- Place orders through an API gateway
- Receive real-time order confirmation with unique order IDs

### 1.2 Architecture Pattern

The application follows the **Microservices Architecture Pattern** with the following principles:

- **Separation of Concerns:** Each service handles a specific business domain
- **Independent Deployment:** Services can be deployed and scaled independently
- **API-First Design:** All inter-service communication is via REST APIs
- **Service Discovery:** Dynamic service location resolution
- **API Gateway Pattern:** Single entry point for all client requests

---

## 2. System Architecture

### 2.1 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client (Web Browser)                      │
│                                                                  │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTP/HTTPS
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Kong API Gateway                              │
│         (Routing, Load Balancing, Rate Limiting)                │
│                                                                  │
│  Routes:                                                         │
│  /api/catalog/* → Food Catalog Service (Port 8080)             │
│  /api/orders/* → Order Service (Port 8081)                     │
│  /* → Cafe UI Service (Port 80)                                │
└────────────┬──────────────────┬──────────────────┬──────────────┘
             │                  │                  │
             ▼                  ▼                  ▼
    ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
    │ Cafe UI (React)  │ │ Food Catalog Svc │ │  Order Service   │
    │  Port: 80        │ │   (Go)           │ │    (Go)          │
    │                  │ │   Port: 8080     │ │   Port: 8081     │
    │ - Menu Display   │ │                  │ │                  │
    │ - Cart Mgmt      │ │ - Menu Items     │ │ - Order Creation │
    │ - Order Placement│ │ - Pricing        │ │ - Order Tracking │
    └──────────────────┘ └──────────────────┘ └──────────────────┘
             │                  │                  │
             └──────────────────┬──────────────────┘
                                │
                    ┌───────────┴───────────┐
                    │                       │
                    ▼                       ▼
        ┌──────────────────────┐  ┌──────────────────────┐
        │  Kubernetes Service  │  │  Kubernetes Service  │
        │   Discovery (DNS)    │  │   Health Checks      │
        └──────────────────────┘  └──────────────────────┘
```

### 2.2 Key Components

#### **Frontend Layer - Cafe UI (React)**

- **Technology:** React.js
- **Port:** 80 (HTTP)
- **Responsibilities:**
  - Display available menu items
  - Manage shopping cart state
  - Handle user interactions
  - Submit orders via API Gateway
  - Display order confirmation messages

#### **Catalog Service (Go)**

- **Technology:** Go with Chi router
- **Port:** 8080
- **Responsibilities:**
  - Serve menu items with pricing
  - Provide health checks for monitoring
  - Return JSON formatted food items
- **Endpoints:**
  - `GET /health` - Service health status
  - `GET /items` - List all available food items

#### **Order Service (Go)**

- **Technology:** Go with Chi router & UUID generation
- **Port:** 8081
- **Responsibilities:**
  - Process incoming orders
  - Generate unique order IDs
  - Maintain order records
  - Perform inter-service communication with Catalog Service
- **Endpoints:**
  - `GET /health` - Service health status
  - `POST /orders` - Create new order

---

## 3. Technology Stack

| Component         | Technology | Version  | Purpose                         |
| ----------------- | ---------- | -------- | ------------------------------- |
| Container Runtime | Docker     | Latest   | Service containerization        |
| Orchestration     | Kubernetes | 1.20+    | Service orchestration & scaling |
| Service Gateway   | Kong       | Latest   | API routing & load balancing    |
| Service Discovery | Consul     | Optional | Service registry (optional)     |
| Backend Services  | Go         | 1.16+    | Microservice implementation     |
| Frontend          | React      | 17+      | User interface                  |
| HTTP Router       | Chi        | v5       | Go HTTP routing                 |
| ID Generation     | UUID       | Standard | Unique order identification     |

---

## 4. Implementation Details

### 4.1 Food Catalog Service Implementation

**File:** `food-catalog-service/main.go`

**Menu Items:**

```
1. Espresso - $2.75
2. Turkey Sandwich - $5.50
3. Blueberry Muffin - $3.50
4. Iced Tea - $2.25
5. Caesar Salad - $6.00
```

**Key Features:**

- Lightweight menu management
- HTTP-based endpoints
- Structured JSON responses
- Middleware logging for debugging

**Code Structure:**

```go
type FoodItem struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### 4.2 Order Service Implementation

**File:** `order-service/main.go`

**Order Data Structure:**

```go
type Order struct {
    ID        string   `json:"id"`       // Unique identifier
    ItemIDs   []string `json:"item_ids"` // Requested items
    Status    string   `json:"status"`   // Processing status
    Timestamp string   `json:"timestamp"` // Order creation time
}
```

**Key Features:**

- UUID-based order identification
- Service-to-service discovery
- Order state management
- Inter-service validation

**Service Discovery:**

- Implements service lookup mechanism
- Supports Kubernetes service names
- Handles service resolution errors

### 4.3 Frontend Implementation

**File:** `cafe-ui/src/App.js`

**Key Features:**

- React hooks for state management
- API Gateway integration
- Shopping cart functionality
- Real-time order feedback

**State Management:**

```javascript
const [items, setItems] = useState([]); // Menu items
const [cart, setCart] = useState([]); // Shopping cart
const [message, setMessage] = useState(""); // User feedback
```

**API Endpoints:**

- `GET /api/catalog/items` - Fetch menu items
- `POST /api/orders/orders` - Place new order

---

## 5. Kubernetes Deployment Configuration

### 5.1 Namespace

```yaml
Namespace: student-cafe
```

Creates an isolated environment for all application resources.

### 5.2 Deployments

#### Food Catalog Deployment

- **Replicas:** 1
- **Image:** food-catalog-service:v1
- **Port:** 8080
- **Pull Policy:** IfNotPresent (for local minikube development)

#### Order Service Deployment

- **Replicas:** 1
- **Image:** order-service:v1
- **Port:** 8081
- **Pull Policy:** IfNotPresent

#### Cafe UI Deployment

- **Replicas:** 1
- **Image:** cafe-ui:v1
- **Port:** 80
- **Pull Policy:** IfNotPresent

### 5.3 Kubernetes Services

Each deployment is exposed via a ClusterIP service:

- `food-catalog-service` (Port 8080)
- `order-service` (Port 8081)
- `cafe-ui-service` (Port 80)

These services enable internal service-to-service communication using Kubernetes DNS (servicename.namespace.svc.cluster.local).

---

## 6. API Gateway Configuration (Kong)

### 6.1 Kong Ingress Routes

| Path             | Backend Service      | Port | Purpose          |
| ---------------- | -------------------- | ---- | ---------------- |
| `/api/catalog/*` | food-catalog-service | 8080 | Menu management  |
| `/api/orders/*`  | order-service        | 8081 | Order processing |
| `/*`             | cafe-ui-service      | 80   | Frontend UI      |

### 6.2 Ingress Features

- **Strip Path:** Enabled - removes path prefix before routing
- **Load Balancing:** Automatic across service replicas
- **Request Routing:** URL-based routing to appropriate services
- **External Access:** Provides single entry point for all clients

---

## 7. Data Flow & Communication

### 7.1 Menu Item Retrieval Flow

```
1. Client (Browser) → Kong Gateway: GET /api/catalog/items
2. Kong → Food Catalog Service: GET /items
3. Food Catalog Service → Kong: JSON array of FoodItem
4. Kong → Client: Formatted response
5. Client: Parse and display menu items
```

### 7.2 Order Placement Flow

```
1. Client (Browser) → Kong Gateway: POST /api/orders/orders
2. Kong → Order Service: POST /orders (with item_ids)
3. Order Service:
   a. Validate request format
   b. Discover Food Catalog Service
   c. Generate unique Order ID (UUID)
   d. Store order with status "received"
   e. Return order details
4. Order Service → Kong: Order confirmation
5. Kong → Client: Success message with Order ID
6. Client: Clear cart, display confirmation
```

### 7.3 Service Discovery Process

```
Order Service needs Catalog Service:
1. Call findService("food-catalog-service")
2. Kubernetes DNS resolves to: food-catalog-service.student-cafe.svc.cluster.local
3. Returns: http://food-catalog-service:8080
4. Attempt inter-service communication
```

---

## 8. Deployment Instructions

### 8.1 Prerequisites

- Docker Desktop with Kubernetes enabled OR Minikube
- kubectl CLI tool
- Helm package manager
- 4GB+ available RAM

### 8.2 Step-by-Step Deployment

**Step 1: Initialize Kubernetes**

```bash
minikube start
eval $(minikube docker-env)
```

**Step 2: Create Namespace**

```bash
kubectl create namespace student-cafe
```

**Step 3: Build Docker Images**

```bash
docker build -t food-catalog-service:v1 ./food-catalog-service/
docker build -t order-service:v1 ./order-service/
docker build -t cafe-ui:v1 ./cafe-ui/
```

**Step 4: Install Kong API Gateway**

```bash
helm repo add kong https://charts.konghq.com
helm install kong kong/kong \
  --set ingressController.installCRDs=false \
  --set admin.enabled=true \
  -n student-cafe
```

**Step 5: Deploy Application Services**

```bash
kubectl apply -f app-deployment.yaml -n student-cafe
```

**Step 6: Configure Ingress Routes**

```bash
kubectl apply -f kong-ingress.yaml -n student-cafe
```

**Step 7: Access Application**

```bash
# Get Minikube IP
minikube ip

# Get Kong proxy port
kubectl get service kong-kong-proxy -n student-cafe

# Access at: http://<minikube-ip>:<kong-port>
```

### 8.3 Verification Commands

```bash
# Check deployments
kubectl get deployments -n student-cafe

# Check services
kubectl get services -n student-cafe

# View logs
kubectl logs -f deployment/food-catalog-deployment -n student-cafe
kubectl logs -f deployment/order-deployment -n student-cafe
kubectl logs -f deployment/cafe-ui-deployment -n student-cafe

# Test API endpoints
kubectl exec -it pod/<cafe-ui-pod> -n student-cafe -- curl http://food-catalog-service:8080/items
```

---

## 9. Key Learning Outcomes

### 9.1 Microservices Architecture

- Decomposed monolith into independent services
- Implemented service boundaries
- Managed service-to-service communication
- Designed for scalability

### 9.2 Container Orchestration

- Defined Kubernetes deployments
- Configured service discovery
- Implemented health checks
- Managed service networking

### 9.3 API Gateway Pattern

- Configured Kong for request routing
- Implemented URL-based routing
- Managed cross-service requests
- Provided unified external interface

### 9.4 Cloud-Native Development

- Built containerized applications
- Implemented stateless services
- Used environment variables for configuration
- Designed for horizontal scaling

---

## 10. Challenges & Solutions

| Challenge                   | Solution                                               |
| --------------------------- | ------------------------------------------------------ |
| Service Discovery           | Used Kubernetes built-in DNS resolution                |
| Inter-service Communication | Implemented service lookup with error handling         |
| Local Development           | Used minikube with local Docker registry               |
| Configuration Management    | Environment variables via Kubernetes (optional Consul) |
| API Routing Complexity      | Kong Ingress for transparent routing                   |

---

## 11. Future Enhancements

1. **Database Integration:** Add persistent storage for orders and menu

   - PostgreSQL for relational data
   - Redis for caching

2. **Authentication & Authorization:** Secure order placement

   - JWT token validation
   - Role-based access control

3. **Advanced Service Discovery:** Implement Consul for advanced features

   - Health-aware routing
   - Automatic failover

4. **Monitoring & Logging:**

   - Prometheus for metrics
   - ELK stack for centralized logging
   - Jaeger for distributed tracing

5. **Message Queue Integration:**

   - RabbitMQ/Kafka for asynchronous order processing
   - Event-driven architecture

6. **Horizontal Auto-scaling:**

   - HPA (Horizontal Pod Autoscaler)
   - Based on CPU/memory metrics

7. **Service Mesh:**
   - Istio for advanced networking
   - Circuit breaking and retry logic

---

## 12. Customizations Made

To differentiate from baseline implementations:

1. **Cafe Name:** Changed from "Student Cafe" to "Campus Cafe Express"
2. **Menu Items:** Expanded from 3 to 5 items with specific names and adjusted pricing
   - Espresso ($2.75), Turkey Sandwich ($5.50), Blueberry Muffin ($3.50)
   - Iced Tea ($2.25), Caesar Salad ($6.00)
3. **Order Struct:** Added Timestamp field for order creation tracking
4. **User Feedback:** Updated empty cart message for clarity

---

## 13. Conclusion

This practical successfully demonstrates the implementation of a production-ready microservices architecture using modern cloud-native technologies. The application showcases:

- **Scalability:** Each service can be scaled independently
- **Resilience:** Redundant services and health checks
- **Maintainability:** Clear service boundaries and API contracts
- **Cloud-Native Deployment:** Kubernetes-ready containerized services

The Campus Cafe Express system serves as a foundation for understanding enterprise-grade microservices architecture and cloud computing principles essential for modern software development.

---

## 14. Appendix: File Structure

```
Web303_p4/
├── README.md                           # Project documentation
├── REPORT.md                           # This report
├── app-deployment.yaml                 # Kubernetes deployment manifests
├── kong-ingress.yaml                   # Kong API Gateway configuration
│
├── cafe-ui/                            # React frontend application
│   ├── Dockerfile                      # Container image definition
│   ├── package.json                    # npm dependencies
│   ├── public/
│   │   ├── index.html                  # HTML entry point
│   │   ├── manifest.json               # PWA manifest
│   │   └── robots.txt                  # SEO configuration
│   └── src/
│       ├── App.js                      # Main React component
│       ├── App.css                     # Styling
│       ├── index.js                    # React DOM mount
│       ├── index.css                   # Global styles
│       ├── App.test.js                 # Component tests
│       ├── setupTests.js               # Test configuration
│       └── reportWebVitals.js          # Performance metrics
│
├── food-catalog-service/               # Go microservice
│   ├── Dockerfile                      # Container image definition
│   ├── main.go                         # Service implementation
│   ├── go.mod                          # Go module definition
│   └── go.sum                          # Dependency checksums
│
└── order-service/                      # Go microservice
    ├── Dockerfile                      # Container image definition
    ├── main.go                         # Service implementation
    ├── go.mod                          # Go module definition
    └── go.sum                          # Dependency checksums
```
