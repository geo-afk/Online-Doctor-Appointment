package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	authentication "github.com/geo-afk/Online-Doctor-Appointment/app/auth"
	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
)

type tokenResponse struct {
	UserRole any    `json:"user_role"`
	Token    string `json:"token"`
}

// @Summary      Register a new user (Doctor or Patient)
// @Description  Registers a new user, handling both doctor and patient profiles.
// @Tags         Register
// @Accept       json
// @Produce      plain
// @Param		 Register body models.User true "User Details"
// @Success      200 {boolean} bool "Successfully registered"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /register [post]
func (s *FiberServer) Register(ctx *fiber.Ctx) error {

	bgContext := context.Background()
	user := ParseAndRegisterUser(ctx)

	fmt.Printf("%#v", user)
	fmt.Printf("\n\n%#v", *user.Contact)

	primaryContact := (*user.Contact).ToBdType()
	primaryContactId, err := s.db.DB().CreateContact(bgContext, primaryContact)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error occurred: %s", err.Error()))
	}

	emergencyContact := (*user.EmergencyContact).ToBdType()
	emergencyContactId, err := s.db.DB().CreateContact(bgContext, emergencyContact)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error occurred: %s", err.Error()))
	}

	userParam := user.ToBdType(primaryContactId, emergencyContactId)
	userId, err := s.db.DB().RegisterUser(bgContext, userParam)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error occurred: %s", err.Error()))
	}

	auth := user.Auth.ToBdType(userId)
	err = s.db.DB().CreateUserAuth(bgContext, auth)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error occurred: %s", err.Error()))
	}

	return nil
}

// @Summary		Login Panel for both Doctor and Patient
// @Tags			Login
// @Accept			json
// @Produce		json
// @Param			login body models.Auth true "Login credentials"
// @Success		200 {object} tokenResponse "Login successful, returns JWT token"
// @Failure		204 {string} string "Invalid credentials for user login"
// @Failure		500 {string} string "Internal server error"
// @Router			/login [post]
func (s *FiberServer) Login(ctx *fiber.Ctx) error {
	auth := models.Auth{}

	err := ctx.BodyParser(&auth)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error: %s", err.Error()))
	}

	userDetails, err := s.db.DB().UserLogin(context.Background(), auth.UserName)
	if err != nil {
		return ctx.Status(http.StatusNotFound).SendString(fmt.Sprintf("Error %s", err.Error()))
	}

	passwordMatch := authentication.ComparePassword(userDetails.Password, auth.Password)

	if !passwordMatch {
		return ctx.Status(fiber.StatusNoContent).SendString("Invalid credentials for user login")
	}

	duration := time.Minute * 30
	token, _, err := s.token.CreateToken(auth.UserName, string(userDetails.UserRole), duration)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error %s", err.Error()))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		MaxAge:   60 * 30,
	})

	return ctx.Status(201).JSON(
		tokenResponse{
			UserRole: userDetails.UserRole,
			Token:    token,
		})
}

// ---@@@Security ApiKeyAuth
