package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tensved/snet-matrix-framework/internal/config"
)

// RegisterFiberRoutes sets up the Fiber routes and middleware for the server.
func (s *FiberServer) RegisterFiberRoutes() {

	// Use CORS middleware with the specified configuration.
	s.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin", // Specifies the headers allowed in CORS requests.
		AllowOrigins:     fmt.Sprintf("https://%s", config.App.Domain),                                                                              // Specifies the allowed origins for CORS requests, formatted with the application domain.
		AllowCredentials: true,                                                                                                                      // Indicates whether the request can include user credentials.
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",                                                                                  // Specifies the allowed methods for CORS requests.
	}))

	// Create a new group for API routes.
	api := s.App.Group("/api")

	// Register the route handlers.
	api.Get("/services", s.GetServices)            // Retrieves a list of services.
	api.Get("/orgs", s.GetOrgs)                    // Retrieves a list of organizations.
	api.Get("/health", s.healthHandler)            // Checks the health of the server.
	api.Get("/payment", s.GetPaymentState)         // Retrieves a payment state.
	api.Put("/payment", s.PatchUpdatePaymentState) // Updates a payment state based on provided fields.
}
