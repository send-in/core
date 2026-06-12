package authentication

import (
	config "core/internal/config"
	logger "core/pkg/log"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

func Configure(cfg *config.AuthConfig) {
	store := sessions.NewCookieStore(
		[]byte(cfg.SessionSecret),
		[]byte(cfg.EncryptionSecret),
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false,
	}

	provider, err := openidConnect.NewNamed(
		"linkedin",
		cfg.ClientID,
		cfg.ClientSecret,
		cfg.CallbackURL,
		"https://www.linkedin.com/oauth/.well-known/openid-configuration",
		"openid",
		"profile",
		"email",
	)


	if err != nil {
		logger.Error("Failed to start auth %s", err)
	}

	gothic.Store = store
	goth.UseProviders(provider)
}