package routes

import (
	"github.com/Elimists/go-app/middleware"
	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {

	// Auth routes

	app.Post("/register", middleware.Limiter(14, 60), Register)
	app.Get("/verify", middleware.Limiter(1, 60), VerifyEmail)

	app.Post("/login", middleware.Limiter(6, 45), Login)
	app.Post("/logout", middleware.Limiter(6, 45), Logout)

	app.Get("/csrf", middleware.Limiter(6, 45), CSRFToken)

	app.Post(
		"/password-reset",
		middleware.Limiter(6, 45),
		PasswordReset)

	app.Post(
		"/password-update",
		middleware.Protected(),
		middleware.Limiter(6, 45),
		PasswordUpdate)

}
