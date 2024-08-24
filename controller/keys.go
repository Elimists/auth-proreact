package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InvalidateJWT(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     fmt.Sprintf("%s_jwt", os.Getenv("API_NAME")),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		SameSite: "Lax",
	})

	return nil
}
