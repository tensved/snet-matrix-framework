package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (s *FiberServer) GetServices(c fiber.Ctx) error {
	services, err := s.db.GetSnetServices()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get services")
	}
	return c.JSON(services)
}

func (s *FiberServer) GetOrgs(c fiber.Ctx) error {
	orgs, err := s.db.GetSnetOrgs()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get orgs")
	}
	return c.JSON(orgs)
}
