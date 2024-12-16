package auth

import (
	"net/http"

	"github.com/rs/zerolog"
)

func getAuthApiKey(opts map[string]interface{}) Auth {
	apiKey := opts["apikey"].(string)

	return Auth{
		Eval: func(r *http.Request, l zerolog.Logger) bool {
			return r.Header.Get("X-Api-Key") == apiKey
		},
	}
}
