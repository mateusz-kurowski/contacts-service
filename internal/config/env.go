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
func NewEnv(dbURL string, isTestEnv bool) (*Env, error) {
	ctx := context.Background()
	conn, connErr := pgxpool.New(ctx, dbURL)
	if connErr != nil {
		return nil, connErr
	}

	//nolint:exhaustruct // necessary because of bucket's conditional occurrence
	env := Env{Logger: slog.Default()}

	// Test the connection
	if pingErr := conn.Ping(ctx); pingErr != nil {
		return nil, pingErr
	}
	env.Queries = db.New(conn)

	if !isTestEnv {
		bucket, err := bucket.OpenFromEnv(ctx)
		if err != nil {
			return nil, err
		}
		env.Bucket = bucket
	}
	return &env, nil
}
