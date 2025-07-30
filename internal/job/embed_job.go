package job

import (
	"context"
	"go-worker/internal/entity"
	"go-worker/internal/service"
	"log"
	"sync"
	"time"
)

func RunEmbedJob(photoSvc service.PhotoService, selfieSvc service.SelfieService) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("🚀 Running embed job")

			ctx := context.Background()

			selfies, _ := selfieSvc.GetSelfiesToEmbed(ctx, 1000)
			photos, _ := photoSvc.GetPhotosToEmbed(ctx, 1000)

			var wg sync.WaitGroup
			selfieSem := make(chan struct{}, 10)
			photoSem := make(chan struct{}, 70)

			// Process selfies
			for _, s := range selfies {
				wg.Add(1)
				selfieSem <- struct{}{}
				go func(selfie entity.Selfie) {
					defer wg.Done()
					defer func() { <-selfieSem }()
					if err := selfieSvc.EmbedSelfie(ctx, &selfie); err != nil {
						log.Printf("❌ selfie embed error: %v", err)
					}
				}(s)
			}

			// Process photos
			for _, p := range photos {
				wg.Add(1)
				photoSem <- struct{}{}
				go func(photo entity.Photo) {
					defer wg.Done()
					defer func() { <-photoSem }()
					if err := photoSvc.EmbedPhoto(ctx, &photo); err != nil {
						log.Printf("❌ photo embed error: %v", err)
					}
				}(p)
			}

			wg.Wait()
			log.Println("✅ Embed job finished")
		}
	}
}
