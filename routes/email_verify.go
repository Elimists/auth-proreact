package routes

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func VerifyEmail(c *fiber.Ctx) error {

	verificationCode := c.Params("verificationCode")
	email := c.Params("email")

	if verificationCode == "" || email == "" {
		rp := models.ResponsePacket{Error: true, Code: "empty_code", Message: "Missing data in url."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	decodedEmail, _ := base64.StdEncoding.DecodeString(email)
	decodedVerificationCode, err := base64.StdEncoding.DecodeString(verificationCode)
	if err != nil {
		rp := models.ResponsePacket{Error: true, Code: "invalid_email", Message: "Invalid email."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var user models.User

	if err := database.DB.Where("email = ?", decodedEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "User not found."}
			return c.Status(fiber.StatusNotFound).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	if user.Verified {
		rp := models.ResponsePacket{Error: true, Code: "already_verified", Message: "Email is already verified."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	var userVerification models.UserVerification

	// Check if verification code has expired. If expired, send a new one.
	if uint(time.Now().Unix()) > userVerification.VerificationExpiry {
		if err := database.DB.Model(&userVerification).Where("email = ?", decodedEmail).Updates(map[string]interface{}{
			"verification_expiry": uint(time.Now().Add(time.Minute * 30).Unix())}).Error; err != nil {
			rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Verification time frame has expired, however server encountered problem while sending a new link. Please try again later."}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}

		verificationLink := fmt.Sprintf("%s/api/v2/verify/%s/%s", os.Getenv("API_URL"), email, verificationCode)

		body := fmt.Sprintf(`{"email": "%s", "verificationLink": "%s"}`, email, verificationLink)
		VERIFIATION_QUEUE.In() <- body

		rp := models.ResponsePacket{Error: true, Code: "expired", Message: "Verfication time frame has expired. A new link has been sent to your email."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	// Check if verification code matches.
	if userVerification.VerificationCode != string(decodedVerificationCode) {
		rp := models.ResponsePacket{Error: true, Code: "code_mismatch", Message: "Verification code does not match."}
		return c.Status(fiber.StatusNotAcceptable).JSON(rp)
	}

	if err := database.DB.Model(&userVerification).Where("email = ?", decodedEmail).Updates(map[string]interface{}{"verified": true, "verification_code": gorm.Expr("NULL")}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp := models.ResponsePacket{Error: true, Code: "not_found", Message: "Unable to update verification status for the user."}
			return c.Status(fiber.StatusInternalServerError).JSON(rp)
		}
		rp := models.ResponsePacket{Error: true, Code: "internal_error", Message: "Internal server error."}
		return c.Status(fiber.StatusInternalServerError).JSON(rp)
	}

	rp := models.ResponsePacket{Error: false, Code: "verified", Message: "Verification successfull."}
	return c.Status(fiber.StatusAccepted).JSON(rp)
}
