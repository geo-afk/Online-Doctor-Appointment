package server

import (
	db "github.com/geo-afk/Online-Doctor-Appointment/internal/db"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	db db.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "Appointment",
			AppName:      "Appointment",
		}),

		db: db.New(),
	}

	return server
}
