package controllers

import (
	"context"
	"daily-150/helper"
	"daily-150/initialisers"
	"daily-150/models"
	"log"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ExtensionLogin(c *fiber.Ctx) error {
	db := initialisers.DB
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	user := models.User{}
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	if !verifyPassword(body.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	token, err := generateJWT(user.Username)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "Login successful",
	})
}

func DidUserJournalToday(c *fiber.Ctx) error {
	db := initialisers.DB
	redis := initialisers.RedisClient

	username, ok := helper.GetUsername(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to view this entry",
		})
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	//Check cache
	ctx := context.Background()
	currentDate := time.Now().UTC().Format("2006-01-02")
	redisKey := fmt.Sprintf("daily-150:journal-today:%s:%d", currentDate, user.ID)
	if result, err := redis.Get(ctx, redisKey).Result(); err == nil {
		if result == "true" {
			log.Println("Serving from cache")
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  true,
				"message": "You have journaled today",
			})
		}
	}

	var entries []models.JournalEntry
	if err := db.Where("user_id = ? AND DATE(date) = CURRENT_DATE", user.ID).Find(&entries).Error; err != nil {
		return helper.HandleError(c, err)
	}

	if len(entries) > 0 {

		//We had a cache miss so update the cache
		_, err := redis.Set(ctx, redisKey, "true", 24*time.Hour).Result()
		if err != nil {
			log.Println("Error saving to cache:", err)
		}

		log.Println("Saved in cache")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "You have journaled today",
		})
	} else {
		_, err := redis.Set(ctx, redisKey, "false", 24*time.Hour).Result()
		if err != nil {
			log.Println("Error saving to cache:", err)
		}

		log.Println("Saved in cache")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  false,
		"message": "You have not journaled today",
	})
}
