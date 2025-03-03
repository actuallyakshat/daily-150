package main

import (
	"daily-150/initialisers"
	"daily-150/middlewares"
	"daily-150/migrate"
	"daily-150/routes"
	"daily-150/routines"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initialisers.LoadEnv()
	initialisers.ConnectDB()
	initialisers.InitRedis()
	migrate.RunMigrations()
}

func main() {

	go routines.ProcessSummaries()

	app := fiber.New()

	setupMiddlewares(app)
	setupRoutes(app)
	setupStaticFiles(app)
	startServer(app)

}

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5173,http://localhost:8080, chrome-extension://jlmohemkiclhpibllpcbggcdopblnodn",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Upgrade, Connection, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions, Sec-WebSocket-Protocol",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge:           3600,
		ExposeHeaders:    "Set-Cookie",
	}))

	cookieEncryptionKey := os.Getenv("COOKIE_ENCRYPTION_KEY")
	if cookieEncryptionKey == "" {
		log.Fatalln("COOKIE_ENCRYPTION_KEY is not set")
	}

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cookieEncryptionKey,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
	}))
}

func setupRoutes(app *fiber.App) {
	app.Use(logger.New())
	api := app.Group("/api")
	api.Use(middlewares.CheckAuth())

	api.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "API rate limit exceeded",
			})
		},
	}))

	routes.IndexRouter(api)
}

func setupStaticFiles(app *fiber.App) {
	app.Static("/", "./client_build")

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./client_build/index.html")
	})
}

func startServer(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//For Railway
	err := app.Listen("0.0.0.0:" + port)

	if err != nil {
		panic(err)
	}
}
