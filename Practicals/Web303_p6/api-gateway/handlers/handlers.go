package handlers

import (
	"net/http"

	"api-gateway/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handlers manages HTTP request handling with backend gRPC service connections
type Handlers struct {
	clients *grpc.ServiceClients
}

// NewHandlers initializes a new request handler with service client connections
func NewHandlers(clients *grpc.ServiceClients) *Handlers {
	return &Handlers{clients: clients}
}

// handleGRPCError translates gRPC status codes to HTTP error responses
func handleGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, return generic internal server error
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Map gRPC status codes to HTTP status codes
	var httpStatus int
	switch st.Code() {
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.FailedPrecondition:
		httpStatus = http.StatusPreconditionFailed
	case codes.Unimplemented:
		httpStatus = http.StatusNotImplemented
	case codes.Unavailable:
		httpStatus = http.StatusServiceUnavailable
	default:
		httpStatus = http.StatusInternalServerError
	}

	http.Error(w, st.Message(), httpStatus)
}
