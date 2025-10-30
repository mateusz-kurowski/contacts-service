//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func MkGetContactByIDRequest(t *testing.T, id int, router *gin.Engine) *httptest.ResponseRecorder {
	return MkJSONRequest(t, "GET", fmt.Sprintf("/api/contacts/%d", id), router, nil)
}

func MkJSONRequest(t *testing.T, method, path string, router *gin.Engine, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	w := httptest.NewRecorder()

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	return w
}
