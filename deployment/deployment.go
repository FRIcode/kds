package deployment

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"

	"github.com/FRIcode/kds/auth"
	"github.com/FRIcode/kds/config"
	"github.com/FRIcode/kds/metrics"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomIdentifier() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func DeploymentHandler(deployment config.DeploymentConfig) http.HandlerFunc {
	l := config.Logger.With().Str("service", "api").Str("name", deployment.Name).Logger()
	auth := auth.GetAuth(deployment.Auth)

	return func(w http.ResponseWriter, r *http.Request) {
		uid := randomIdentifier()
		log := l.With().Str("id", uid).Logger()
		log.Info().Msg("Received request")
		metrics.MetricRequestsTotal.WithLabelValues(r.URL.Path).Inc()

		if !auth.Eval(r, log) {
			log.Warn().Str("addr", r.RemoteAddr).Str("x-forwarded-for", r.Header.Get("X-Forwarded-For")).Msg("Unauthorized")
			metrics.MetricRequestsAuthFailed.WithLabelValues(r.URL.Path).Inc()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		metrics.MetricRequestsAuthSuccess.WithLabelValues(r.URL.Path).Inc()

		runDeployment(deployment, uid, log)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": uid})
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log := config.Logger.With().Str("service", "api").Logger()
	log.Info().Msg("Received request")

	path := strings.Split(r.URL.Path, "/")
	path = path[:len(path)-1]
	metrics.MetricRequestsTotal.WithLabelValues(strings.Join(path, "/")).Inc()

	id := r.PathValue("id")
	if id == "" {
		log.Warn().Msg("Missing id")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		return
	}

	s := metrics.GetStatusEntry(id)
	if s == nil {
		log.Warn().Str("id", id).Msg("Status not found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}
