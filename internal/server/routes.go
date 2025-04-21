package server

import (
	"context"
	"fmt"
	"log"

	"github.com/geo-afk/Online-Doctor-Appointment/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/geo-afk/Online-Doctor-Appointment/cmd/docs"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Added Logging to the routes
	s.Use(logger.New())

	// Add Swagger Route for code documentation
	s.Group("/docs").Get("*", swagger.HandlerDefault)

	s.App.Get("/health", s.healthHandler)
	s.App.Get("/login", s.Login)

}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func HandleError[T any](value T, err error) T {

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	return value
}

func (s *FiberServer) Login(ctx *fiber.Ctx) error {

	auth := utils.Auth{}

	err := ctx.BodyParser(&auth)
	if err != nil {
		return ctx.Status(500).SendString(fmt.Sprintf("Error: %s", err.Error()))
	}

	userLogin := utils.AuthToUserLoginParams(auth)

	userDetails := HandleError(s.db.DB().UserLogin(context.Background(), userLogin))
	auth.UserID = userDetails.ID

	token := HandleError(utils.GenerateToken(uint(userDetails.ID)))

	return ctx.Status(201).JSON(
		fiber.Map{
			"user_role": userDetails.UserRole,
			"token":     token,
		})
}
