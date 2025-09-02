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

type UserLoginReq struct {
	SessionId             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	UserId                int32     `json:"user_id"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
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
// @Success		200 {object} UserLoginReq "Login successful, returns JWT token"
// @Failure		204 {string} string "Invalid credentials for user login"
// @Failure		500 {string} string "Internal server error"
// @Router			/login [post]
func (s *FiberServer) Login(ctx *fiber.Ctx) error {

	auth := models.Auth{}

	err := ctx.BodyParser(&auth)

	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error: %s", err.Error()))
	}

	userDetails, err := s.db.DB().UserLogin(s.ctx, auth.UserName)
	if err != nil {
		return ctx.Status(http.StatusNotFound).SendString(fmt.Sprintf("Error %s", err.Error()))
	}

	passwordMatch := authentication.ComparePassword(userDetails.Password, auth.Password)

	if !passwordMatch {
		return ctx.Status(fiber.StatusNoContent).SendString("Invalid credentials for user login")
	}

	var userRole = string(userDetails.UserRole)
	token, claims, err := s.token.CreateToken(userDetails.ID, userRole, time.Minute*15)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error %s", err.Error()))
	}

	refreshToken, refreshClaims, err := s.token.CreateToken(userDetails.ID, userRole, time.Hour*9)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error %s", err.Error()))
	}

	sn := models.Session{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserId:       userDetails.ID,
		UserRole:     userRole,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	session := models.CreateSession(sn)
	s.db.DB().InsertSession(s.ctx, session)

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Minute * 15),
	})

	return ctx.Status(201).JSON(
		UserLoginReq{
			SessionId:             session.ID,
			AccessToken:           token,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  claims.RegisteredClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
			UserId:                session.UserID,
		})
}

// @Summary		When a logged in user want to logout
// @Description Logs out the user that is currently logged in
// @Tags		Logout
// @Produce 	json
// @Param 		logout path int true "Path Variable"
// @Success 	204 {string} string "No Content"
// @Failure 	400 {string} string "Bad Request: Missing Id"
// @Failure 	500 {string} string "Internal Server Error: when deleting session"
// @Router 	/logout [delete]
func (s *FiberServer) Logout(ctx *fiber.Ctx) error {

	sessionId := ctx.Params("id", "")

	if sessionId == "" {

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("missing session id %v", fiber.ErrBadRequest),
		})
	}

	if err := s.db.DB().DeleteSession(s.ctx, sessionId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Err deleting the session%v", fiber.ErrInternalServerError),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary		Renew user token
// @Description Logged in users can use to get a new session token
// @Produce 	json
// @Tags		RenewToken
// @Param 		RenewToken body RenewAccessTokenReq true "access token request"
// @Success 	200 {object} RenewAccessTokenRes "Success"
// @Failure 	400 {string} string "Bad Request"
// @Failure 	401 {string} string "Un Authorized"
// @Failure 	500 {string} string "Internal Server Error: when deleting session"
// @Router 	/renew/token [post]
func (s *FiberServer) RenewAccessToken(ctx *fiber.Ctx) error {

	refreshToken := RenewAccessTokenReq{}

	err := ctx.BodyParser(&refreshToken)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("error decoding token: %s", err.Error()))
	}

	refreshClaims, err := s.token.VerifyToken(refreshToken.RefreshToken)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("error verifying token: %s", err.Error()))
	}

	session, err := s.db.DB().GetSession(s.ctx, refreshClaims.RegisteredClaims.ID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error getting session: %s", err.Error()))
	}

	if session.UserID != refreshClaims.UserId {
		return ctx.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Invalid session: %s", err.Error()))
	}

	token, claims, err := s.token.CreateToken(refreshClaims.UserId, refreshClaims.UserRole, 15*time.Minute)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error creating token: %s", err.Error()))
	}

	res := RenewAccessTokenRes{
		RefreshToken:         token,
		AccessTokenExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// @Summary		When a logged in user want to revoke their session
// @Description Logs out the user that is currently logged in
// @Tags		RevokeSession
// @Produce 	json
// @Param 		RevokeSession path int true "Path Variable"
// @Success 	204 {string} string "No Content"
// @Failure 	400 {string} string "Bad Request: Missing Id"
// @Failure 	500 {string} string "Internal Server Error: when deleting session"
// @Router 		/revoke/token [delete]
func (s *FiberServer) RevokeSession(ctx *fiber.Ctx) error {

	sessionId := ctx.Params("id", "")

	if sessionId == "" {

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("missing session id %v", fiber.ErrBadRequest),
		})
	}

	if err := s.db.DB().RevokeSession(s.ctx, sessionId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Err revoking the session%v", fiber.ErrInternalServerError),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// change and reset password
