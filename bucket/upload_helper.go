package bucket

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ValidateFile(header *multipart.FileHeader, maxFileSize int64) error {
	if header.Size == 0 {
		return fmt.Errorf("file cannot be empty")
	}

	if header.Size > maxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	return nil
}

func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().Format("20060102150405")
	uuid := uuid.New().String()[:8]

	return fmt.Sprintf("%s_%s_%s%s",
		sanitizeFilename(nameWithoutExt),
		timestamp,
		uuid,
		ext,
	)
}

func sanitizeFilename(filename string) string {
	// Replace any character that's not alphanumeric, dash, or underscore with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
	return reg.ReplaceAllString(filename, "_")
}
