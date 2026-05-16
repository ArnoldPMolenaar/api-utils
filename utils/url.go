package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
)

// URLEncode converts an arbitrary string into an ASCII-only, URL-safe slug.
// Rules:
// - Trim leading/trailing whitespace;
// - Lowercase ASCII only;
// - Strip diacritics by normalizing to NFD and removing combining marks;
// - Replace any run of non-ASCII alphanumeric characters (including spaces/punctuation) with a single '-';
// - Collapse repeated '-' and trim '-' from both ends.
func URLEncode(s string) string {
	// Fast path for empty input.
	if s == "" {
		return ""
	}

	// Normalize whitespace.
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Lowercase for consistency (ASCII only target).
	s = strings.ToLower(s)

	// Strip diacritics: decompose (NFD) and remove combining marks (Mn).
	// Example: "Café" -> decompose into 'Cafe' + combining acute, then drop combining marks -> "Cafe".
	decomposed := norm.NFD.String(s)
	s = runes.Remove(runes.In(unicode.Mn)).String(decomposed)

	// Replace any sequence of non-ASCII alphanumeric characters with '-'.
	// Keep only [a-z0-9]; everything else (including spaces and punctuation) becomes '-'.
	var nonAsciiAlnum = regexp.MustCompile(`[^a-z0-9]+`)
	s = nonAsciiAlnum.ReplaceAllString(s, "-")

	// Collapse multiple '-' that might have been introduced and trim edges.
	var dashes = regexp.MustCompile(`-+`)
	s = dashes.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	return s
}
