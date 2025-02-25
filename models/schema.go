package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string         `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Password       string         `gorm:"not null;size:255" json:"-"`
	JournalEntries []JournalEntry `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"journal_entries"`
	Summaries      []Summary      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"summaries"`
}

type JournalEntry struct {
	gorm.Model
	UserID           uint      `gorm:"not null" json:"user_id"`
	Date             time.Time `gorm:"not null" json:"date"`
	EncryptedContent string    `gorm:"not null" json:"content"`
}

type Summary struct {
	gorm.Model
	UserID     uint   `gorm:"not null" json:"user_id"`
	WeekNumber uint   `gorm:"not null" json:"week_number"`
	Year       uint   `gorm:"not null" json:"year"`
	Summary    string `gorm:"not null" json:"summary"`
}

func getEncryptionKey() ([]byte, error) {
	key := os.Getenv("JOURNAL_ENCRYPTION_KEY")
	if key == "" {
		log.Println("ERROR: JOURNAL_ENCRYPTION_KEY environment variable not set")
		return nil, fmt.Errorf("JOURNAL_ENCRYPTION_KEY environment variable not set")
	}

	// Decode from base64 if you stored it that way
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		log.Printf("ERROR: invalid encryption key format: %v\n", err)
		return nil, fmt.Errorf("invalid encryption key format: %v", err)
	}

	// Verify key length for AES-256
	if len(decoded) != 32 {
		log.Printf("ERROR: encryption key must be exactly 32 bytes (got %d)\n", len(decoded))
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes (got %d)", len(decoded))
	}

	log.Println("INFO: Encryption key successfully retrieved")
	return decoded, nil
}

func Encrypt(text string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("ERROR: failed to create cipher block: %v\n", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("ERROR: failed to create GCM: %v\n", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("ERROR: failed to read nonce: %v\n", err)
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}

func Decrypt(encryptedText string) (string, error) {

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		log.Printf("ERROR: failed to decode base64 ciphertext: %v\n", err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("ERROR: failed to create cipher block: %v\n", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Printf("ERROR: failed to create GCM: %v\n", err)
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		log.Println("ERROR: ciphertext too short")
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Printf("ERROR: failed to decrypt text: %v\n", err)
		return "", err
	}

	decryptedText := string(plaintext)

	return decryptedText, nil
}
