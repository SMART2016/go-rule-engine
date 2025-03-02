package store

import "log"

func InitializeEventStateStore(dsn string) (Querier, error) {
	// Initialize database
	database, err := NewDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
		return nil, err
	}
	defer database.Close() // Ensure closure of DB connection

	// Initialize SQLC Queries
	eventStore := New(database.DB)
	return eventStore, nil
}
