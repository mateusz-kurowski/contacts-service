package config

import (
	"context"

	"contactsAI/contacts/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	*db.Queries
}

// NewEnv Create a new Env instance.
func NewEnv(dbURL string) (*Env, error) {
	ctx := context.Background()
	conn, connErr := pgxpool.New(ctx, dbURL)
	if connErr != nil {
		return nil, connErr
	}

	// Test the connection
	if pingErr := conn.Ping(ctx); pingErr != nil {
		return nil, pingErr
	}

	return &Env{Queries: db.New(conn)}, nil
}
