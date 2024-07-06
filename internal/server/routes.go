package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"matrix-ai-framework/internal/config"
)

func (s *FiberServer) RegisterFiberRoutes() {

	s.App.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     fmt.Sprintf("https://%s", config.App.Domain),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	api := s.App.Group("/api")

	api.Get("/services", s.GetServices)
	api.Get("/orgs", s.GetOrgs)
	api.Get("/health", s.healthHandler)
	api.Get("/payment", s.GetPaymentState)
	api.Put("/payment", s.PatchUpdatePaymentState)

}
