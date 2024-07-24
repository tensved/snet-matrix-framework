package server

import "github.com/gofiber/fiber/v2"

// healthHandler handles the health check endpoint.
// It retrieves the health status from the database and returns it as a JSON response.
//
// Parameters:
//   - c: The Fiber context which provides methods to interact with the request and response.
//
// Returns:
//   - error: An error if the operation fails.
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
