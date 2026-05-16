package cache

import (
	"errors"
	"testing"

	"github.com/valkey-io/valkey-go"
)

func TestValkeyConnection(t *testing.T) {
	origBuilder := connectionURLBuilderFn
	origNewClient := newClientFn
	defer func() {
		connectionURLBuilderFn = origBuilder
		newClientFn = origNewClient
	}()

	t.Run("url builder error", func(t *testing.T) {
		connectionURLBuilderFn = func(name string) (string, error) {
			return "", errors.New("builder failed")
		}
		newClientFn = func(opt valkey.ClientOption) (valkey.Client, error) {
			t.Fatal("new client should not be called when URL builder fails")
			return nil, nil
		}

		client, err := ValkeyConnection()
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

		connectionURLBuilderFn = func(name string) (string, error) {
			if name != "valkey" {
				t.Fatalf("expected valkey name, got %q", name)
			}
			return "127.0.0.1:6379", nil
		}

		var gotOpt valkey.ClientOption
		newClientFn = func(opt valkey.ClientOption) (valkey.Client, error) {
			gotOpt = opt
			return nil, errors.New("dial failed")
		}

		client, err := ValkeyConnection()
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
}

