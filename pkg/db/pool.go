package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// CreateDBConnectionPool creates a new PostgreSQL connection pool.
func CreateDBConnectionPool() (*pgxpool.Pool, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	return pgxpool.Connect(context.Background(), connStr)
}
