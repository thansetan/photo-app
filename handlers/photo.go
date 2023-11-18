package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"photo-app/controllers"
	"photo-app/dtos"
	"photo-app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PhotoHandler struct {
	c controllers.PhotoController
}

func NewPhotoHandler(c controllers.PhotoController) *PhotoHandler {
	return &PhotoHandler{c}
}

// AddPhoto godoc
//
//	@Summary		add photo
//	@Description	add photo for current user
//	@Tags			Photos
//	@Param			form	formData	dtos.CreatePhotoRequest	true	"data required to create new photo"
//	@Param			photo	formData	file					true	"the picture file"
//	@Produce		json
//	@Success		201	{object}	dtos.CreatePhotoResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos [post]
//	@Security		Bearer
func (h *PhotoHandler) Create(ctx *gin.Context) {
	var data dtos.CreatePhotoRequest
	if err := ctx.ShouldBind(&data); err != nil {
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

	if !helpers.IsImage(data.Photo.Header.Get("Content-Type")) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "photo must be an image",
		})
		return
	}

	resp, err := h.c.Create(ctx, data)
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

// GetAllPhotos godoc
//
//	@Summary		get all public photos
//	@Description	get all public photos
//	@Tags			Photos
//	@Produce		json
//	@Success		200	{object}	helpers.PhotosResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		404	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos [get]
func (h *PhotoHandler) GetAll(ctx *gin.Context) {
	photos, err := h.c.GetAll(ctx)
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

	ctx.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})
}

// GetMyPhotos godoc
//
//	@Summary		get all photos of current user
//	@Description	get all available photos of current user
//	@Tags			Photos
//	@Produce		json
//	@Success		200	{object}	helpers.PhotosResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		404	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos/my [get]
//	@Security		Bearer
func (h *PhotoHandler) GetMine(ctx *gin.Context) {
	id := ctx.GetString("id")
	fmt.Println(id)

	photos, err := h.c.GetByUserID(ctx, id)
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

	ctx.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})
}

// GetPhotoByOwner godoc
//
//	@Summary		get all public photo owned by a user
//	@Description	get all public photo owned by specified user by providing their username
//	@Tags			Photos
//	@Param			username	path	string	true	"owner's username"
//	@Produce		json
//	@Success		200	{object}	helpers.PhotosResponse
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		404	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos/by/{username} [get]
func (h *PhotoHandler) GetByOwner(ctx *gin.Context) {
	username := ctx.Param("username")

	photos, err := h.c.GetByOwner(ctx, username)
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

	ctx.JSON(http.StatusOK, gin.H{
		"photos": photos,
	})
}

// UpdatePhoto godoc
//
//	@Summary		update data of a photo
//	@Description	update data of a photo by given ID
//	@Tags			Photos
//	@Accept			json
//	@Param			id		path	string					true	"photo ID"
//	@Param			Body	body	dtos.UpdatePhotoRequest	true	"data required to update a photo"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		401	{object}	helpers.ErrorResponse
//	@Failure		404	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos/{id} [put]
//	@Security		Bearer
func (h *PhotoHandler) Update(ctx *gin.Context) {
	var data dtos.UpdatePhotoRequest
	photoID := ctx.Param("id")
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

	err := h.c.Update(ctx, data, photoID)
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

// Delete Photo godoc
//
//	@Summary		delete photo
//	@Description	delete a photo by given ID
//	@Tags			Photos
//	@Param			id	path	string	true	"photo ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	helpers.ErrorResponse
//	@Failure		401	{object}	helpers.ErrorResponse
//	@Failure		404	{object}	helpers.ErrorResponse
//	@Failure		500	{object}	helpers.ErrorResponse
//	@Router			/photos/{id} [delete]
//	@Security		Bearer
func (h *PhotoHandler) Delete(ctx *gin.Context) {
	photoID := ctx.Param("id")

	err := h.c.Delete(ctx, photoID)
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
