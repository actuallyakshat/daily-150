package middlewares

import (
	"fmt"
	"log"
	"os"

	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var ignoredRoutes = []string{"/api/register", "/api/login", "/api/generate-summary", "/api/extension/login"}
var extension_routes = []string{"/api/extension/login", "/api/extension/did-user-journal-today", "/api/extension/me"}

func isIgnoredRoute(c *fiber.Ctx) bool {
	return slices.Contains(ignoredRoutes, c.Path())
}

func isExtensionRoute(c *fiber.Ctx) bool {
	return slices.Contains(extension_routes, c.Path())
}

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		log.Println("REQUEST RECEIVED")

		if isIgnoredRoute(c) {
			return c.Next()
		}

		log.Println("CHECKING AUTH")

		var tokenString string
		if isExtensionRoute(c) {
			log.Println("EXTENSION ROUTE ACTIVATED")
			tokenString = getExtensionRouteToken(c)
			log.Println("EXTENSION ROUTE TOKEN: ", tokenString)
		} else {
			log.Println("NORMAL ROUTE")
			tokenString = c.Cookies("token")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token cookie",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Println("Token Parsing Error:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username, ok := claims["username"].(string)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims",
				})
			}

			c.Locals("username", username)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
}

func getExtensionRouteToken(c *fiber.Ctx) string {
	token := c.Get("Authorization")
	if token == "" || len(token) < 8 || token[:7] != "Bearer " {
		return ""
	}

	return token[7:]
}
