package helpers

import "photo-app/dtos"

type ErrorResponse struct {
	Error string `json:"error"`
}

type PhotosResponse struct {
	Photos []dtos.PhotoResponse `json:"photos"`
}
