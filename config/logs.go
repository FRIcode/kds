package config

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger() {
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	switch Config.Logging.Level {
	case "trace":
		Logger = Logger.Level(zerolog.TraceLevel)
	case "debug":
		Logger = Logger.Level(zerolog.DebugLevel)
	case "info":
		Logger = Logger.Level(zerolog.InfoLevel)
	case "warn":
		Logger = Logger.Level(zerolog.WarnLevel)
	case "error":
		Logger = Logger.Level(zerolog.ErrorLevel)
	default:
		panic("Invalid log level")
	}
}
