package database

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestPostgresSQLConnection(t *testing.T) {
	deps := PostgresDeps{}

	t.Run("url builder error", func(t *testing.T) {
		deps = PostgresDeps{
			BuildConnectionURL: func(name string) (string, error) {
			return "", errors.New("builder failed")
			},
			Open: func(dsn string) (*gorm.DB, error) {
				t.Fatal("gorm open should not be called when URL builder fails")
				return nil, nil
			},
		}

		db, err := PostgresSQLConnectionWithConfig(PostgresConfig{}, deps)
		if err == nil {
			t.Fatal("expected error")
		}
		if db != nil {
			t.Fatalf("expected nil db, got %v", db)
		}
	})

	t.Run("gorm open error", func(t *testing.T) {
		deps = PostgresDeps{
			BuildConnectionURL: func(name string) (string, error) {
				if name != "postgres" {
					t.Fatalf("expected postgres name, got %q", name)
				}
				return "host=localhost", nil
			},
			Open: func(dsn string) (*gorm.DB, error) {
				if dsn != "host=localhost" {
					t.Fatalf("unexpected dsn: %q", dsn)
				}
				return nil, errors.New("dial failed")
			},
		}

		db, err := PostgresSQLConnectionWithConfig(PostgresConfig{}, deps)
		if err == nil {
			t.Fatal("expected error")
		}
		if db != nil {
			t.Fatalf("expected nil db, got %v", db)
		}
	})

	t.Run("invalid deps", func(t *testing.T) {
		_, err := PostgresSQLConnectionWithConfig(PostgresConfig{}, PostgresDeps{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestParseConnMaxLifetimeFromEnv(t *testing.T) {
	t.Run("duration env", func(t *testing.T) {
		t.Setenv("DB_MAX_LIFETIME", "2s")

		got, err := parseConnMaxLifetimeFromEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != 2*time.Second {
			t.Fatalf("expected 2s, got %v", got)
		}
	})

	t.Run("invalid duration env", func(t *testing.T) {
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

func TestPostgresConfigFromEnv(t *testing.T) {
	t.Run("parses full config", func(t *testing.T) {
		t.Setenv("DB_MAX_CONNECTIONS", "100")
		t.Setenv("DB_MAX_IDLE_CONNECTIONS", "10")
		t.Setenv("DB_MAX_LIFETIME", "2s")

		cfg, err := PostgresConfigFromEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.MaxConnections != 100 || cfg.MaxIdleConnections != 10 || cfg.MaxLifetime != 2*time.Second {
			t.Fatalf("unexpected config: %+v", cfg)
		}
	})

	t.Run("invalid max connections", func(t *testing.T) {
		t.Setenv("DB_MAX_CONNECTIONS", "oops")

		if _, err := PostgresConfigFromEnv(); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid max idle connections", func(t *testing.T) {
		t.Setenv("DB_MAX_CONNECTIONS", "100")
		t.Setenv("DB_MAX_IDLE_CONNECTIONS", "oops")

		if _, err := PostgresConfigFromEnv(); err == nil {
			t.Fatal("expected error")
		}
	})
}



