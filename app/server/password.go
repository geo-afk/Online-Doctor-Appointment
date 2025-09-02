package server

import (
	"fmt"
	"time"

	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
)

// Request structures
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// Request structures
type ResetPasswordRequest struct {
	Pass string `json:"password"`
}

func (s *FiberServer) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	userFound, err := s.db.DB().UserByEmail(s.ctx, models.ToPgText(req.Email))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintln("something went wrong: ", err.Error()),
		})
	}

	if !userFound {
		// Don't reveal whether email exists or not (security best practice)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "If your email exists in our system, you will receive a password reset link",
		})
	}

	// Generate token
	token, err := GenerateSecureToken(32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate reset token",
		})
	}

	duration := time.Now().Add(time.Minute * 3)
	request := models.CreateCreateRequestParams("password_reset", token, req.Email, duration)
	s.db.DB().CreateRequest(s.ctx, request)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If your email exists in our system, you will receive a password reset link",
	})
}

func (s *FiberServer) RecoverPassword(ctx *fiber.Ctx) error {

	passRequest := ResetPasswordRequest{}
	ctx.BodyParser(&passRequest)

	token := ctx.Params("token", "")

	if token == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	request, err := s.db.DB().GetRequest(s.ctx, token)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	expired := request.ExpiresAt.Time

	time.Now().After(expired)
	if time.Now().After(expired) {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	passwordParam := models.CreateForgetPasswordParams(passRequest.Pass, request.Token)
	passwordParam.Password = GeneratePassword(passwordParam.Password)
	err = s.db.DB().ForgetPassword(s.ctx, passwordParam)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s *FiberServer) ChangePassword(ctx *fiber.Ctx) error {

	token := ctx.Cookies("jwt", "")

	if token == "" {
		return ctx.Status(fiber.StatusNoContent).SendString("Invalid token")
	}

	claims, err := s.token.VerifyToken(token)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error Parsing token %v", err.Error()))
	}

	passRequest := ResetPasswordRequest{}
	err = ctx.BodyParser(&passRequest)

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	user_id, ok := ctx.Locals("user_id").(int32)

	if !ok {
		ctx.SendStatus(fiber.StatusNoContent)
	}

	if claims.UserId != user_id {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	passwordParam := models.CreateChangePasswordParams(passRequest.Pass, user_id)
	passwordParam.Password = GeneratePassword(passwordParam.Password)
	err = s.db.DB().ChangePassword(s.ctx, passwordParam)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": fmt.Sprintf("Error changing password %v", err.Error()),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)

}
