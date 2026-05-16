package cache

import (
	"os"
	"strconv"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/valkey-io/valkey-go"
)

var (
	connectionURLBuilderFn = utils.ConnectionURLBuilder
	newClientFn            = valkey.NewClient
)

// ValkeyConnection func for connect to Valkey server.
func ValkeyConnection() (valkey.Client, error) {
	// Define Valkey database number.
	dbNumber, _ := strconv.Atoi(os.Getenv("VALKEY_DB_NUMBER"))

	// Build Valkey URL.
	valkeyURL, err := connectionURLBuilderFn("valkey")
	if err != nil {
		return nil, err
	}

	// Set Valkey options.
	options := valkey.ClientOption{
		InitAddress:       []string{valkeyURL},
		ForceSingleClient: true,
		SelectDB:          dbNumber,
		ClientName:        os.Getenv("VALKEY_CLIENT_NAME"),
		Username:          os.Getenv("VALKEY_USERNAME"),
		Password:          os.Getenv("VALKEY_PASSWORD"),
	}

	client, err := newClientFn(options)
	if err != nil {
		return nil, err
	}

	return client, nil
}
