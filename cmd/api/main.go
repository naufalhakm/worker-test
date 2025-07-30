package main

import (
	"fmt"
	"go-worker/internal/handler"
	"go-worker/internal/repository"
	"go-worker/internal/service"
	"go-worker/pkg/database"
	"go-worker/pkg/rabbitmq"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := database.NewMySQLClient()
	if err != nil {
		log.Fatal("Could not connect to MySQL:", err)
	}

	router := gin.Default()
	router.Use(gin.Logger(), CORS())

	publisher, err := rabbitmq.NewPublisher("photo-upload-queue")
	if err != nil {
		log.Fatal("Failed to init RabbitMQ:", err)
	}
	defer publisher.Close()

	photoRepo := repository.NewPhotoRepository(db)
	selfieRepo := repository.NewSelfieRepository(db)
	photoService := service.NewPhotoService(photoRepo, selfieRepo, publisher)

	photoHandler := handler.NewPhotoHandler(photoService)

	router.GET("/", func(ctx *gin.Context) {
		currentYear := time.Now().Year()
		message := fmt.Sprintf("API Photo Service %d", currentYear)

		ctx.JSON(http.StatusOK, message)
	})

	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{
			v1.POST("/photos", photoHandler.UploadPhoto)
			v1.POST("/selfies", photoHandler.UploadSelfie)
		}
	}

	router.Run(":8080")
}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, accept, access-control-allow-origin, access-control-allow-headers")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}
