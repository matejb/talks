package database

import (
	"database/sql"
	"fmt"
)

// DB wraps a database connection with helper methods.
// In the Kennedy layout, platform packages provide foundational,
// reusable infrastructure code - NOT business-specific logic.
type DB struct {
	*sql.DB
}

// Open connects to the database.
func Open(driverName, connString string) (*DB, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	return &DB{DB: db}, nil
}

// Close closes the underlying database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}
