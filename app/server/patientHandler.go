package server

import (
	"context"

	"github.com/geo-afk/Online-Doctor-Appointment/app/models"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Appointment ;patients come and book appointments
// @Tags		BookAppointment
// @Security 	bearerToken
// @Accept		json
// @Produce		json
// @Param		appointment body models.Appointment true "Book Appointment"
// @Success		200 {boolean} true "Login successfully booked appointment"
// @Failure		204 {string} string "Unable to book appointment"
// @Failure		500 {string} string "Internal server error"
// @Router		/api/v1/book_appointment [post]
func (s *FiberServer) BookAppointment(ctx *fiber.Ctx) error {

	user_id, ok := ctx.Locals("user_id").(int32)

	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not find credentials",
		})
	}

	appointment := models.Appointment{}

	err := ctx.BodyParser(&appointment)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	appointment.UserID = user_id
	appointmentParam := appointment.GetBookAppointmentParams()
	err = s.db.DB().BookAppointment(context.Background(), appointmentParam)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).SendString("appointment booked")
}
