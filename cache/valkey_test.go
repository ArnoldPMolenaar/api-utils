package cache

import (
	"errors"
	"testing"

	"github.com/valkey-io/valkey-go"
)

func TestValkeyConnection(t *testing.T) {
	deps := ValkeyDeps{}

	t.Run("url builder error", func(t *testing.T) {
		deps = ValkeyDeps{
			BuildConnectionURL: func(name string) (string, error) {
				return "", errors.New("builder failed")
			},
			NewClient: func(opt valkey.ClientOption) (valkey.Client, error) {
				t.Fatal("new client should not be called when URL builder fails")
				return nil, nil
			},
		}

		client, err := ValkeyConnectionWithConfig(ValkeyConfig{}, deps)
		if err == nil {
			t.Fatal("expected error")
		}
		if client != nil {
			t.Fatalf("expected nil client, got %v", client)
		}
	})

	t.Run("constructs options and returns constructor error", func(t *testing.T) {
		t.Setenv("VALKEY_DB_NUMBER", "7")
		t.Setenv("VALKEY_CLIENT_NAME", "api-utils")
		t.Setenv("VALKEY_USERNAME", "user")
		t.Setenv("VALKEY_PASSWORD", "pass")

		var gotOpt valkey.ClientOption
		deps = ValkeyDeps{
			BuildConnectionURL: func(name string) (string, error) {
				if name != "valkey" {
					t.Fatalf("expected valkey name, got %q", name)
				}
				return "127.0.0.1:6379", nil
			},
			NewClient: func(opt valkey.ClientOption) (valkey.Client, error) {
				gotOpt = opt
				return nil, errors.New("dial failed")
			},
		}

		client, err := ValkeyConnectionWithConfig(valkeyConfigFromEnvMust(t), deps)
		if err == nil {
			t.Fatal("expected constructor error")
		}
		if client != nil {
			t.Fatalf("expected nil client, got %v", client)
		}
		if len(gotOpt.InitAddress) != 1 || gotOpt.InitAddress[0] != "127.0.0.1:6379" {
			t.Fatalf("unexpected init address: %+v", gotOpt.InitAddress)
		}
		if !gotOpt.ForceSingleClient || gotOpt.SelectDB != 7 {
			t.Fatalf("unexpected options: %+v", gotOpt)
		}
		if gotOpt.ClientName != "api-utils" || gotOpt.Username != "user" || gotOpt.Password != "pass" {
			t.Fatalf("unexpected auth/client options: %+v", gotOpt)
		}
	})

	t.Run("invalid deps", func(t *testing.T) {
		_, err := ValkeyConnectionWithConfig(ValkeyConfig{}, ValkeyDeps{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestValkeyConfigFromEnv(t *testing.T) {
	t.Run("parses config", func(t *testing.T) {
		t.Setenv("VALKEY_DB_NUMBER", "7")
		t.Setenv("VALKEY_CLIENT_NAME", "api-utils")
		t.Setenv("VALKEY_USERNAME", "user")
		t.Setenv("VALKEY_PASSWORD", "pass")

		cfg, err := ValkeyConfigFromEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.DBNumber != 7 || cfg.ClientName != "api-utils" || cfg.Username != "user" || cfg.Password != "pass" {
			t.Fatalf("unexpected config: %+v", cfg)
		}
	})

	t.Run("invalid db number", func(t *testing.T) {
		t.Setenv("VALKEY_DB_NUMBER", "oops")

		if _, err := ValkeyConfigFromEnv(); err == nil {
			t.Fatal("expected error")
		}
	})
}

func valkeyConfigFromEnvMust(t *testing.T) ValkeyConfig {
	t.Helper()
	cfg, err := ValkeyConfigFromEnv()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	return cfg
}

