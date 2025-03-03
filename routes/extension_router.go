package routes

import (
	controllers "daily-150/controller"

	"github.com/gofiber/fiber/v2"
)

func ExtensionRouter(api fiber.Router) {
	api.Post("/extension/login", controllers.ExtensionLogin)
	api.Get("/extension/did-user-journal-today", controllers.DidUserJournalToday)
	api.Get("/extension/me", controllers.Me)
}
