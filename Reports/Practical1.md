# Practical One: gRPC Microservices with Docker

## Overview

This practical demonstrates building and deploying a microservices architecture using **gRPC** (gRPC Remote Procedure Call) and **Docker**. The project consists of two independent services that communicate with each other over gRPC.

## Project Structure

```
practical-one/
├── docker-compose.yml          # Orchestrates multi-container setup
├── greeter-service/
│   ├── Dockerfile             # Container configuration for greeter service
│   ├── main.go                # Greeter service implementation
│   └── go.mod                 # Go module dependencies
├── time-service/
│   ├── Dockerfile             # Container configuration for time service
│   ├── main.go                # Time service implementation
│   └── go.mod                 # Go module dependencies
└── proto/
    ├── greeter.proto          # gRPC service definition for greeter
    ├── time.proto             # gRPC service definition for time
    └── gen/
        └── proto/
            ├── greeter.pb.go        # Generated protobuf code
            ├── greeter_grpc.pb.go   # Generated gRPC code for greeter
            ├── time.pb.go           # Generated protobuf code
            └── time_grpc.pb.go      # Generated gRPC code for time
```

## Services

### 1. **Time Service** (Port 50052)

- **Purpose**: Provides current time in RFC3339 format
- **Endpoint**: Listens on `time-service:50052`
- **RPC Method**: `GetTime(TimeRequest) → TimeResponse`
- **Returns**: Current system time

### 2. **Greeter Service** (Port 50051)

- **Purpose**: Greets users with a personalized message that includes the current time
- **Endpoint**: Listens on `greeter-service:50051` (exposed to host as port 50051)
- **RPC Method**: `SayHello(HelloRequest) → HelloResponse`
- **Functionality**:
  - Accepts a name as input
  - Calls the Time Service to get current time
  - Returns: `"Hello {name}! The current time is {timestamp}"`

## Communication Flow

```
User/Client
    ↓
Greeter Service (Port 50051)
    ↓
[Calls via gRPC]
    ↓
Time Service (Port 50052)
    ↓
Returns current time
    ↓
Greeter Service formats response
    ↓
Response sent back to client
```

## Protocol Buffer Definitions

### greeter.proto

```proto3
service GreeterService {
  rpc SayHello(HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

### time.proto

```proto3
service TimeService {
  rpc GetTime(TimeRequest) returns (TimeResponse);
}

message TimeRequest {}

message TimeResponse {
  string current_time = 1;
}
```

## Key Implementation Details

### Greeter Service (`greeter-service/main.go`)

- Establishes a gRPC client connection to `time-service:50052`
- Implements `SayHello` RPC method
- Calls the Time Service to fetch current time
- Combines greeting with time information

### Time Service (`time-service/main.go`)

- Implements `GetTime` RPC method
- Returns current system time in RFC3339 format
- Runs independently and can be called by other services

### Docker Composition (`docker-compose.yml`)

- **Services**:
  - `time-service`: Backend service (no exposed ports)
  - `greeter-service`: Frontend service (exposes port 50051)
- **Networking**: Docker automatically creates a network allowing service-to-service communication via hostname
- **Dependencies**: Greeter service depends on Time service startup

## How to Run

### Prerequisites

- Docker
- Docker Compose

### Steps

1. **Build and start services**:

```bash
docker-compose up --build
```

2. **Services will start in order**:

   - Time Service starts first (no dependencies)
   - Greeter Service starts and connects to Time Service

3. **Test the Greeter Service** (from another terminal):

```bash
# Example using grpcurl (if installed)
grpcurl -plaintext -d '{"name":"World"}' localhost:50051 greeter.GreeterService/SayHello
```

4. **Stop services**:

```bash
docker-compose down
```

## Learning Outcomes

This practical demonstrates:

**gRPC Protocol**: Building services using Protocol Buffers and gRPC  
**Service-to-Service Communication**: Inter-service communication over gRPC  
**Docker Containerization**: Packaging Go applications in containers  
**Docker Compose**: Orchestrating multi-container applications  
**Service Dependencies**: Managing startup order and inter-service dependencies  
**Microservices Architecture**: Building loosely-coupled, independently deployable services

## Technologies Used

- **Language**: Go (Golang)
- **RPC Framework**: gRPC
- **Serialization**: Protocol Buffers (Proto3)
- **Containerization**: Docker
- **Orchestration**: Docker Compose

## Logs

When running with `docker-compose up`, you can monitor service logs:

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f time-service
docker-compose logs -f greeter-service
```

## Troubleshooting

| Issue                | Solution                                                                        |
| -------------------- | ------------------------------------------------------------------------------- |
| "Connection refused" | Ensure Time Service has fully started before Greeter Service tries to connect   |
| "Service not found"  | Check Docker Compose networking; services should use service names as hostnames |
| Port conflicts       | Change port mappings in docker-compose.yml if ports 50051 or 50052 are in use   |
