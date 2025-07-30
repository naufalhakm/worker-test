package main

import (
	"context"
	"go-worker/internal/repository"
	"go-worker/internal/service"
	"go-worker/pkg/database"
	"go-worker/pkg/rabbitmq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT/SIGTERM agar bisa shutdown dengan graceful
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	db, err := database.NewMySQLClient()
	if err != nil {
		log.Fatal("Could not connect to MySQL:", err)
	}
	photoRepo := repository.NewPhotoRepository(db)
	selfieRepo := repository.NewSelfieRepository(db)
	photoService := service.NewPhotoService(photoRepo, selfieRepo, nil)

	consumer := rabbitmq.NewConsumer("photo-upload-queue", 60, func(ctx context.Context, payload rabbitmq.VerifyPayload) error {
		return photoService.VerifyPhotoMatch(ctx, payload.UserID)
	})

	// log.Println("Starting consumer...")
	// for {
	// 	err := consumer.Start(ctx)
	// 	if err != nil {
	// 		log.Printf("Consumer error: %v. Retrying in 5s...", err)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Consumer exited with error: %v", err)
	}
}
