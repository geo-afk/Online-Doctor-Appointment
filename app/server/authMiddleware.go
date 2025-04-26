package server

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/geo-afk/Online-Doctor-Appointment/app/auth"
	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *FiberServer) AuthMiddleware() fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		cookieToken := ctx.Cookies("jwt")
		var err error
		var cookieString string

		if cookieToken != "" {
			log.Warn("Found cookie, consuming it..!")
			cookieString = cookieToken
		} else {
			log.Warn("Found Not Found, trying Auth header")

			authHeader := ctx.Cookies("Authorization")

			cookieString, err = getAuthToken(authHeader)

			if err != nil {
				return unAuthError(ctx)
			}

		}
		claims, err := auth.NewJwt().VerifyToken(cookieString)

		expired := time.Now().Before(claims.ExpiresAt.Time)
		if err != nil || expired {
			ctx.ClearCookie()
			return unAuthError(ctx)
		}

		param := models.GetIsUserLoggedInParams(claims.UserName, claims.UserRole)
		found, err := s.db.DB().IsUserLoggedIn(context.Background(), param)

		if found {
			ctx.Locals("user_name", claims.UserName)
		} else {
			ctx.ClearCookie()
			return ctx.SendStatus(fiber.StatusNoContent)
		}
		return ctx.Next()
	}
}

func getAuthToken(authHeader string) (string, error) {

	err := errors.New("unauthorized")
	if authHeader == "" {

		return "", err

	}

	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {

		return "", err
	}

	return tokenParts[1], nil
}

func unAuthError(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"error": "",
		})
}
