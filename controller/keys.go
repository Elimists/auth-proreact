package controller

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

func GetPrivateKey() *rsa.PrivateKey {
	return PrivateKey
}

func GetPublicKey() *rsa.PublicKey {
	return PublicKey
}

func LoadKeys() {
	privateKeyData, err := os.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		log.Fatal("Error reading private key: " + err.Error())
	}

	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatal("Error parsing private key: " + err.Error())
	}

	publicKeyData, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		log.Fatal("Error reading public key: " + err.Error())
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Fatal("Error parsing public key: " + err.Error())
	}

}

func InvalidateJWT(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     fmt.Sprintf("%s_jwt", os.Getenv("API_NAME")),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		SameSite: "Lax",
	})

	return nil
}
