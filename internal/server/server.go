package server

import (
	"github.com/gofiber/fiber/v2"
	"matrix-ai-framework/pkg/db"
)

type FiberServer struct {
	*fiber.App
	db db.Service
}

func New(db db.Service) *FiberServer {
	server := &FiberServer{
		App: fiber.New(),
		db:  db,
	}

	return server
}
