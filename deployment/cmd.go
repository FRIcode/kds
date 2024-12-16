package deployment

import (
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/FRIcode/kds/config"
	"github.com/FRIcode/kds/metrics"
)

func runDeployment(deployment config.DeploymentConfig, uid string, l zerolog.Logger) string {
	metrics.AddStatusEntry(metrics.StatusEntry{
		ID:     uid,
		Status: metrics.StatusRunning,
	})

	metrics.MetricDeploymentsTotal.WithLabelValues(deployment.Name).Inc()
	started := time.Now()

	go func() {
		log := l.With().Str("service", "runner").Str("name", deployment.Name).Logger()
		log.Info().Msg("Running deployment")

		failed := false
		for i, script := range deployment.Script {
			log.Debug().Int("index", i).Str("script", script).Msg("Running script")

			cmd := exec.Command("sh", "-c", script)
			cmd.Dir = deployment.WorkDir

			for _, env := range deployment.Env {
				cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
			}

			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Error().Err(err).Int("index", i).Msg("Failed to run script")
				failed = true
				break
			}

			log.Info().Int("index", i).Str("script", script).Str("output", strings.TrimSpace(string(out))).Msg("Script output")
		}

		duration := time.Since(started)
		metrics.MetricDeploymentsDuration.WithLabelValues(deployment.Name).Observe(duration.Seconds())

		if failed {
			metrics.UpdateStatusEntry(uid, metrics.StatusFailed)
			log.Info().Msg("Deployment failed")
			metrics.MetricDeploymentsFailed.WithLabelValues(deployment.Name).Inc()
		} else {
			metrics.UpdateStatusEntry(uid, metrics.StatusSuccess)
			log.Info().Msg("Deployment finished")
			metrics.MetricDeploymentsSuccess.WithLabelValues(deployment.Name).Inc()
		}
	}()

	return uid
}
