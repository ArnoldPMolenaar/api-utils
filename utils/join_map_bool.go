package utils

import (
	"fmt"
	"strings"
)

// JoinMapBool joins a map of string keys and boolean values into a single string.
func JoinMapBool(m map[string]bool, sep string) string {
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s:%t", k, v))
	}
	return strings.Join(parts, sep)
}
