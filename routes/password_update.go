package routes

import (
	"errors"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PasswordUpdate(c *fiber.Ctx) error {

	if err := controller.ValidateCSRFToken(c); err != nil {
		return err
	}

	email := c.FormValue("email")
	oldpassword := string(c.FormValue("password"))
	newpassword := string(c.FormValue("newpassword"))

	if email == "" || oldpassword == "" || newpassword == "" {
		rp := models.ResponsePacket{Error: true, Code: "empty_fields", Message: "Missing required fields."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if string(oldpassword) == (string(newpassword)) {
		rp := models.ResponsePacket{Error: true, Code: "same_password", Message: "New password cannot be the same as old password."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !controller.PasswordIsValid(string(newpassword)) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_password", Message: "Password is not strong enough."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	newHashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newpassword), 12)

	var auth models.User

	if err := database.DB.Model(&auth).Where("email = ?", email).Updates(map[string]interface{}{"password": string(newHashedPassword)}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Unable to update password for user. User not found."}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not update password"}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthroized"})
}
