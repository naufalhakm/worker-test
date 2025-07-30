package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	baseURL             = "http://localhost:8080/api/v1"
	totalUsers          = 50000
	batchSize           = 5000
	concurrencyPerBatch = 100
	delayBetweenBatches = 30 * time.Second
)

type UploadPhotoRequest struct {
	UserID   uint64 `json:"user_id"`
	PhotoURL string `json:"photo_url"`
}

type UploadSelfieRequest struct {
	UserID    uint64 `json:"user_id"`
	SelfieURL string `json:"selfie_url"`
}

var samplePhotos = []string{
	"https://kumparan24.com/wp-content/uploads/2025/04/8aae01e5b7379cc1877294bebf3a3e06.jpg",
	"https://upload.wikimedia.org/wikipedia/commons/e/e6/Jackie_Chan_Cannes_2012.jpg",
}

var sampleSelfies = []string{
	"https://kumparan24.com/wp-content/uploads/2025/04/8aae01e5b7379cc1877294bebf3a3e06.jpg",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Log jumlah goroutine setiap 5 detik
	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Printf("🧵 Active Goroutines: %d\n", runtime.NumGoroutine())
		}
	}()

	userID := 1
	for batch := 1; userID <= totalUsers; batch++ {
		fmt.Printf("\n🚀 Starting batch %d (%d - %d)\n", batch, userID, min(userID+batchSize-1, totalUsers))

		var wg sync.WaitGroup
		sem := make(chan struct{}, concurrencyPerBatch)

		// Rate limiter: max 100 requests per second, burst 20
		limiter := rate.NewLimiter(rate.Limit(100), 20)

		batchEnd := min(userID+batchSize, totalUsers+1)
		for i := userID; i < batchEnd; i++ {
			wg.Add(1)
			sem <- struct{}{}

			go func(uid uint64) {
				defer wg.Done()
				defer func() { <-sem }()

				ctx := context.Background()

				// Limit each user processing
				if err := limiter.Wait(ctx); err != nil {
					fmt.Printf("❌ Rate limit error: %v\n", err)
					return
				}
				uploadSelfie(uid)

				for j := 0; j < 7; j++ {
					if err := limiter.Wait(ctx); err != nil {
						fmt.Printf("❌ Rate limit error: %v\n", err)
						return
					}
					uploadPhoto(uid)
				}
			}(uint64(i))
		}

		wg.Wait()
		fmt.Printf("✅ Completed batch %d (%d - %d)\n", batch, userID, batchEnd-1)
		userID = batchEnd

		// Delay before next batch
		if userID <= totalUsers {
			fmt.Printf("⏳ Waiting %v before next batch...\n", delayBetweenBatches)
			time.Sleep(delayBetweenBatches)
		}
	}

	fmt.Println("🏁 Load testing completed.")
}

func uploadSelfie(userID uint64) {
	payload := UploadSelfieRequest{
		UserID:    userID,
		SelfieURL: sampleSelfies[rand.Intn(len(sampleSelfies))],
	}
	sendPostRequest("/selfies", payload)
}

func uploadPhoto(userID uint64) {
	payload := UploadPhotoRequest{
		UserID:   userID,
		PhotoURL: samplePhotos[rand.Intn(len(samplePhotos))],
	}
	sendPostRequest("/photos", payload)
}

func sendPostRequest(path string, payload interface{}) {
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error POST %s: %v\n", path, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("⚠️  Failed POST %s. Status: %d\n", path, resp.StatusCode)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
