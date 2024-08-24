package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/joho/godotenv"
)

func init() {

	envFile := ".env"

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	controller.LoadAPIConfig()
	controller.LoadSMTPConfig()
	controller.LoadKeys()

	database.Connect()
	go controller.EmailVerificationWorker() // Start the email verification worker
}

func main() {

	app := fiber.New()

	app.Static("/", "./public")

	var corsConfig cors.Config
	if controller.APIEnvironment == "production" || controller.APIEnvironment == "staging" {
		corsConfig = cors.Config{
			AllowCredentials: true,
			AllowOrigins:     "http://localhost:3000", // Restrict to localhost:3000 in production
			AllowMethods:     "GET,POST,PUT,DELETE",
		}
	} else {
		corsConfig = cors.Config{
			AllowCredentials: true,
			AllowOrigins:     "*", // Allow any origin in development
			AllowMethods:     "GET,POST,PUT,DELETE",
		}
	}

	app.Use(cors.New(corsConfig))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:_csrf",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     60 * time.Minute,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "CSRF token mismatch",
			})
		},
	}))

	routes.AllRoutes(app)

	app.Listen(fmt.Sprintf(":%s", controller.APIPort))
}
