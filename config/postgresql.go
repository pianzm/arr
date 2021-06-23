package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	connectionPattern = "host=%s user=%s password=%s dbname=%s sslmode=disable"
)

// WritePostgresDB function for creating database connection for write-access
func WritePostgresDB(cfg *Config) *sql.DB {
	return CreateDBConnection(fmt.Sprintf(connectionPattern, cfg.PostgreWriteDBHost, cfg.PostgreWriteDBUser, cfg.PostgreWriteDBPassword, cfg.PostresWriteDBName))
}

// CreateDBConnection function for creating database connection
func CreateDBConnection(descriptor string) *sql.DB {
	db, err := sql.Open("postgres", descriptor)
	if err != nil {
		defer db.Close()
		return nil
	}
	return db
}
