package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tensved/snet-matrix-framework/pkg/db"
)

// FiberServer represents the server that uses the Fiber web framework.
type FiberServer struct {
	App *fiber.App // Embeds the Fiber application instance.
	db  db.Service // Database service for the server.
}

// New creates and returns a new instance of FiberServer.
// It initializes the Fiber application and sets up the database service.
//
// Parameters:
//   - db: An instance of db.Service which provides database-related functionalities.
//
// Returns:
//   - A pointer to the initialized FiberServer instance.
func New(db db.Service) *FiberServer {
	server := &FiberServer{
		App: fiber.New(), // Initializes the Fiber application.
		db:  db,          // Sets the provided database service.
	}

	return server
}
