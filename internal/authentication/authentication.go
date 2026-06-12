package authentication

import (
	config "core/internal/config"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/linkedin"
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

	gothic.Store = store
	goth.UseProviders(
		linkedin.New(
			cfg.ClientID,
			cfg.ClientSecret,
			cfg.CallbackURL,

			"openid",
			"profile",
			"email",
		),
	)
}