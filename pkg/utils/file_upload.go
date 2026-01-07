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

// UploadFile handles file upload and returns the file path
func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// Validate file size
	if file.Size > MaxUploadSize {
		return "", errors.New("file size exceeds maximum limit of 5MB")
	}

	// Validate file type
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return "", errors.New("invalid file type. Only images are allowed")
	}

	// Create upload directory if not exists
	uploadDir := filepath.Join(UploadPath, folder)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return relative path
	return strings.Replace(filePath, "\\", "/", -1), nil
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		// Ignore error if file doesn't exist
		if !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// GetFileURL returns the full URL for a file
func GetFileURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	// In production, this should return CDN URL
	return "/" + strings.TrimPrefix(filePath, "./")
}
