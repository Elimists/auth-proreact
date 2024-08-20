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
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

func LoadKeys() {
	privateKeyData, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		log.Fatal("Error reading private key: " + err.Error())
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatal("Error parsing private key: " + err.Error())
	}

	publicKeyData, err := os.ReadFile("keys/public_key.pem")
	if err != nil {
		log.Fatal("Error reading public key: " + err.Error())
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
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
