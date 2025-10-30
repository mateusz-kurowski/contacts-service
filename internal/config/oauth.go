package config

import (
	"fmt"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
)

func SetupProviders() {
	useAuthentik()
}

func useAuthentik() error {
	openidConnect, err := openidConnect.New(
		os.Getenv("AUTHENTIK_CLIENT_ID"),
		os.Getenv("AUTHENTIK_CLIENT_SECRET"),
		"http://localhost:33500/api/auth/callback?provider=openid-connect",
		os.Getenv("AUTHENTIK_DISCOVERY_URL"),
		"openid", "profile", "email", // Required scopes
	)
	if err != nil {
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	goth.UseProviders(openidConnect)
	return nil
}
