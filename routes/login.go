package routes

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {

	if err := controller.ValidateCSRFToken(c); err != nil {
		return err
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		rp := models.ResponsePacket{Error: true, Code: "empty_fields", Message: "Missing required fields."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var user models.User

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "account_not_found", Message: "Account not found!"}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		rp := models.ResponsePacket{Error: true, Code: "incorrect_password", Message: "Password is not correct"}
		return c.Status(fiber.StatusBadRequest).JSON(rp)
	}

	longerLogin := c.FormValue("longerlogin")
	expiry := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	if longerLogin == "true" {
		expiry = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}

	claims := jwt.MapClaims{
		"email":     user.Email,
		"id":        user.ID,
		"verified":  user.Verified,
		"privilege": user.Privilege,
		"exp":       expiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(controller.GetPrivateKey())

	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Could not sign token."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	database.DB.Model(&user).Where("email = ?", email).Update("updated_at", time.Now()) // update the last logged in datetime

	c.Cookie(&fiber.Cookie{
		Name:     fmt.Sprintf("%s_jwt", os.Getenv("API_NAME")),
		Value:    signedToken,
		Expires:  expiry.Time,
		SameSite: "Lax",
	})

	rp := models.ResponsePacket{Error: false, Code: "successfull", Message: "Login successfull"}
	return c.Status(fiber.StatusOK).JSON(rp)
}
