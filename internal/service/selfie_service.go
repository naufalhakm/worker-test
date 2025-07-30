package service

import (
	"context"
	"go-worker/internal/entity"
	"go-worker/internal/params"
	"go-worker/internal/repository"
	"time"
)

type SelfieService interface {
	UploadSelfie(ctx context.Context, req *params.UploadSelfieRequest) error
}

type SelfieServiceImpl struct {
	SelfieRepo repository.SelfieRepository
}

func NewSelfieService(SelfieRepository repository.SelfieRepository) SelfieService {
	return &SelfieServiceImpl{
		SelfieRepo: SelfieRepository,
	}
}

func (svc *SelfieServiceImpl) UploadSelfie(ctx context.Context, req *params.UploadSelfieRequest) error {
	selfie := entity.Selfie{
		UserID:    req.UserID,
		URL:       req.SelfieURL,
		CreatedAt: time.Now(),
	}

	err := svc.SelfieRepo.Create(ctx, &selfie)
	if err != nil {
		return err
	}
	return nil
}
