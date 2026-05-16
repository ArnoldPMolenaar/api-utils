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

var (
	connectionURLBuilderFn = utils.ConnectionURLBuilder
	gormOpenFn             = func(dsn string) (*gorm.DB, error) {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
)

func PostgresSQLConnection() (*gorm.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, err := parseConnMaxLifetimeFromEnv()
	if err != nil {
		return nil, err
	}

	// Build PostgresSQL connection URL.
	postgresConnectionURL, err := connectionURLBuilderFn("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgresSQL.
	db, err := gormOpenFn(postgresConnectionURL)
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
	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetConnMaxLifetime(maxLifetimeConn)

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
