package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresConfig defines pooled connection settings.
type PostgresConfig struct {
	MaxConnections     int
	MaxIdleConnections int
	MaxLifetime        time.Duration
}

// PostgresDeps defines external dependencies to keep connection setup testable.
type PostgresDeps struct {
	BuildConnectionURL func(name string) (string, error)
	Open               func(dsn string) (*gorm.DB, error)
}

// DefaultPostgresDeps returns production defaults for Postgres dependencies.
func DefaultPostgresDeps() PostgresDeps {
	return PostgresDeps{
		BuildConnectionURL: utils.ConnectionURLBuilder,
		Open: func(dsn string) (*gorm.DB, error) {
			return gorm.Open(postgres.Open(dsn), &gorm.Config{})
		},
	}
}

func (d PostgresDeps) validate() error {
	if d.BuildConnectionURL == nil {
		return fmt.Errorf("postgres dependency BuildConnectionURL cannot be nil")
	}
	if d.Open == nil {
		return fmt.Errorf("postgres dependency Open cannot be nil")
	}
	return nil
}

func PostgresSQLConnection() (*gorm.DB, error) {
	cfg, err := PostgresConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return PostgresSQLConnectionWithConfig(cfg, DefaultPostgresDeps())
}

// PostgresSQLConnectionWithConfig starts a Postgres connection using explicit config and deps.
func PostgresSQLConnectionWithConfig(cfg PostgresConfig, deps PostgresDeps) (*gorm.DB, error) {
	if err := deps.validate(); err != nil {
		return nil, err
	}

	// Build PostgresSQL connection URL.
	postgresConnectionURL, err := deps.BuildConnectionURL("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgresSQL.
	db, err := deps.Open(postgresConnectionURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConn: the default is 0 (unlimited)
	// 	- SetMaxIdleConn: defaultMaxIdleConn = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Try to ping database.
	if err := sqlDB.Ping(); err != nil {
		// close database connection
		defer func(sqlDB *sql.DB) {
			err := sqlDB.Close()
			if err != nil {
				panic(fmt.Sprintf("error, not closed database connection, %v\n", err))
			}
		}(sqlDB)
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}

// PostgresConfigFromEnv reads database pool config from environment variables.
func PostgresConfigFromEnv() (PostgresConfig, error) {
	maxConnRaw := strings.TrimSpace(os.Getenv("DB_MAX_CONNECTIONS"))
	maxConn, err := parseIntEnv(maxConnRaw, "DB_MAX_CONNECTIONS")
	if err != nil {
		return PostgresConfig{}, err
	}

	maxIdleRaw := strings.TrimSpace(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxIdleConn, err := parseIntEnv(maxIdleRaw, "DB_MAX_IDLE_CONNECTIONS")
	if err != nil {
		return PostgresConfig{}, err
	}

	maxLifetimeConn, err := parseConnMaxLifetimeFromEnv()
	if err != nil {
		return PostgresConfig{}, err
	}

	return PostgresConfig{
		MaxConnections:     maxConn,
		MaxIdleConnections: maxIdleConn,
		MaxLifetime:        maxLifetimeConn,
	}, nil
}

func parseConnMaxLifetimeFromEnv() (time.Duration, error) {
	if raw := strings.TrimSpace(os.Getenv("DB_MAX_LIFETIME")); raw != "" {
		lifetime, err := time.ParseDuration(raw)
		if err != nil {
			return 0, fmt.Errorf("invalid DB_MAX_LIFETIME value %q: %w", raw, err)
		}
		return lifetime, nil
	}

	return 0, nil
}

func parseIntEnv(raw string, key string) (int, error) {
	if raw == "" {
		return 0, nil
	}

	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, raw, err)
	}

	return v, nil
}

