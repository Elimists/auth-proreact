package config

import (
	"os"
)

var (
	smtpUserName string
	smtpPassword string
	smtpHost     string
)

func LoadSMTP() {
	smtpUserName = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	smtpHost = os.Getenv("SMTP_HOST")
}

func GetSMTPUserName() string {
	return smtpUserName
}

func GetSMTPPassword() string {
	return smtpPassword
}

func GetSMTPHost() string {
	return smtpHost
}
