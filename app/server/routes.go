package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"

	_ "github.com/geo-afk/Online-Doctor-Appointment/docs"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Use(cache.New())

	s.App.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	// Initialize default config for rate limiter
	s.App.Use(limiter.New())

	// Added Logging to the routes
	s.Use(logger.New())

	// Add Swagger Route for code documentation
	s.Group("/docs").Get("*", swagger.HandlerDefault)

	s.App.Get("/health", s.healthHandler)
	s.App.Get("/metrics", monitor.New())

	s.App.Post("/login", s.Login)
	// s.App.Post("/reset", s.ResetPassword)
	s.App.Post("/register", s.Register)

	v1 := s.App.Group("/api/v1")
	protected := v1.Use(s.AuthMiddleware)
	protected.Post("/book_appointment", s.BookAppointment)
}

// @Summary		Get the Database health
// @Description Get the Database health
// @Tags		Health
// @Produce		json
// @Success		200	{map}	map[string]string
// @Router		/health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
