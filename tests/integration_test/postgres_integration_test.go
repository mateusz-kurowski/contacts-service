//go:build integration

package integration_test

import (
	"context"
	"testing"

	integration "contactsAI/contacts/tests/integration_test"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()

	dbContainer, err := integration.SetupTestDB(ctx)
	testcontainers.CleanupContainer(t, dbContainer)
	require.NoError(t, err, "Setup db failed. Check if docker is available")

	connStr, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	t.Logf("Connection string: %s", connStr)
}
