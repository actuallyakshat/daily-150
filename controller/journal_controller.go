package controllers

import (
	"daily-150/helper"
	"daily-150/initialisers"
	"daily-150/models"
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

	// Validate required fields
	if body.Content == "" || body.Date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both date and content are required",
		})
	}

	// Get user
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return helper.HandleError(c, err)
	}

	// Parse date
	parsedDate, err := time.Parse(time.RFC3339, body.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":        "Invalid date format. Please use RFC3339 format",
			"receivedDate": body.Date,
		})
	}

	// Encrypt content
	encryptedContent, err := models.Encrypt(body.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process entry",
		})
	}

	// Create new entry
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

	// Create response with decrypted content
	type EntryResponse struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Date      string `json:"date"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}

	response := EntryResponse{
		ID:        newEntry.ID,
		UserID:    newEntry.UserID,
		Date:      newEntry.Date.Format("2006-01-02"),
		Content:   body.Content, // Use original content instead of decrypting
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

	// Encrypt the new content before saving
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entry updated successfully",
	})
}

func GetEntryByDate(c *fiber.Ctx) error {
	db := initialisers.DB
	date := c.Params("date")
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
	if err := db.Where("date = ? AND user_id = ?", date, user.ID).First(&entry).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving entry",
		})
	}

	// Decrypt the content before sending
	decryptedContent, err := models.Decrypt(entry.EncryptedContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decrypting content",
		})
	}

	// Create a response struct that includes decrypted content
	type EntryResponse struct {
		ID        uint   `json:"id"`
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

	// Create a slice to hold decrypted entries
	type EntryResponse struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Date      string `json:"date"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	decryptedEntries := make([]EntryResponse, 0, len(entries))

	// Decrypt each entry's content
	for _, entry := range entries {
		decryptedContent, err := models.Decrypt(entry.EncryptedContent)
		if err != nil {
			continue // Skip entries that fail to decrypt
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
