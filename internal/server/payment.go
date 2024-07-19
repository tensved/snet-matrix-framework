package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tensved/snet-matrix-framework/pkg/db"
)

// GetPaymentState retrieves payment information based on the provided UUID from the query parameter.
// It parses the UUID from the request, retrieves the payment state from the database, and returns it as a JSON response.
//
// Parameters:
//   - c: The Fiber context which provides query parameters and methods to interact with the request and response.
//
// Returns:
//   - error: An error if the operation fails.
func (s *FiberServer) GetPaymentState(c *fiber.Ctx) (err error) {
	paymentUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse uuid")
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse uuid")
	}
	p, err := s.db.GetPaymentState(paymentUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get payment state")
		return c.Status(fiber.StatusBadRequest).SendString("Failed to get payment state")
	}
	return c.JSON(fiber.Map{
		"data": p,
	})
}

// PatchUpdatePaymentState updates specific fields of a payment state based on the provided input.
// It parses the input from the request body, updates the payment state in the database, and returns an appropriate response.
//
// Parameters:
//   - c: The Fiber context which provides methods to interact with the request and response.
//
// Returns:
//   - error: An error if the operation fails.
func (s *FiberServer) PatchUpdatePaymentState(c *fiber.Ctx) error {
	var input db.PaymentState
	if err := c.BodyParser(&input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
	}

	if err := s.db.PatchUpdatePaymentState(&input); err != nil {
		log.Error().Err(err).Msg("Failed to update payment state")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update payment state")
	}

	return c.SendStatus(fiber.StatusOK)
}
