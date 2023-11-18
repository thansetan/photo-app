package dtos

import (
	"mime/multipart"
)

type CreatePhotoRequest struct {
	Title     string                `form:"title" binding:"required" example:"I'm cool"`
	Caption   string                `form:"caption" example:"A cool photo of me"`
	IsPrivate bool                  `form:"is_private" description:"Do you want to make this photo private?"`
	Photo     *multipart.FileHeader `form:"photo" swaggerignore:"true"`
}

type CreatePhotoResponse struct {
	ID string `json:"photo_id"`
}

type UpdatePhotoRequest struct {
	Title     *string `json:"title" example:"I'm very cool"`
	IsPrivate *bool   `json:"is_private"`
	Caption   *string `json:"caption" example:"A very cool photo of me"`
}

type PhotoResponse struct {
	ID        string `json:"photo_id"`
	Title     string `json:"title"`
	Caption   string `json:"caption,omitempty"`
	PhotoPath string `json:"photo_path"`

	Owner *UserResponse `json:"owner,omitempty"`
}
