package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configures the global zerolog logger.
func Init(env string) {
	zerolog.TimeFieldFormat = time.RFC3339
	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger()
	}
}

func L() *zerolog.Logger {
	return &log.Logger
}
