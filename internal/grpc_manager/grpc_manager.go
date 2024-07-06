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

// GRPCClientManager manages a pool of GRPCService connections.
type GRPCClientManager struct {
	clients map[string]*GRPCService
	mu      sync.Mutex
}

// NewGRPCClientManager creates a new GRPCClientManager.
func NewGRPCClientManager() *GRPCClientManager {
	return &GRPCClientManager{
		clients: make(map[string]*GRPCService),
	}
}

type GRPCService struct {
	Target      string
	DialOptions grpc.DialOption
	Conn        *grpc.ClientConn
}

// GetClient provides a managed connection to a gRPC service.
func (manager *GRPCClientManager) GetClient(target string) (*GRPCService, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if client, exists := manager.clients[target]; exists {
		// Check if connection is still alive
		if client.Conn.GetState() == connectivity.Ready {
			return client, nil
		}
		// Close and delete stale connection
		client.Close()
		delete(manager.clients, target)
	}

	// Create new client if not existing or deleted
	newClient, err := NewGRPCService(target)
	if err != nil {
		return nil, err
	}
	manager.clients[target] = newClient
	return newClient, nil
}

// NewGRPCService creates and returns a new GRPCService.
func NewGRPCService(target string) (*GRPCService, error) {
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(target, dialOptions)
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
		err := s.Conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close connection")
		}
	}
}

// CallMethod invoke grpc method
func (s *GRPCService) CallMethod(method string, req, resp interface{}, md metadata.MD) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, md)
	return s.Conn.Invoke(ctx, method, req, resp)
}
