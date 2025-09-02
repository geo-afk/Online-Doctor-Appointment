package server

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/geo-afk/Online-Doctor-Appointment/app/auth"
	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
)

func ParseAndRegisterUser(ctx *fiber.Ctx) models.User {

	var user models.User

	ctx.BodyParser(&user)

	log.Println(user)

	password := GeneratePassword(user.Auth.Password)
	user.Auth.Password = password

	return user

}


func GeneratePassword(password string) string{

	return auth.GeneratePassword(password)
}

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
