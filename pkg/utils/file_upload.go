package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxUploadSize = 5 << 20 // 5MB
	UploadPath    = "./uploads"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	if file.Size > MaxUploadSize {
		return "", errors.New("file size exceeds maximum limit of 5MB")
	}

	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return "", errors.New("invalid file type. Only images are allowed")
	}

	uploadDir := filepath.Join(UploadPath, folder)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return strings.Replace(filePath, "\\", "/", -1), nil
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func GetFileURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	return "/" + strings.TrimPrefix(filePath, "./")
}
