package utils

import (
	"fmt"
	"strings"
	"time"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

// GenerateTransactionNumber generates a unique transaction number
func GenerateTransactionNumber() string {
	return fmt.Sprintf("TRX-%d", time.Now().UnixNano())
}

// GenerateDownloadToken generates a secure download token
func GenerateDownloadToken(userID, productID uint) string {
	return fmt.Sprintf("%d-%d-%d", userID, productID, time.Now().Unix())
}
