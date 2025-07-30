package repository

import (
	"context"
	"fmt"
	"go-worker/internal/entity"

	"gorm.io/gorm"
)

type PhotoRepository interface {
	Create(ctx context.Context, photo *entity.Photo) error
	FindByID(ctx context.Context, photo *entity.Photo) error
	Update(ctx context.Context, photo *entity.Photo) error
	FindByAll(ctx context.Context, photo *[]entity.Photo, userID uint64) error
}

type PhotoRepositoryImpl struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &PhotoRepositoryImpl{
		db: db,
	}
}

func (repo *PhotoRepositoryImpl) Create(ctx context.Context, photo *entity.Photo) error {
	err := repo.db.WithContext(ctx).Create(&photo).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *PhotoRepositoryImpl) FindByID(ctx context.Context, photo *entity.Photo) error {
	err := repo.db.WithContext(ctx).First(&photo, photo.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("photo with id %d not found", photo.ID)
		}
		return err
	}
	return nil
}

func (repo *PhotoRepositoryImpl) Update(ctx context.Context, photo *entity.Photo) error {
	err := repo.db.WithContext(ctx).Save(&photo).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *PhotoRepositoryImpl) FindByAll(ctx context.Context, photos *[]entity.Photo, userID uint64) error {
	err := repo.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "pending").Find(photos).Limit(7).Error
	if err != nil {
		return err
	}
	return nil
}
