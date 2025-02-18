package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string         `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Password       string         `gorm:"not null;size:255" json:"-"`
	JournalEntries []JournalEntry `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"journal_entries"`
}

type JournalEntry struct {
	gorm.Model
	UserID  uint      `gorm:"not null" json:"user_id"`
	Date    time.Time `gorm:"not null" json:"date"`
	Content string    `gorm:"not null" json:"content"`
}
