package routines

import (
	"bytes"
	"context"
	"daily-150/initialisers"
	"daily-150/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

func ProcessSummaries() {
	log.Println("PROCESSING SUMMARIES ROUTINE ACTIVE")
	ctx := context.Background()
	queueName := "summary_tasks"
	expressServerURL := os.Getenv("SUMMARY_SERVICE_URL")
	SUMMARISER_KEY := os.Getenv("SUMMARISER_KEY")
	db := initialisers.DB
	redisClient := initialisers.RedisClient

	if redisClient == nil {
		log.Println("Redis client not initialised")
		return
	}

	//Config for the worker.
	// maxConcurrent := 5 // Maximum number of concurrent batches process
	batchSize := 10 // Number of users to process in one API call
	// retryLimit := 3    // Maximum number of retries for a failed task

	for {
		batchTasks := []models.SummaryTask{}

		//Get tasks until we get a batch or the queue is empty
		for range batchSize {
			result, err := redisClient.BLPop(ctx, 1*time.Second, queueName).Result()
			if err != nil {
				// If it's a timeout, just continue to the next iteration
				if err == redis.Nil {
					break
				}
				log.Println("Error getting task from Redis queue: ", err)
				break
			}

			// The result contains [queueName, value]
			taskJSON := result[1]

			var task models.SummaryTask
			if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
				log.Println("Error unmarshalling task: ", err)
				continue
			}

			batchTasks = append(batchTasks, task)
		}

		if len(batchTasks) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		// Process the batch
		userEntries := make(map[uint][]string)
		for _, task := range batchTasks {
			userEntries[task.UserID] = task.Entries
		}

		// Send the batch to the Express server
		payload, err := json.Marshal(userEntries)
		if err != nil {
			log.Printf("Error marshalling payload: %v\n", err)
			// Put tasks back in the queue
			for _, task := range batchTasks {
				taskJSON, _ := json.Marshal(task)
				redisClient.RPush(ctx, queueName, taskJSON)
			}
			continue
		}

		req, err := http.NewRequest("POST", expressServerURL, bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Error creating request: %v\n", err)
			// Put tasks back in the queue
			for _, task := range batchTasks {
				taskJSON, _ := json.Marshal(task)
				redisClient.RPush(ctx, queueName, taskJSON)
			}
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", SUMMARISER_KEY)

		client := &http.Client{
			Timeout: 120 * time.Second, // Increased timeout for larger batches
		}

		log.Println("SENDING REQUEST TO SUMMARY SERVICE")
		resp, err := client.Do(req)

		if err != nil {
			log.Printf("Error sending data to Express server: %v\n", err)
			// Put tasks back in the queue
			for _, task := range batchTasks {
				taskJSON, _ := json.Marshal(task)
				redisClient.RPush(ctx, queueName, taskJSON)
			}
			continue
		}
		defer resp.Body.Close()

		var result map[uint]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Error decoding response from Express server: %v\n", err)
			// Put tasks back in the queue
			for _, task := range batchTasks {
				taskJSON, _ := json.Marshal(task)
				redisClient.RPush(ctx, queueName, taskJSON)
			}
			continue
		}

		now := time.Now().UTC()
		previousWeek := now.AddDate(0, 0, -7)
		year, week := previousWeek.ISOWeek()

		log.Println("SAVING SUMMARIES")
		for userID, summary := range result {
			encryptedSummary, err := models.Encrypt(summary)
			if err != nil {
				log.Printf("Error encrypting summary for user %d: %v\n", userID, err)
				continue
			}

			newSummary := models.Summary{
				UserID:     userID,
				WeekNumber: uint(week),
				Year:       uint(year),
				Summary:    encryptedSummary,
			}

			if err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "user_id"}, {Name: "week_number"}, {Name: "year"}},
				DoUpdates: clause.AssignmentColumns([]string{"summary"}),
			}).Create(&newSummary).Error; err != nil {
				log.Printf("Error saving summary for user %d: %v\n", userID, err)
			}
		}

		log.Printf("Processed batch of %d tasks successfully\n", len(batchTasks))
	}

}
