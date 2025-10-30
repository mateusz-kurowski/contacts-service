//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/db"
	"contactsAI/contacts/internal/handlers"
	"contactsAI/contacts/internal/routing"
	integration "contactsAI/contacts/tests/integration_test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func setupSuite(t *testing.T) (*gin.Engine, func(t *testing.T)) {
	ctx := context.Background()
	dbContainer, err := integration.SetupTestDB(ctx)
	require.NoError(t, err, "testcontainer creation failed")

	dbURL, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get conn string")

	env, err := config.NewEnv(dbURL)
	require.NoError(t, err, "db connection failed")

	router := routing.SetupRouter(env)

	return router, func(t *testing.T) {
		testcontainers.CleanupContainer(t, dbContainer)
	}
}

func TestIntegration(t *testing.T) {
	router, teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	t.Run("GET /api/contacts", func(t *testing.T) {
		w := integration.MkJSONRequest(t, "GET", "/api/contacts/", router, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var response []handlers.ContactResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Unmarshaling response failed")

		assert.Len(t, response, 8)

		first := response[0]
		assert.Equal(t, "Agnieszka Szymańska", first.Name)
		assert.Equal(t, "888-999-000", first.Phone)
	})

	t.Run("GET /api/contacts/:id", func(t *testing.T) {
		w := integration.MkGetContactByIDRequest(t, 2, router)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ContactResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Unmarshaling response failed")

		assert.Equal(t, "Anna Nowak", response.Name)
		assert.Equal(t, "987-654-321", response.Phone)
	})

	t.Run("POST /api/contacts with dial", func(t *testing.T) {
		newContactBody := db.CreateContactParams{
			Name:  "testName",
			Phone: "+48 123 123 123",
		}

		w := integration.MkJSONRequest(t, "POST", "/api/contacts/", router, newContactBody)

		assert.Exactly(t, http.StatusCreated, w.Code)

		var response handlers.ContactResponse

		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Unmarshaling response failed")

		assert.Exactly(t, newContactBody.Name, response.Name)
		assert.Exactly(t, newContactBody.Phone, response.Phone)
		assert.NotNil(t, response.ID)
	})

	t.Run("PUT /api/contacts", func(t *testing.T) {
		contactID := 1
		updateBody := handlers.UpdateContactBody{
			Name:  "newname",
			Phone: "+48123456789",
		}

		w := integration.MkJSONRequest(t, "PUT", fmt.Sprintf("/api/contacts/%d", contactID), router, updateBody)
		assert.Equal(t, http.StatusOK, w.Code)

		var response handlers.ContactResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

		assert.Exactly(t, updateBody.Name, response.Name)
		assert.Exactly(t, updateBody.Phone, response.Phone)
		assert.Exactly(t, int32(contactID), response.ID)

		wUpdatedContact := integration.MkGetContactByIDRequest(t, contactID, router)

		assert.Equal(t, http.StatusOK, wUpdatedContact.Code)

		var updatedContactResponse handlers.ContactResponse
		assert.NoError(t, json.Unmarshal(wUpdatedContact.Body.Bytes(), &updatedContactResponse))

		assert.Exactly(t, updateBody.Name, updatedContactResponse.Name)
		assert.Exactly(t, updateBody.Phone, updatedContactResponse.Phone)
		assert.Exactly(t, int32(contactID), updatedContactResponse.ID)
	})

	t.Run("DELETE /api/contacts", func(t *testing.T) {
		contactID := 2

		wDel := integration.MkJSONRequest(t, "DELETE",
			fmt.Sprintf("/api/contacts/%d", contactID), router, nil)
		assert.Equal(t, http.StatusNoContent, wDel.Code)

		wAfterDel := integration.MkGetContactByIDRequest(t, contactID, router)
		assert.Exactly(t, http.StatusNotFound, wAfterDel.Code)
	})
}
