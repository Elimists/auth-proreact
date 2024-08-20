package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Elimists/go-app/controller"
	"github.com/Elimists/go-app/database"
	"github.com/Elimists/go-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {

	controller.LoadKeys()

	envFile := ".env" // Defaults to local dev environment

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	database.Connect()
	go controller.EmailVerificationWorker() // Start the email verification worker
}

func main() {

	app := fiber.New()

	app.Static("/", "./public")

	var corsConfig cors.Config
	if os.Getenv("ENVIRONMENT") == "production" || os.Getenv("ENVIRONMENT") == "staging" {
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

	routes.AllRoutes(app)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("API_PORT")))
}
