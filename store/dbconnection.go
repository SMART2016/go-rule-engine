package store

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
)

// Database struct to manage DB connection
type Database struct {
	DB *sql.DB
}

// NewDatabase initializes a new database connection
func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the database is reachable
	if err := db.Ping(); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	log.Println("Database connection established")
	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() {
	if d.DB != nil {
		err := d.DB.Close()
		if err != nil {
			return
		}
		log.Println("Database connection closed")
	}
}
