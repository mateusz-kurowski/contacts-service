package config

import (
	"context"
	"log/slog"

	"contactsAI/contacts/internal/bucket"
	"contactsAI/contacts/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	*db.Queries
	*slog.Logger

	Bucket *bucket.Store
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
	bucket, err := bucket.OpenFromEnv(ctx)
	if err != nil {
		return nil, err
	}
	logger := slog.Default()
	return &Env{Queries: db.New(conn), Bucket: bucket, Logger: logger}, nil
}
