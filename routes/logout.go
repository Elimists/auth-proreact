package routes

import (
	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {

	controller.InvalidateCSRFToken(c)
	controller.InvalidateJWT(c)

	rp := models.ResponsePacket{Error: false, Code: "successfull", Message: "Logged out successfully."}
	return c.Status(fiber.StatusOK).JSON(rp)
}
