package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// GetServices handles the endpoint for retrieving a list of services.
// It fetches the services from the database and returns them as a JSON response.
//
// Parameters:
//   - c: The Fiber context which provides methods to interact with the request and response.
//
// Returns:
//   - error: An error if the operation fails or nil if the operation is successful.
func (s *FiberServer) GetServices(c *fiber.Ctx) error {
	services, err := s.db.GetSnetServices()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get services")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve services")
	}
	return c.JSON(services)
}

// GetOrgs handles the endpoint for retrieving a list of organizations.
// It fetches the organizations from the database and returns them as a JSON response.
//
// Parameters:
//   - c: The Fiber context which provides methods to interact with the request and response.
//
// Returns:
//   - error: An error if the operation fails or nil if the operation is successful.
func (s *FiberServer) GetOrgs(c *fiber.Ctx) error {
	orgs, err := s.db.GetSnetOrgs()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get orgs")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve organizations")
	}
	return c.JSON(orgs)
}
