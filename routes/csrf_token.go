package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
)

func CSRFToken(c *fiber.Ctx) error {
	csrfToken, err := controller.GenerateCSRFToken()
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not generate CSRF token."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	c.Cookie(&fiber.Cookie{
		Name:     fmt.Sprintf("%s_csrf", os.Getenv("API_NAME")),
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		fmt.Sprintf("%s_csrf", os.Getenv("API_NAME")): csrfToken,
	})
}
