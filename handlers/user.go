package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"photo-app/controllers"
	"photo-app/dtos"
	"photo-app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	c controllers.UserController
}

func NewUserHandler(uc controllers.UserController) *UserHandler {
	return &UserHandler{uc}
}

// UserRegister godoc
//
//	@Summary		user register
//	@Description	create a new user account
//	@Tags			Users
//	@Accept			json
//	@Param			Body	body	dtos.UserRegister	true	"data required to create a new user"
//	@Produce		json
//	@Success		201	{object}	dtos.RegisterResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		429	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/users/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var data dtos.UserRegister

	if err := ctx.ShouldBindJSON(&data); err != nil {
		var errValidation validator.ValidationErrors
		if errors.As(err, &errValidation) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": helpers.GetValidationError(errValidation),
			})
			return
		}
		var errJSON *json.SyntaxError
		if errors.As(err, &errJSON) {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": "unable to parse JSON/invalid JSON format",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	resp, err := h.c.Register(ctx, data)
	if err != nil {
		var errController helpers.ResponseError
		if errors.As(err, &errController) {
			ctx.AbortWithStatusJSON(errController.Code(), gin.H{
				"error": errController.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": helpers.ErrInternal,
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// UserLogin godoc
//
//	@Summary		user login
//	@Description	login user. returns JWT
//	@Tags			Users
//	@Param			Body	body	dtos.UserLogin	true	"data required to login"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.LoginResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		401	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/users/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var data dtos.UserLogin

	if err := ctx.ShouldBindJSON(&data); err != nil {
		var errValidation validator.ValidationErrors
		if errors.As(err, &errValidation) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": helpers.GetValidationError(errValidation),
			})
			return
		}
		var errJSON *json.SyntaxError
		if errors.As(err, &errJSON) {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": "unable to parse JSON/invalid JSON format",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	resp, err := h.c.Login(ctx, data)
	if err != nil {
		var errController helpers.ResponseError
		if errors.As(err, &errController) {
			ctx.AbortWithStatusJSON(errController.Code(), gin.H{
				"error": errController.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": helpers.ErrInternal,
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UserUpdate godoc
//
//	@Summary		user update
//	@Description	update user data
//	@Tags			Users
//	@Param			Body	body	dtos.UserUpdateRequest	true	"data required to update user data"
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		401	{object}	helpers.ErrorResponse
//	@Failure		429	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/users/me [put]
//	@Security		Bearer
func (h *UserHandler) Update(ctx *gin.Context) {
	var data dtos.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		var errValidation validator.ValidationErrors
		if errors.As(err, &errValidation) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": helpers.GetValidationError(errValidation),
			})
			return
		}
		var errJSON *json.SyntaxError
		if errors.As(err, &errJSON) {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error": "unable to parse JSON/invalid JSON format",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	err := h.c.Update(ctx, data)
	if err != nil {
		var errController helpers.ResponseError
		if errors.As(err, &errController) {
			ctx.AbortWithStatusJSON(errController.Code(), gin.H{
				"error": errController.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": helpers.ErrInternal,
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteUser godoc
//
//	@Summary		delete update
//	@Description	delete user data and all photos related to this user
//	@Tags			Users
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		401	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/users/me [delete]
//	@Security		Bearer
func (h *UserHandler) Delete(ctx *gin.Context) {
	err := h.c.Delete(ctx)
	if err != nil {
		var errController helpers.ResponseError
		if errors.As(err, &errController) {
			ctx.AbortWithStatusJSON(errController.Code(), gin.H{
				"error": errController.Error(),
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": helpers.ErrInternal,
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
