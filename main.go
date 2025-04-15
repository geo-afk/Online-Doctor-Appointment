package main

import (
	"context"
	"log"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/geo-afk/Online-Doctor-Appointment/auth/db"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault) // get one user by ID
}

func databaseInit() {

	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)

	}
	defer conn.Close(ctx)

	queries := db.New(conn)
}

func main() {
	databaseInit()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	SwaggerRoute(app)

	patientRoute(app)
	doctorRoute(app)

	err := app.Listen(":3400")
	if err != nil {
		log.Fatal("Error Listening: ", err.Error())
	}

}

func patientRoute(app *fiber.App) {

	patient := app.Group("/patient")
	patient.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("patient")
	})
}

func doctorRoute(app *fiber.App) {

	doctor := app.Group("/doctor")
	doctor.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("doctor")
	})
}
