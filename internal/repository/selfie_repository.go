package repository

import (
	"context"
	"go-worker/internal/entity"

	"gorm.io/gorm"
)

type SelfieRepository interface {
	Create(ctx context.Context, selfie *entity.Selfie) error
	FindByUserID(ctx context.Context, selfie *entity.Selfie, userID uint64) error
	Update(ctx context.Context, selfie *entity.Selfie) error
}

type SelfieRepositoryImpl struct {
	db *gorm.DB
}

func NewSelfieRepository(db *gorm.DB) SelfieRepository {
	return &SelfieRepositoryImpl{
		db: db,
	}
}

func (repo *SelfieRepositoryImpl) Create(ctx context.Context, selfie *entity.Selfie) error {
	err := repo.db.WithContext(ctx).Create(&selfie).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *SelfieRepositoryImpl) FindByUserID(ctx context.Context, selfie *entity.Selfie, userID uint64) error {
	err := repo.db.WithContext(ctx).Where("user_id = ?", userID).First(&selfie).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *SelfieRepositoryImpl) Update(ctx context.Context, selfie *entity.Selfie) error {
	err := repo.db.WithContext(ctx).Save(&selfie).Error
	if err != nil {
		return err
	}
	return nil
}
