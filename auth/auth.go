package auth

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/FRIcode/kds/config"
)

type Auth struct {
	Eval func(r *http.Request, l zerolog.Logger) bool
}

func GetAuth(cfg config.AuthConfig) Auth {
	switch cfg.Type {
	case "none":
		return getAuthNone(cfg.Opts)
	case "apikey":
		return getAuthApiKey(cfg.Opts)
	case "jwt":
		return getAuthJWT(cfg.Opts)
	default:
		panic("Unknown auth type")
	}
}
