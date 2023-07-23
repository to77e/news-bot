// Package database - provides database connection and database operations.
package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres - creates new connection to postgres database.
func NewPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool connections: %w", err)
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire connection from the pool: %w", err)
	}
	defer conn.Release()

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping the database: %w", err)
	}

	return pool, nil
}
