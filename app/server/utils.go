package server

import (
	"log"

	"github.com/geo-afk/Online-Doctor-Appointment/app/auth"
	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
)

func ParseAndRegisterUser(ctx *fiber.Ctx) models.User {

	var user models.User

	ctx.BodyParser(&user)

	log.Println(user)

	password := auth.GeneratePassword(user.Auth.Password)
	user.Auth.Password = password

	return user

}
