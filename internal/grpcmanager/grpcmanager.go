package grpcmanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net/url"
	"sync"
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
func (manager *GRPCClientManager) GetClient(serviceURL string) (*GRPCService, error) {
	manager.mu.Lock() // Locks the mutex to ensure thread-safe access.
	defer manager.mu.Unlock()

	if client, exists := manager.clients[serviceURL]; exists {
		if client.Conn == nil {
			return nil, fmt.Errorf("grpc client for service %s is not available", serviceURL)
		}
		// Check if the existing connection is still healthy.
		if client.Conn.GetState() == connectivity.Ready {
			return client, nil
		}
		// Close and delete the stale connection.
		client.Close()
		delete(manager.clients, serviceURL)
	}

	// Create a new client if none exists or the existing one was stale.
	newClient, err := NewGRPCService(serviceURL)
	if err != nil {
		return nil, err
	}
	manager.clients[serviceURL] = newClient // Store the new client in the map.
	return newClient, nil
}

// NewGRPCService creates and returns a new GRPCService for the specified target address.
func NewGRPCService(serviceURL string) (*GRPCService, error) {
	parsedURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}
	if parsedURL.Hostname() == "" || parsedURL.Port() == "" {
		log.Error().Msgf("invalid service URL: %s", serviceURL)
		return nil, fmt.Errorf("invalid service URL: %s", serviceURL)
	}
	target := parsedURL.Hostname() + ":" + parsedURL.Port()

	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials()) // Use insecure credentials for the connection.
	if parsedURL.Scheme == "https" {
		dialOptions = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")) // Use tls credentials for the connection.
	}
	callOpts := grpc.WithDefaultCallOptions(
		grpc.MaxCallSendMsgSize(1024*1024*16), // 16 MB send size limit
		grpc.MaxCallRecvMsgSize(1024*1024*16), // 16 MB receive size limit
	)

	conn, err := grpc.NewClient(target, dialOptions, callOpts) // Dial the target address.
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("failed to create grpc connection")
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
			log.Error().Err(err).Msg("failed to close connection")
		}
	}
}

// CallMethod invokes a gRPC method on the service.
// It creates a new context with a timeout and attaches metadata to the context.
func (s *GRPCService) CallMethod(method string, req, resp interface{}, md metadata.MD) error {
	err := s.Conn.Invoke(metadata.NewOutgoingContext(context.Background(), md), method, req, resp)
	if err != nil {
		log.Error().Err(err).Msg("failed to invoke method")
		return err
	}
	return nil
}
