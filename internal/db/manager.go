package db

import (
	"context"
	"fmt"

	"github.com/hydrocode-de/datailama/internal/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Manager handles database connections and queries
type Manager struct {
	*sql.Queries
	pool *pgxpool.Pool
}

// New creates a new database manager
func New(ctx context.Context, connString string) (*Manager, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	queries := sql.New(pool)
	return &Manager{
		Queries: queries,
		pool:    pool,
	}, nil
}

// Close closes the database connection pool
func (m *Manager) Close() {
	if m.pool != nil {
		m.pool.Close()
	}
}

// GetPool returns the underlying connection pool
func (m *Manager) GetPool() *pgxpool.Pool {
	return m.pool
}
