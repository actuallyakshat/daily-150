package routes

import (
	controllers "daily-150/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(api fiber.Router) {
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/logout", controllers.Logout)
	api.Get("/me", controllers.Me)
}
