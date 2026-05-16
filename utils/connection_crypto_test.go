package utils

import (
	"encoding/base64"
	"testing"
)

func TestConnectionURLBuilder(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "app")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "appdb")
	t.Setenv("DB_SSL_MODE", "disable")
	t.Setenv("VALKEY_HOST", "127.0.0.1")
	t.Setenv("VALKEY_PORT", "6379")
	t.Setenv("SERVER_HOST", "0.0.0.0")
	t.Setenv("SERVER_PORT", "8080")

	postgres, err := ConnectionURLBuilder("postgres")
	if err != nil {
		t.Fatalf("expected no error for postgres, got %v", err)
	}
	if postgres != "host=localhost port=5432 user=app password=secret dbname=appdb sslmode=disable" {
		t.Fatalf("unexpected postgres url: %q", postgres)
	}

	valkey, err := ConnectionURLBuilder("valkey")
	if err != nil {
		t.Fatalf("expected no error for valkey, got %v", err)
	}
	if valkey != "127.0.0.1:6379" {
		t.Fatalf("unexpected valkey url: %q", valkey)
	}

	fiberURL, err := ConnectionURLBuilder("fiber")
	if err != nil {
		t.Fatalf("expected no error for fiber, got %v", err)
	}
	if fiberURL != "0.0.0.0:8080" {
		t.Fatalf("unexpected fiber url: %q", fiberURL)
	}

	if _, err := ConnectionURLBuilder("unknown"); err == nil {
		t.Fatal("expected error for unsupported connection name")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := "01234567890123456789012345678901" // 32 bytes
	plaintext := "hello secure world"

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("expected encrypt success, got %v", err)
	}
	if ciphertext == "" {
		t.Fatal("expected non-empty ciphertext")
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("expected decrypt success, got %v", err)
	}
	if decrypted != plaintext {
		t.Fatalf("expected decrypted text %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptDecryptErrors(t *testing.T) {
	if _, err := Encrypt("short", "abc"); err == nil {
		t.Fatal("expected encrypt to fail with invalid key length")
	}

	key := "01234567890123456789012345678901"
	if _, err := Decrypt(key, "not-base64"); err == nil {
		t.Fatal("expected decrypt to fail for invalid base64")
	}

	shortCipher := base64.StdEncoding.EncodeToString([]byte("tiny"))
	if _, err := Decrypt(key, shortCipher); err == nil {
		t.Fatal("expected decrypt to fail for short ciphertext")
	}

	ciphertext, err := Encrypt(key, "secret")
	if err != nil {
		t.Fatalf("unexpected encrypt error: %v", err)
	}

	wrongKey := "abcdefghijklmnopqrstuvwxyzABCDEF" // 32 bytes, wrong key
	if _, err := Decrypt(wrongKey, ciphertext); err == nil {
		t.Fatal("expected decrypt to fail with wrong key")
	}
}

