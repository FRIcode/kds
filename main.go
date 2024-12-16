package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/FRIcode/kds/config"
	"github.com/FRIcode/kds/deployment"
)

func main() {
	config.InitConfig()
	config.InitLogger()

	logger := config.Logger.With().Str("module", "main").Logger()
	logger.Info().Msg("Starting server")

	for _, deploy := range config.Config.Deployments {
		logger.Debug().Str("name", deploy.Name).Msg("Starting deployment")
		http.HandleFunc("/api/v1/deploy/"+deploy.Name, deployment.DeploymentHandler(deploy))
	}

	http.HandleFunc("/api/v1/status/{id}", deployment.StatusHandler)

	http.Handle("/metrics", promhttp.Handler())

	logger.Info().Str("host", config.Config.Server.Host).Msg("Listening")
	http.ListenAndServe(config.Config.Server.Host, nil)
}
