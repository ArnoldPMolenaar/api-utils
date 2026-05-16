package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

// IsBase64 checks if a string is a valid base64 encoded string.
func IsBase64(base64String string) bool {
	// Split the base64 string into metadata and data parts.
	parts := strings.Split(base64String, ",")
	if len(parts) < 2 {
		return false
	}

	// Test the string for base64 validity.
	_, err := base64.StdEncoding.DecodeString(parts[1])

	return err == nil
}

// ExtractExtensionFromBase64 extracts the file extension from a base64 string.
func ExtractExtensionFromBase64(base64String string) (string, error) {
	// Split the base64 string into metadata and data parts.
	parts := strings.Split(base64String, ";")
	if len(parts) < 2 {
		return "", errors.New("invalid base64 string format")
	}

	// Split the MIME type to get the extension.
	mimeParts := strings.Split(parts[0], "/")
	if len(mimeParts) != 2 {
		return "", errors.New("invalid MIME type format")
	}

	return mimeParts[1], nil
}
