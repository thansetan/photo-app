package controllers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"photo-app/dtos"
	"photo-app/helpers"
	"photo-app/models"
	"photo-app/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoController interface {
	GetAll(context.Context) ([]dtos.PhotoResponse, error)
	GetByOwner(context.Context, string) ([]dtos.PhotoResponse, error)
	GetByUserID(context.Context, string) ([]dtos.PhotoResponse, error)
	Create(context.Context, dtos.CreatePhotoRequest) (dtos.CreatePhotoResponse, error)
	Update(context.Context, dtos.UpdatePhotoRequest, string) error
	Delete(context.Context, string) error
	IsAllowedToView(context.Context, string) (bool, error)
}

type photoController struct {
	repo     repositories.PhotoRepository
	userRepo repositories.UserRepository
	logger   *slog.Logger
}

func NewPhotoController(repo repositories.PhotoRepository, userRepo repositories.UserRepository, logger *slog.Logger) PhotoController {
	return &photoController{repo, userRepo, logger}
}

func (c *photoController) GetAll(ctx context.Context) ([]dtos.PhotoResponse, error) {
	photos, err := c.repo.FindAll(ctx)
	if err != nil {
		c.logger.Error("Photos [GET ALL]", "error", err.Error())
		return nil, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	data := make([]dtos.PhotoResponse, len(photos))
	for i, photo := range photos {
		data[i] = dtos.PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoPath: filepath.ToSlash(filepath.Join("/photos", photo.PhotoPath)),
			Owner: &dtos.UserResponse{
				Usename: photo.User.Username,
			},
		}
	}

	return data, nil
}

func (c *photoController) Create(ctx context.Context, data dtos.CreatePhotoRequest) (dtos.CreatePhotoResponse, error) {
	var res dtos.CreatePhotoResponse

	id, ok := ctx.Value("id").(string)
	if !ok {
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	photoID := uuid.NewString()

	filePath, err := helpers.SaveFile(ctx, data.Photo, photoID)
	if err != nil {
		c.logger.Error("Photos [CREATE]", "error", err.Error())
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	photo := models.Photo{
		Title:     data.Title,
		Caption:   data.Caption,
		PhotoPath: filePath,
		IsPrivate: data.IsPrivate,
		UserID:    id,
	}

	photo.ID = photoID
	_, err = c.repo.Create(ctx, photo)
	if err != nil {
		c.logger.Error("Photos [CREATE]", "error", err.Error())
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	res.ID = photoID
	return res, nil
}

func (c *photoController) Update(ctx context.Context, data dtos.UpdatePhotoRequest, id string) error {
	photo, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Error("Photos [UPDATE]", "error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.NewResponseError(errors.New("photo with specified ID can't be found"), http.StatusNotFound)
		}
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	userID, ok := ctx.Value("id").(string)
	if !ok {
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	if photo.UserID != userID {
		return helpers.NewResponseError(helpers.ErrNotAllowed, http.StatusUnauthorized)
	}

	toUpdate := make(map[string]any)
	if data.Caption != nil {
		toUpdate["caption"] = *data.Caption
	}
	if data.Title != nil {
		toUpdate["title"] = data.Title
	}
	if data.IsPrivate != nil {
		toUpdate["is_private"] = data.IsPrivate
	}

	err = c.repo.Update(ctx, photo, toUpdate)
	if err != nil {
		c.logger.Error("Photos [UPDATE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (c *photoController) Delete(ctx context.Context, id string) error {
	photo, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Error("Photos [DELETE]", "error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.NewResponseError(errors.New("photo with specified ID can't be found"), http.StatusNotFound)
		}
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	userID, ok := ctx.Value("id").(string)
	if !ok {
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	if photo.UserID != userID {
		return helpers.NewResponseError(helpers.ErrNotAllowed, http.StatusUnauthorized)
	}

	err = helpers.RemoveFile(photo.PhotoPath)
	if err != nil {
		c.logger.Error("Photos [DELETE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	err = c.repo.Delete(ctx, photo)
	if err != nil {
		c.logger.Error("Photos [DELETE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (c *photoController) GetByOwner(ctx context.Context, username string) ([]dtos.PhotoResponse, error) {
	user, err := c.userRepo.FindByUsername(ctx, username)
	if err != nil {
		c.logger.Error("Photos [GET BY OWNER]", "error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewResponseError(errors.New("user with specified username can't be found"), http.StatusNotFound)
		}
		return nil, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	data, err := c.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *photoController) GetByUserID(ctx context.Context, userID string) ([]dtos.PhotoResponse, error) {
	photos, err := c.repo.FindByUserID(ctx, userID)
	if err != nil {
		c.logger.Error("Photos [GET BY USER ID]", "error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewResponseError(errors.New("user with specified user_id can't be found"), http.StatusNotFound)
		}
		return nil, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	data := make([]dtos.PhotoResponse, len(photos))
	for i, photo := range photos {
		data[i] = dtos.PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoPath: filepath.ToSlash(filepath.Join("/photos", photo.PhotoPath)),
		}
	}

	return data, nil
}

func (c *photoController) IsAllowedToView(ctx context.Context, photoID string) (bool, error) {
	fmt.Println("USER ID = ", ctx.Value("id"))
	photo, err := c.repo.FindByID(ctx, photoID)
	if err != nil {
		return false, err
	}

	return (!photo.IsPrivate || photo.UserID == ctx.Value("id").(string)), nil
}
