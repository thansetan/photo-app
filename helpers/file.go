package helpers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

var (
	wd, _ = os.Getwd()
)

func SaveFile(ctx context.Context, file *multipart.FileHeader, photoID string) (string, error) {
	outputDir := filepath.Join(wd, os.Getenv("PHOTO_DIR"), ctx.Value("id").(string))
	outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s%s", photoID, filepath.Ext(file.Filename)))

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return outputFilePath, err
	}

	dst, err := os.Create(outputFilePath)
	if err != nil {
		return outputFilePath, err
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return outputFilePath, err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return outputDir, err
	}

	return filepath.ToSlash(strings.TrimPrefix(outputFilePath, filepath.Join(wd, os.Getenv("PHOTO_DIR")))), nil
}

func RemoveFile(filePath string) error {
	fullPath := filepath.Join(wd, os.Getenv("PHOTO_DIR"), filePath)
	err := os.Remove(fullPath)
	if err != nil {
		return err
	}

	return nil
}

func IsImage(fileType string) bool {
	return strings.HasPrefix(fileType, "image/")
}
