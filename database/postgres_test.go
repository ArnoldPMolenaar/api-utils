package database

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestPostgresSQLConnection(t *testing.T) {
	origBuilder := connectionURLBuilderFn
	origOpen := gormOpenFn
	defer func() {
		connectionURLBuilderFn = origBuilder
		gormOpenFn = origOpen
	}()

	t.Run("url builder error", func(t *testing.T) {
		connectionURLBuilderFn = func(name string) (string, error) {
			return "", errors.New("builder failed")
		}
		gormOpenFn = func(dsn string) (*gorm.DB, error) {
			t.Fatal("gorm open should not be called when URL builder fails")
			return nil, nil
		}

		db, err := PostgresSQLConnection()
		if err == nil {
			t.Fatal("expected error")
		}
		if db != nil {
			t.Fatalf("expected nil db, got %v", db)
		}
	})

	t.Run("gorm open error", func(t *testing.T) {
		connectionURLBuilderFn = func(name string) (string, error) {
			if name != "postgres" {
				t.Fatalf("expected postgres name, got %q", name)
			}
			return "host=localhost", nil
		}
		gormOpenFn = func(dsn string) (*gorm.DB, error) {
			if dsn != "host=localhost" {
				t.Fatalf("unexpected dsn: %q", dsn)
			}
			return nil, errors.New("dial failed")
		}

		db, err := PostgresSQLConnection()
		if err == nil {
			t.Fatal("expected error")
		}
		if db != nil {
			t.Fatalf("expected nil db, got %v", db)
		}
	})
}

func TestParseConnMaxLifetimeFromEnv(t *testing.T) {
	t.Run("new duration env", func(t *testing.T) {
		t.Setenv("DB_MAX_LIFETIME", "2s")

		got, err := parseConnMaxLifetimeFromEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != 2*time.Second {
			t.Fatalf("expected 2s, got %v", got)
		}
	})

	t.Run("invalid new duration env", func(t *testing.T) {
		t.Setenv("DB_MAX_LIFETIME", "not-a-duration")

		if _, err := parseConnMaxLifetimeFromEnv(); err == nil {
			t.Fatal("expected error for invalid DB_MAX_LIFETIME")
		}
	})

	t.Run("unset env defaults to zero", func(t *testing.T) {
		t.Setenv("DB_MAX_LIFETIME", "")

		got, err := parseConnMaxLifetimeFromEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != 0 {
			t.Fatalf("expected 0 duration, got %v", got)
		}
	})
}



