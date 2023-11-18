package repositories

import (
	"context"
	"photo-app/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(context.Context, models.User) (string, error)
	FindByEmail(context.Context, string) (models.User, error)
	FindByUsername(context.Context, string) (models.User, error)
	FindByID(context.Context, string) (models.User, error)
	Update(context.Context, models.User) error
	Delete(context.Context, models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) Create(ctx context.Context, data models.User) (string, error) {
	err := repo.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data.ID, err
	}

	return data.ID, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *userRepository) FindByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *userRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	var user models.User

	err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *userRepository) Update(ctx context.Context, data models.User) error {
	err := repo.db.WithContext(ctx).Updates(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Delete(ctx context.Context, data models.User) error {
	err := repo.db.WithContext(ctx).Delete(&data).Error
	if err != nil {
		return err
	}

	return nil
}
