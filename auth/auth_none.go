package auth

import (
	"net/http"

	"github.com/rs/zerolog"
)

func getAuthNone(_ map[string]interface{}) Auth {
	return Auth{
		Eval: func(r *http.Request, l zerolog.Logger) bool {
			return true
		},
	}
}
