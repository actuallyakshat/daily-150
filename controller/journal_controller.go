package controllers

import (
	"context"
	"daily-150/helper"
	"daily-150/initialisers"
	"daily-150/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateEntry(c *fiber.Ctx) error {
	db := initialisers.DB
	username, ok := helper.GetUsername(c)

	if !ok {
		return helper.HandleError(c, fiber.ErrUnauthorized)
	}

	type RequestBody struct {
		Date    string `json:"date"`
		Content string `json:"content"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	if body.Content == "" || body.Date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both date and content are required",
		})
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	parsedDate, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":        "Invalid date format. Please use RFC3339 format",
			"receivedDate": body.Date,
		})
	}

	encryptedContent, err := models.Encrypt(body.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process entry",
		})
	}

	var existingEntry models.JournalEntry
	if err := db.Where("user_id = ? AND date = ?", user.ID, parsedDate).First(&existingEntry).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Entry for this date already exists",
		})
	}

	newEntry := models.JournalEntry{
		UserID:           user.ID,
		Date:             parsedDate,
		EncryptedContent: encryptedContent,
	}

	if err := db.Create(&newEntry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create entry",
		})
	}

	type EntryResponse struct {
		ID        uint   `json:"ID"`
		UserID    uint   `json:"user_id"`
		Date      string `json:"date"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	response := EntryResponse{
		ID:        newEntry.ID,
		UserID:    newEntry.UserID,
		Date:      newEntry.Date.Format("2006-01-02"),
		Content:   body.Content,
		CreatedAt: newEntry.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Entry created successfully",
		"entry":   response,
	})
}

func DeleteEntry(c *fiber.Ctx) error {

	db := initialisers.DB
	username, ok := helper.GetUsername(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this entry",
		})
	}
	id := c.Params("id")
	entry := models.JournalEntry{}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	err := initialisers.DB.First(&entry, id).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving entry",
		})
	}

	if entry.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this entry",
		})
	}

	err = db.Delete(&entry).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting entry",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry deleted successfully",
	})
}

func UpdateEntry(c *fiber.Ctx) error {
	db := initialisers.DB
	username, ok := helper.GetUsername(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to update this entry",
		})
	}

	id := c.Params("id")
	entry := models.JournalEntry{}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	if err := db.First(&entry, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving entry",
		})
	}

	if entry.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to update this entry",
		})
	}

	type UpdateEntryRequest struct {
		Content string `json:"content"`
	}

	var updateEntryRequest UpdateEntryRequest
	if err := c.BodyParser(&updateEntryRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body",
		})
	}

	plainContent := updateEntryRequest.Content
	encryptedContent, err := models.Encrypt(updateEntryRequest.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error encrypting content",
		})
	}

	entry.EncryptedContent = encryptedContent

	if err := db.Save(&entry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating entry",
		})
	}

	entry.EncryptedContent = plainContent

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry updated successfully",
		"entry":   entry,
	})
}

func GetEntryByID(c *fiber.Ctx) error {
	db := initialisers.DB
	id := c.Params("id")
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

	var entry models.JournalEntry
	if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&entry).Error; err != nil {
		return helper.HandleError(c, err)
	}

	decryptedContent, err := models.Decrypt(entry.EncryptedContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decrypting content",
		})
	}

	type EntryResponse struct {
		ID        uint   `json:"ID"`
		UserID    uint   `json:"user_id"`
		Date      string `json:"date"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	response := EntryResponse{
		ID:        entry.ID,
		UserID:    entry.UserID,
		Date:      entry.Date.Format("2006-01-02"),
		Content:   decryptedContent,
		CreatedAt: entry.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: entry.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"entry": response,
	})
}

func GetAllEntries(c *fiber.Ctx) error {
	db := initialisers.DB
	username, ok := helper.GetUsername(c)

	if !ok {
		return helper.HandleError(c, fiber.ErrUnauthorized)
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	var entries []models.JournalEntry
	if err := db.Where("user_id = ?", user.ID).Find(&entries).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving entries",
		})
	}

	type EntryResponse struct {
		ID        uint   `json:"ID"`
		UserID    uint   `json:"user_id"`
		Date      string `json:"date"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	decryptedEntries := make([]EntryResponse, 0, len(entries))

	for _, entry := range entries {
		decryptedContent, err := models.Decrypt(entry.EncryptedContent)
		if err != nil {
			continue
		}

		decryptedEntry := EntryResponse{
			ID:        entry.ID,
			UserID:    entry.UserID,
			Date:      entry.Date.Format("2006-01-02"),
			Content:   decryptedContent,
			CreatedAt: entry.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: entry.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		decryptedEntries = append(decryptedEntries, decryptedEntry)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entries fetched successfully",
		"entries": decryptedEntries,
	})
}

func GenerateWeeklySummary(c *fiber.Ctx) error {
	db := initialisers.DB
	redisClient := initialisers.RedisClient
	now := time.Now().UTC()

	// Calculate the start and end of the previous week because cron job will run every monday
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())-7)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	CRON_ACTIVATION_KEY := os.Getenv("CRON_ACTIVATION_KEY")

	RECEIVED_CRON_ACTIVATION_KEY := c.Get("x-api-key")
	if RECEIVED_CRON_ACTIVATION_KEY != CRON_ACTIVATION_KEY {
		log.Println("RECEIVED_CRON_ACTIVATION_KEY: ", RECEIVED_CRON_ACTIVATION_KEY)
		log.Println("CRON_ACTIVATION_KEY: ", CRON_ACTIVATION_KEY)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var entries []models.JournalEntry
	if err := db.Where("date BETWEEN ? AND ?", startOfWeek, endOfWeek).Find(&entries).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving entries",
		})
	}

	log.Println("GOT THESE MANY ENTRIES: ", len(entries))
	log.Println("DECRYPTING ENTRIES")

	userEntries := make(map[uint][]string)
	for _, entry := range entries {
		decryptedContent, err := models.Decrypt(entry.EncryptedContent)

		if err != nil {
			continue
		}

		userEntries[entry.UserID] = append(userEntries[entry.UserID], decryptedContent)
	}

	ctx := context.Background()
	queueName := "summary_tasks"
	taskCount := 0

	for userID, entries := range userEntries {
		task := models.SummaryTask{
			UserID:  userID,
			Entries: entries,
		}

		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Println("Error marshalling task: ", err)
			continue
		}

		log.Println("PUSHING TASK INTO QUEUE ")

		if err := redisClient.RPush(ctx, queueName, taskJSON).Err(); err != nil {
			log.Println("Error pushing task to Redis queue: ", err)
			continue
		}
		taskCount++

	}

	// expressServerURL := os.Getenv("SUMMARY_SERVICE_URL")
	// payload, err := json.Marshal(userEntries)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Error marshalling payload",
	// 	})
	// }

	// SUMMARISER_KEY := os.Getenv("SUMMARISER_KEY")

	// req, err := http.NewRequest("POST", expressServerURL, bytes.NewBuffer(payload))
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Error creating request",
	// 	})
	// }

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("x-api-key", SUMMARISER_KEY)

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Error sending data to Express server" + err.Error(),
	// 	})
	// }
	// defer resp.Body.Close()

	// var result map[uint]string
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Error decoding response from Express server",
	// 	})
	// }

	previousWeek := now.AddDate(0, 0, -7)
	year, week := previousWeek.ISOWeek()
	batchInfoJSON, _ := json.Marshal(map[string]interface{}{
		"year":      year,
		"week":      week,
		"taskCount": taskCount,
		"timestamp": now.Format(time.RFC3339),
	})

	log.Println("SETTING REDIS BATCH INFO")
	redisClient.Set(ctx, fmt.Sprintf("batch:%d:%d", year, week), batchInfoJSON, 30*24*time.Hour)

	// for userID, summary := range result {
	// 	encryptedSummary, err := models.Encrypt(summary)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Error encrypting summary",
	// 		})
	// 	}
	// 	newSummary := models.Summary{
	// 		UserID:     userID,
	// 		WeekNumber: uint(week),
	// 		Year:       uint(year),
	// 		Summary:    encryptedSummary,
	// 	}
	// 	if err := db.Clauses(clause.OnConflict{
	// 		Columns:   []clause.Column{{Name: "user_id"}, {Name: "week_number"}, {Name: "year"}}, // Conflict target (unique constraint)
	// 		DoUpdates: clause.AssignmentColumns([]string{"summary"}),                             // Update the "summary" column on conflict
	// 	}).Create(&newSummary).Error; err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Error saving summary",
	// 		})
	// 	}
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Summary Request Queued.",
	})
}

func GetSummariesForUser(c *fiber.Ctx) error {
	db := initialisers.DB
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

	var summaries []models.Summary
	if err := db.Where("user_id = ?", user.ID).Find(&summaries).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving summaries",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"summaries": summaries,
	})
}

func GetSummaryByID(c *fiber.Ctx) error {
	db := initialisers.DB
	id := c.Params("id")
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

	var summary models.Summary
	if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&summary).Error; err != nil {
		return helper.HandleError(c, err)
	}

	decryptedSummary, err := models.Decrypt(summary.Summary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decrypting summary",
		})
	}

	type SummaryResponse struct {
		ID         uint   `json:"ID"`
		UserID     uint   `json:"user_id"`
		WeekNumber uint   `json:"week_number"`
		Summary    string `json:"summary"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	response := SummaryResponse{
		ID:         summary.ID,
		UserID:     summary.UserID,
		WeekNumber: summary.WeekNumber,
		Summary:    decryptedSummary,
		CreatedAt:  summary.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  summary.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"summary": response,
	})
}
