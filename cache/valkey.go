package cache

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/valkey-io/valkey-go"
)

// ValkeyConfig defines connection behavior for the shared client.
type ValkeyConfig struct {
	DBNumber   int
	ClientName string
	Username   string
	Password   string
}

// ValkeyDeps defines client bootstrap dependencies to keep setup testable.
type ValkeyDeps struct {
	BuildConnectionURL func(name string) (string, error)
	NewClient          func(opt valkey.ClientOption) (valkey.Client, error)
}

// DefaultValkeyDeps returns production defaults for Valkey dependencies.
func DefaultValkeyDeps() ValkeyDeps {
	return ValkeyDeps{
		BuildConnectionURL: utils.ConnectionURLBuilder,
		NewClient:          valkey.NewClient,
	}
}

func (d ValkeyDeps) validate() error {
	if d.BuildConnectionURL == nil {
		return fmt.Errorf("valkey dependency BuildConnectionURL cannot be nil")
	}
	if d.NewClient == nil {
		return fmt.Errorf("valkey dependency NewClient cannot be nil")
	}
	return nil
}

// ValkeyConnection func for connect to Valkey server.
func ValkeyConnection() (valkey.Client, error) {
	cfg, err := ValkeyConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return ValkeyConnectionWithConfig(cfg, DefaultValkeyDeps())
}

// ValkeyConnectionWithConfig connects using explicit config and dependencies.
func ValkeyConnectionWithConfig(cfg ValkeyConfig, deps ValkeyDeps) (valkey.Client, error) {
	if err := deps.validate(); err != nil {
		return nil, err
	}

	// Build Valkey URL.
	valkeyURL, err := deps.BuildConnectionURL("valkey")
	if err != nil {
		return nil, err
	}

	// Set Valkey options.
	options := valkey.ClientOption{
		InitAddress:       []string{valkeyURL},
		ForceSingleClient: true,
		SelectDB:          cfg.DBNumber,
		ClientName:        cfg.ClientName,
		Username:          cfg.Username,
		Password:          cfg.Password,
	}

	client, err := deps.NewClient(options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ValkeyConfigFromEnv reads Valkey client config from environment variables.
func ValkeyConfigFromEnv() (ValkeyConfig, error) {
	dbNumber, err := parseValkeyDBNumber(os.Getenv("VALKEY_DB_NUMBER"))
	if err != nil {
		return ValkeyConfig{}, err
	}

	return ValkeyConfig{
		DBNumber:   dbNumber,
		ClientName: os.Getenv("VALKEY_CLIENT_NAME"),
		Username:   os.Getenv("VALKEY_USERNAME"),
		Password:   os.Getenv("VALKEY_PASSWORD"),
	}, nil
}

func parseValkeyDBNumber(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid VALKEY_DB_NUMBER value %q: %w", raw, err)
	}

	return v, nil
}
