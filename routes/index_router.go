package routes

import (
	"github.com/gofiber/fiber/v2"
)

func IndexRouter(api fiber.Router) {
	AuthRouter(api)
	JournalRouter(api)
	ExtensionRouter(api)
}
