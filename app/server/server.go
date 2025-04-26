package server

import (
	jwt "github.com/geo-afk/Online-Doctor-Appointment/app/auth"
	db "github.com/geo-afk/Online-Doctor-Appointment/app/db"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	db    db.Service
	token *jwt.JwtSecret
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "Appointment",
			AppName:      "Appointment",
		}),

		db:    db.New(),
		token: jwt.NewJwt(),
	}

	return server
}
