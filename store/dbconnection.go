package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

// Database struct to manage the DB connection
type Database struct {
	Conn *pgx.Conn
}

// NewDatabase initializes a new database connection using pgx.Conn
func NewDatabase(dsn string) (*Database, error) {
	ctx := context.Background()

	// Create a direct connection
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// Check if the database is reachable
	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, err
	}

	return &Database{Conn: conn}, nil
}

// Close closes the database connection
func (d *Database) Close() {
	if d.Conn != nil {
		d.Conn.Close(context.Background())
	}
}
