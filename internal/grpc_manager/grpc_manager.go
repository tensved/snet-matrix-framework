package grpc_manager

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"
)

// GRPCClientManager manages a pool of gRPC service connections.
type GRPCClientManager struct {
	clients map[string]*GRPCService // A map to store gRPC service connections by target address.
	mu      sync.Mutex              // A mutex to ensure thread-safe access to the clients map.
}

// NewGRPCClientManager creates and returns a new GRPCClientManager.
func NewGRPCClientManager() *GRPCClientManager {
	return &GRPCClientManager{
		clients: make(map[string]*GRPCService), // Initializes the clients map.
	}
}

// GRPCService represents a gRPC service with a connection to a target address.
type GRPCService struct {
	Target      string           // The target address of the gRPC service.
	DialOptions grpc.DialOption  // The options for dialing the gRPC service.
	Conn        *grpc.ClientConn // The client connection to the gRPC service.
}

// GetClient provides a managed connection to a gRPC service.
// If a connection to the target already exists and is healthy, it returns the existing connection.
// Otherwise, it creates a new connection.
func (manager *GRPCClientManager) GetClient(target string) (*GRPCService, error) {
	manager.mu.Lock() // Locks the mutex to ensure thread-safe access.
	defer manager.mu.Unlock()

	if client, exists := manager.clients[target]; exists {
		// Check if the existing connection is still healthy.
		if client.Conn.GetState() == connectivity.Ready {
			return client, nil
		}
		// Close and delete the stale connection.
		client.Close()
		delete(manager.clients, target)
	}

	// Create a new client if none exists or the existing one was stale.
	newClient, err := NewGRPCService(target)
	if err != nil {
		return nil, err
	}
	manager.clients[target] = newClient // Store the new client in the map.
	return newClient, nil
}

// NewGRPCService creates and returns a new GRPCService for the specified target address.
func NewGRPCService(target string) (*GRPCService, error) {
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials()) // Use insecure credentials for the connection.
	conn, err := grpc.Dial(target, dialOptions)                             // Dial the target address.
	if err != nil {
		return nil, err
	}
	return &GRPCService{
		Target:      target,
		DialOptions: dialOptions,
		Conn:        conn,
	}, nil
}

// Close closes the GRPCService connection.
func (s *GRPCService) Close() {
	if s.Conn != nil {
		err := s.Conn.Close() // Close the connection.
		if err != nil {
			log.Error().Err(err).Msg("Failed to close connection") // Log an error if closing fails.
		}
	}
}

// CallMethod invokes a gRPC method on the service.
// It creates a new context with a timeout and attaches metadata to the context.
func (s *GRPCService) CallMethod(method string, req, resp interface{}, md metadata.MD) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10) // Create a context with a timeout.
	defer cancel()                                                           // Ensure the context is cancelled to free resources.
	ctx = metadata.NewOutgoingContext(ctx, md)                               // Attach metadata to the context.
	return s.Conn.Invoke(ctx, method, req, resp)                             // Invoke the gRPC method with the context, request, and response.
}
