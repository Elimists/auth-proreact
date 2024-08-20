package routes

import (
	"fmt"
	"strings"
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/eapache/channels"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var VERIFIATION_QUEUE = channels.NewInfiniteChannel()

func Register(c *fiber.Ctx) error {

	if err := controller.ValidateCSRFToken(c); err != nil {
		return err
	}

	email := c.FormValue("email")
	password := c.FormValue("password")
	password2 := c.FormValue("password2")

	if email == "" || password == "" {
		rp := models.ResponsePacket{Error: true, Code: "empty_fields", Message: "Missing required fields."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if !controller.EmailIsValid(email) {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Email is not valid."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if password != password2 {
		rp := models.ResponsePacket{Error: true, Code: "password_mismatch", Message: "Passwords do not match."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	// TODO check if password is valid

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	verificationCode := controller.GenerateVerificationCode()
	auth := models.User{
		Email:     string(email),
		Password:  hashedPassword,
		Privilege: 9, // General user.
		Verified:  false,
		UserVerification: models.UserVerification{
			VerificationCode:   verificationCode,
			VerificationExpiry: uint(time.Now().Add(time.Minute * 30).Unix()),
		},
	}

	userErr := database.DB.Create(&auth).Error

	//print userErr
	fmt.Println(userErr)

	if userErr != nil {
		if strings.Contains(userErr.Error(), "Duplicate entry") {
			rp := models.ResponsePacket{Error: true, Code: "duplicate_email", Message: "Email already exists!"}
			return c.Status(fiber.StatusNotAcceptable).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error. Could not register user."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	body := fmt.Sprintf(`{"email": "%s", "verificationCode": "%s"}`, string(email), string(verificationCode))
	VERIFIATION_QUEUE.In() <- body

	rp := models.ResponsePacket{Error: false, Code: "user_registered", Message: "User registered successfully."}
	return c.Status(fiber.StatusOK).JSON(rp)
}
