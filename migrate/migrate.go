package migrate

import (
	"daily-150/initialisers"
	"daily-150/models"
)

func RunMigrations() {
	initialisers.DB.AutoMigrate(&models.User{}, &models.JournalEntry{})
}
