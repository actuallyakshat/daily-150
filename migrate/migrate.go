package main

import (
	"daily-150/initialisers"
	"daily-150/models"
)

func init() {
	initialisers.LoadEnv()
	initialisers.ConnectDB()
}

func main() {
	initialisers.DB.AutoMigrate(&models.User{}, &models.JournalEntry{})
}
