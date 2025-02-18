package routes

import (
	controllers "daily-150/controller"

	"github.com/gofiber/fiber/v2"
)

func JournalRouter(api fiber.Router) {
	api.Post("/entry", controllers.CreateEntry)
	api.Get("/entry", controllers.GetAllEntries)
	api.Get("/entry/:id", controllers.GetEntryByDate)
	api.Patch("/entry/:id", controllers.UpdateEntry)
	api.Delete("/entry/:id", controllers.DeleteEntry)
}
