package routes

import (
	"errors"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func PasswordReset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "empty_body", Message: "Nothing in body"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if len(data["email"]) <= 0 {
		rp := models.ResponsePacket{Error: true, Code: "missing_data", Message: "Form is missing required data!"}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !controller.EmailIsValid(data["email"]) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Email is not valid."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "User not found."}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	err := controller.SendPasswordResetEmail(user.Email)
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not send email."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"message": "Email sent!"})
}
