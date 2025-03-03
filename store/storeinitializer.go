package store

import "log"

func InitializeEventStateStore(dsn string) (Querier, error) {
	// Initialize database connection
	//TODO: Should instantiate Event Store ones and use it to get new connection from the pool and connect to it
	database, err := NewDatabase(dsn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
		return nil, err
	}
	//defer database.Close() // Ensure closure of DB connection

	// Initialize SQLC Queries
	//TODO: Add strategy to remove older events
	eventStore := New(database.DB)
	return eventStore, nil
}
