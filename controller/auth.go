package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
	"unicode"

	"github.com/eapache/channels"
)

var VERIFIATION_QUEUE = channels.NewInfiniteChannel()

func EmailVerificationWorker() {
	for payload := range VERIFIATION_QUEUE.Out() {
		var data map[string]string
		err := json.Unmarshal([]byte(payload.(string)), &data)
		if err != nil {
			log.Printf("Error unmarshalling email payload: %s", err.Error())
			continue
		}
		email := data["email"]
		verificationCode := data["verificationCode"]

		//Send email
		if err := SendVerificationEmail(email, verificationCode); err != nil {
			log.Printf("Error sending verification email: %s", err.Error())
		}
	}
}

func SendPasswordResetEmail(email string) error {

	smtpPlainAuth := smtp.PlainAuth(
		"",
		SMTPUsername,
		SMTPPassword,
		SMTPHost,
	)

	to := []string{email}
	subject := "Subject: Password Reset\n"
	from := "no-reply@example.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<div  style="font-size:20px; font-family: Arial, serif;">
				<p>Hi there,</p>
				<p>It looks like you requested a password reset. If this was you, please click the link below to reset your password.</p>
				<p>If you did not request a password reset, please ignore this email.</p>
				<p>Reset your password here: <code style="font-weight: bold;"><button>Reset Password</button></p>
			</div>
		</html>
		`)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", smtpPlainAuth, from, to, msg)
	return err
}

func EmailIsValid(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}

func PasswordIsValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func GenerateVerificationCode() string {
	min := 10101
	max := 99999
	return strconv.Itoa((rand.Intn(max-min+1) + min))
}

func SendVerificationEmail(email string, verificationCode string) error {
	smtpPlainAuth := smtp.PlainAuth(
		"",
		SMTPUsername,
		SMTPPassword,
		SMTPHost,
	)

	to := []string{email}
	subject := "Subject: Verify your email\n"
	from := "no-reply@example.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<div  style="font-size:20px; font-family: Arial, serif;">
				<p>Please see your verification code below:</p>
				<p>Verification Code: "%s"</p>
			</div>
		</html>
		`, verificationCode)
	msg := []byte(subject + mime + body)

	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", smtpPlainAuth, from, to, msg)
	return err
}
