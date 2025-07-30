// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// const (
// 	baseURL     = "http://localhost:8080/api/v1"
// 	totalUsers  = 200
// 	concurrency = 20
// )

// type UploadPhotoRequest struct {
// 	UserID   uint64 `json:"user_id"`
// 	PhotoURL string `json:"photo_url"`
// }

// type UploadSelfieRequest struct {
// 	UserID    uint64 `json:"user_id"`
// 	SelfieURL string `json:"selfie_url"`
// }

// var samplePhotos = []string{
// 	"https://kumparan24.com/wp-content/uploads/2025/04/8aae01e5b7379cc1877294bebf3a3e06.jpg",
// 	"https://upload.wikimedia.org/wikipedia/commons/e/e6/Jackie_Chan_Cannes_2012.jpg",
// }

// var sampleSelfies = []string{
// 	"https://kumparan24.com/wp-content/uploads/2025/04/8aae01e5b7379cc1877294bebf3a3e06.jpg",
// 	// "https://media.discordapp.net/attachments/855653642659954712/1397115003982844004/IMG_8368.webp?ex=68808c3e&is=687f3abe&hm=727c29740d5a0e1987077aea06ba4c2cb06d3a5b572ea4a8ad1a6ff988cc9c4e&=&format=webp",
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())
// 	var wg sync.WaitGroup
// 	sem := make(chan struct{}, concurrency)

// 	userid := 1
// 	for i := 1; i <= totalUsers; i++ {
// 		if userid > 2 {
// 			userid = 1
// 		}
// 		wg.Add(1)
// 		sem <- struct{}{} // limit concurrency

// 		go func(userID uint64) {
// 			defer wg.Done()
// 			defer func() { <-sem }()

// 			uploadSelfie(userID)
// 			for j := 0; j < 7; j++ {
// 				uploadPhoto(userID)
// 			}
// 		}(uint64(i))
// 		userid++
// 	}

// 	wg.Wait()
// 	fmt.Println("Load testing completed.")
// }

// func uploadSelfie(userID uint64) {
// 	selfie := UploadSelfieRequest{
// 		UserID:    userID,
// 		SelfieURL: sampleSelfies[rand.Intn(len(sampleSelfies))],
// 	}
// 	sendPostRequest("/selfies", selfie)
// }

// func uploadPhoto(userID uint64) {
// 	photo := UploadPhotoRequest{
// 		UserID:   userID,
// 		PhotoURL: samplePhotos[rand.Intn(len(samplePhotos))],
// 	}
// 	sendPostRequest("/photos", photo)
// }

// func sendPostRequest(path string, payload interface{}) {
// 	jsonData, _ := json.Marshal(payload)
// 	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Printf("Error POST %s: %v\n", path, err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
// 		fmt.Printf("Failed POST %s. Status: %d\n", path, resp.StatusCode)
// 	}
// }
