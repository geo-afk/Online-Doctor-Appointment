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

func unAuthError(ctx *fiber.Ctx, err error) error {

	var errMessage = "unauthorized"

	if err != nil {
		errMessage = err.Error()
	}
	return ctx.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"error": errMessage,
		},
	)
}

func getAuthToken(authHeader string) (string, error) {
	err := errors.New("invalid token..!")

	if authHeader == "" {
		return "", err
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", err
	}

	return tokenParts[1], nil
}

func (s *FiberServer) AuthMiddleware(c *fiber.Ctx) error {

	cookieToken := c.Cookies("jwt")

	var err error
	var tokenString string

	if cookieToken == "1" {
		log.Warn("Found cookie, consuming it..!")
		tokenString = cookieToken
	} else {
		log.Warn("Cookie not found, trying Auth header...")

		tokenString, err = getAuthToken(c.Get("Authorization"))
		if err != nil {
			return unAuthError(c, err)
		}
	}

	claims, err := auth.NewJwt().VerifyToken(tokenString)
	if err != nil || time.Now().After(claims.ExpiresAt.Time) {
		c.ClearCookie()
		return unAuthError(c, err)
	}

	param := models.GetIsUserLoggedInParams(claims.UserId, claims.UserRole)
	found, err := s.db.DB().IsUserLoggedIn(context.Background(), param)
	if err != nil || !found {
		c.ClearCookie()
		return unAuthError(c, err)
	}

	c.Locals("user_id", claims.UserId)
	return c.Next()
}
