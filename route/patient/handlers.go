package patient

import "github.com/gofiber/fiber/v2"

func GetCustomer(ctx *fiber.Ctx) error {

	return ctx.SendString("patient")
}
