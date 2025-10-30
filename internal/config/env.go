package config

import (
	"context"
	"net/http"
	"os"

	"contactsAI/contacts/internal/db"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth/gothic"
)

type Env struct {
	*db.Queries
	CookieStore cookie.Store
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

	cookieStore := setupSessionStorage()
	gothic.Store = cookieStore
	return &Env{Queries: db.New(conn), CookieStore: cookieStore}, nil
}

func setupSessionStorage() cookie.Store {
	sessionSecret, sessionKey := []byte(os.Getenv("SESSION_SECRET")), []byte(os.Getenv("SESSION_KEY"))

	cookieStore := cookie.NewStore(sessionSecret, sessionKey)

	cookieStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,      // Prevent XSS attacks
		Secure:   true,      // Use only with HTTPS in production
		SameSite: http.SameSiteLaxMode,
	})

	return cookieStore
}
