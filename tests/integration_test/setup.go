//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const dbImage = "postgres:18"

func SetupTestDB(ctx context.Context) (*postgres.PostgresContainer, error) {
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)

	schemaFilePath := filepath.Join(testDir, "..", "..", "sql", "schema.sql")
	startupFilePath := filepath.Join(testDir, "..", "testdata", "init-db.sql")

	dbContainer, err := postgres.Run(
		ctx,
		dbImage,
		postgres.WithOrderedInitScripts(schemaFilePath, startupFilePath),
		postgres.WithDatabase("contacts_db"),
		postgres.WithUsername("root"),
		postgres.WithPassword("example"),
		postgres.BasicWaitStrategies(),
	)

	return dbContainer, err
}
