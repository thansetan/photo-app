package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"photo-app/dtos"
	"photo-app/helpers"
	"photo-app/models"
	"photo-app/repositories"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserController interface {
	Register(context.Context, dtos.UserRegister) (dtos.RegisterResponse, error)
	Login(context.Context, dtos.UserLogin) (dtos.LoginResponse, error)
	Update(context.Context, dtos.UserUpdateRequest) error
	Delete(context.Context) error
}

type userController struct {
	repo   repositories.UserRepository
	logger *slog.Logger
}

func NewUserController(repo repositories.UserRepository, logger *slog.Logger) UserController {
	return &userController{repo, logger}
}

func (c *userController) Register(ctx context.Context, data dtos.UserRegister) (dtos.RegisterResponse, error) {
	var res dtos.RegisterResponse
	user := models.User{
		Username: data.Username,
		Email:    data.Email,
	}

	h, err := helpers.HashPassword([]byte(data.Password))
	if err != nil {
		c.logger.Error("User [REGISTER]", "error", err.Error())
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}
	user.ID = uuid.NewString()
	user.Password = string(h)

	userID, err := c.repo.Create(ctx, user)
	if err != nil {
		c.logger.Error("User [REGISTER]", "error", err.Error())
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return res, helpers.NewResponseError(errors.New("user with provided username/email already exists"), http.StatusConflict)
		}
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	res.ID = userID

	return res, nil
}

func (c *userController) Login(ctx context.Context, data dtos.UserLogin) (dtos.LoginResponse, error) {
	var res dtos.LoginResponse

	user, err := c.repo.FindByEmail(ctx, data.Email)
	if err != nil {
		c.logger.Error("User [LOGIN]", "error", err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, helpers.NewResponseError(errors.New("incorrect email/password"), http.StatusUnauthorized)
		}
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)

	}

	err = helpers.ComparePassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		c.logger.Error("User [LOGIN]", "error", err.Error())
		return res, helpers.NewResponseError(errors.New("incorrect email/password"), http.StatusUnauthorized)
	}

	res.Token, err = helpers.GenerateJWT(user.ID)
	if err != nil {
		c.logger.Error("User [LOGIN]", "error", err.Error())
		return res, helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return res, nil
}

func (c *userController) Update(ctx context.Context, data dtos.UserUpdateRequest) error {
	id, ok := ctx.Value("id").(string)
	if !ok {
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	if data.Password == data.NewPassword {
		return helpers.NewResponseError(errors.New("old and new password can't be the same"), http.StatusBadRequest)
	}

	user, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Error("User [UPDATE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	err = helpers.ComparePassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		c.logger.Error("User [UPDATE]", "error", err.Error())
		return helpers.NewResponseError(errors.New("invalid password"), http.StatusUnauthorized)
	}

	// if the value doesn't change, gorm will automatically handles it (not updating the data).
	user.Email = data.Email
	if data.NewPassword != "" {
		newPass, err := helpers.HashPassword([]byte(data.NewPassword))
		if err != nil {
			c.logger.Error("User [UPDATE]", "error", err.Error())
			return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
		}
		user.Password = string(newPass)
	}
	user.Username = data.Username

	err = c.repo.Update(ctx, user)
	if err != nil {
		c.logger.Error("User [UPDATE]", "error", err.Error())
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			return helpers.NewResponseError(errors.New("user with provided username/email already exist"), http.StatusConflict)
		}
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (c *userController) Delete(ctx context.Context) error {
	id, ok := ctx.Value("id").(string)
	if !ok {
		helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	user, err := c.repo.FindByID(ctx, id)
	if err != nil {
		c.logger.Error("User [DELETE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	err = c.repo.Delete(ctx, user)
	if err != nil {
		c.logger.Error("User [DELETE]", "error", err.Error())
		return helpers.NewResponseError(helpers.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}
