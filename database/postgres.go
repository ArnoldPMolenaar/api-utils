package database

import (
	"database/sql"
	"fmt"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func PostgresSQLConnection() (*gorm.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build PostgresSQL connection URL.
	postgresConnectionURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgresSQL.
	db, err := gorm.Open(postgres.Open(postgresConnectionURL), &gorm.Config{})
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
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

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
