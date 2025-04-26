package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	_ "github.com/geo-afk/Online-Doctor-Appointment/docs"
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

	// protected := s.Use(s.AuthMiddleware)

	// Add Swagger Route for code documentation
	s.Group("/docs").Get("*", swagger.HandlerDefault)

	s.App.Get("/health", s.healthHandler)
	s.App.Post("/login", s.Login)
	s.App.Post("/register", s.Register)
}

// @Summary	get the Database health
// @ID			health
// @Produce	json
// @Success	200	{object}	map[string]string
// @Router		/health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
