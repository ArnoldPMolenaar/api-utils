package utils

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestIsBase64(t *testing.T) {
	if !IsBase64("data:text/plain;base64,SGVsbG8=") {
		t.Fatal("expected valid base64 data URI")
	}
	if IsBase64("SGVsbG8=") {
		t.Fatal("expected false when metadata prefix is missing")
	}
	if IsBase64("data:text/plain;base64,not-base64!") {
		t.Fatal("expected invalid base64 to be rejected")
	}
}

func TestExtractExtensionFromBase64(t *testing.T) {
	ext, err := ExtractExtensionFromBase64("data:image/png;base64,AAAA")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ext != "png" {
		t.Fatalf("expected png, got %q", ext)
	}

	if _, err := ExtractExtensionFromBase64("invalid"); err == nil {
		t.Fatal("expected error for invalid base64 format")
	}
	if _, err := ExtractExtensionFromBase64("data:imagepng;base64,AAAA"); err == nil {
		t.Fatal("expected error for invalid mime type format")
	}
}

func TestBuildQueryParams(t *testing.T) {
	if got := BuildQueryParams("appName", []string{"foo", "bar"}); got != "appName=foo&appName=bar" {
		t.Fatalf("unexpected query params: %q", got)
	}
	if got := BuildQueryParams("k", nil); got != "" {
		t.Fatalf("expected empty query params for nil values, got %q", got)
	}
}

func TestCaseConversionHelpers(t *testing.T) {
	if got := CamelcaseToPascalCase("helloWorld"); got != "HelloWorld" {
		t.Fatalf("unexpected pascal case: %q", got)
	}
	if got := CamelcaseToPascalCase(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}

	if got := PascalCaseToCamelcase("HelloWorld"); got != "helloWorld" {
		t.Fatalf("unexpected camel case: %q", got)
	}
	if got := PascalCaseToCamelcase(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestJoinMapBool(t *testing.T) {
	joined := JoinMapBool(map[string]bool{"a": true, "b": false}, ",")
	parts := strings.Split(joined, ",")
	sort.Strings(parts)
	if !reflect.DeepEqual(parts, []string{"a:true", "b:false"}) {
		t.Fatalf("unexpected joined map values: %v", parts)
	}
}

func TestStringToUint(t *testing.T) {
	v, err := StringToUint("42")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}

	if _, err := StringToUint("-1"); err == nil {
		t.Fatal("expected error for negative number")
	}
}

func TestUintOrZero(t *testing.T) {
	if got := UintOrZero(nil); got != 0 {
		t.Fatalf("expected 0 for nil pointer, got %d", got)
	}
	v := uint(9)
	if got := UintOrZero(&v); got != 9 {
		t.Fatalf("expected 9, got %d", got)
	}
}

func TestURLEncode(t *testing.T) {
	if got := URLEncode("  Café Déjà Vu!  "); got != "cafe-deja-vu" {
		t.Fatalf("unexpected slug: %q", got)
	}
	if got := URLEncode("___Hello---World___"); got != "hello-world" {
		t.Fatalf("unexpected slug: %q", got)
	}
	if got := URLEncode("   "); got != "" {
		t.Fatalf("expected empty slug for whitespace input, got %q", got)
	}
}

