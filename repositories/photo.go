package repositories

import (
	"context"
	"photo-app/models"

	"gorm.io/gorm"
)

type PhotoRepository interface {
	Create(context.Context, models.Photo) (string, error)
	FindAll(context.Context) ([]models.Photo, error)
	FindByID(context.Context, string) (models.Photo, error)
	FindByUserID(context.Context, string) ([]models.Photo, error)
	Update(context.Context, models.Photo, map[string]any) error
	Delete(context.Context, models.Photo) error
}

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{db}
}

func (repo *photoRepository) Create(ctx context.Context, data models.Photo) (string, error) {
	err := repo.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data.ID, err
	}

	return data.ID, nil
}

func (repo *photoRepository) FindAll(ctx context.Context) ([]models.Photo, error) {
	var photos []models.Photo
	err := repo.db.WithContext(ctx).Preload("User").Find(&photos, "NOT is_private").Error
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (repo *photoRepository) FindByUserID(ctx context.Context, userID string) ([]models.Photo, error) {
	var photos []models.Photo

	err := repo.db.WithContext(ctx).Find(&photos, "(user_id = ? AND NOT is_private) OR user_id = ?", userID, ctx.Value("id")).Error
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (repo *photoRepository) FindByID(ctx context.Context, id string) (models.Photo, error) {
	var photo models.Photo

	err := repo.db.WithContext(ctx).First(&photo, "id = ? AND (NOT is_private OR user_id = ?)", id, ctx.Value("id")).Error
	if err != nil {
		return photo, err
	}

	return photo, nil
}

func (repo *photoRepository) Update(ctx context.Context, data models.Photo, toUpdate map[string]any) error {
	err := repo.db.WithContext(ctx).Model(&data).Updates(toUpdate).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *photoRepository) Delete(ctx context.Context, data models.Photo) error {
	err := repo.db.WithContext(ctx).Delete(&data).Error
	if err != nil {
		return err
	}

	return nil
}
