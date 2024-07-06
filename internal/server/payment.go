package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"matrix-ai-framework/pkg/db"
)

// GetPaymentState retrieves a payment info
func (s *FiberServer) GetPaymentState(c *fiber.Ctx) (err error) {
	paymentUUID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse uuid")
		return
	}
	p, err := s.db.GetPaymentState(paymentUUID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get payment state")
		return
	}
	return c.JSON(fiber.Map{
		"data": p,
	})
}

// PatchUpdatePaymentState patch updates a payment state based on provided fields
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
