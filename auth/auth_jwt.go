package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/cel-go/cel"
	"github.com/rs/zerolog"

	"github.com/FRIcode/kds/config"
)

func getAuthJWT(opts map[string]interface{}) Auth {
	jwks := opts["jwks"].(string)
	expr := opts["expr"].(string)
	aud := opts["aud"].(string)

	log := config.Logger.With().Str("service", "auth").Logger()
	log.Trace().Msgf("JWKS: %s", jwks)
	log.Trace().Msgf("Expression: %s", expr)
	log.Trace().Msgf("Audience: %s", aud)

	// Validate expression
	env, err := cel.NewEnv(
		cel.Variable("token", cel.AnyType),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create environment")
		panic(err)
	}

	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		log.Error().Err(iss.Err()).Msg("Failed to compile expression")
		panic(iss.Err())
	}

	return Auth{
		Eval: func(r *http.Request, l zerolog.Logger) bool {
			log := l.With().Str("service", "auth").Logger()

			// Validate JWT token
			authHeader := r.Header.Get("Authorization")
			log.Trace().Msgf("Authorization header: %s", authHeader)
			if authHeader == "" || len(authHeader) < len("Bearer ") {
				log.Warn().Msg("Missing or invalid Authorization header")
				return false
			}

			tokenString := authHeader[len("Bearer "):]
			jwks, err := keyfunc.NewDefault([]string{jwks})
			if err != nil {
				log.Error().Err(err).Msg("Failed to create JWKS keyfunc")
				return false
			}

			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, jwks.Keyfunc)
			if err != nil || !token.Valid {
				log.Warn().Err(err).Msg("Invalid token")
				return false
			}

			tokenaud, err := token.Claims.GetAudience()
			log.Trace().Msgf("Token audience: %v", tokenaud)
			if err != nil || len(tokenaud) != 1 || tokenaud[0] != aud {
				log.Warn().Err(err).Msg("Invalid audience")
				return false
			}

			// Get data from token
			data := map[string]interface{}{}
			tokenDataBase64 := strings.Split(tokenString, ".")[1]

			tokenData, err := base64.RawStdEncoding.DecodeString(tokenDataBase64)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to decode token data")
				return false
			}

			err = json.Unmarshal(tokenData, &data)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to unmarshal token data")
				return false
			}

			// Evaluate expression
			prg, err := env.Program(ast)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create program")
				return false
			}

			out, _, err := prg.Eval(map[string]any{"token": data})
			if err != nil {
				log.Error().Err(err).Msg("Failed to evaluate program")
				return false
			}

			if out.Type() != cel.BoolType {
				log.Warn().Msg("Invalid program output")
				return false
			}

			log.Trace().Msgf("Program output: %v", out.Value())
			return out.Value().(bool)
		},
	}
}
