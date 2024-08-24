package controller

import (
	"os"
)

var (
	SMTPUsername string
	SMTPPassword string
	SMTPHost     string
)

func LoadSMTPConfig() {
	SMTPUsername = os.Getenv("SMTP_USERNAME")
	SMTPPassword = os.Getenv("SMTP_PASSWORD")
	SMTPHost = os.Getenv("SMTP_HOST")
}
