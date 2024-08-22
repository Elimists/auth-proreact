package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
)

func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func ValidateCSRFToken(c *fiber.Ctx) error {
	csrfToken := c.FormValue(fmt.Sprintf("%s_csrf", os.Getenv("API_NAME")))

	if csrfToken == "" {
		rp := models.ResponsePacket{Error: true, Code: "status_disallowed", Message: "Not allowed."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if csrfToken != c.Cookies(fmt.Sprintf("%s_csrf", os.Getenv("API_NAME"))) {
		rp := models.ResponsePacket{Error: true, Code: "csrf_mismatch", Message: "CSRF token mismatch."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	return c.Next()
}

func InvalidateCSRFToken(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "csrf_",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set to a past time to expire the cookie
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return nil
}
