package cmd

import (
	"fmt"
	"os"

	"github.com/billy4479/mc-runner/internal"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() error {
	noDotenv := os.Getenv("DONT_LOAD_DOTENV")

	var dotenvErr error = nil
	if noDotenv != "yes" {
		dotenvErr = godotenv.Load()
		if dotenvErr != nil {
			dotenvErr = fmt.Errorf("load .env: %w", dotenvErr)
		}
	}

	config := internal.Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Environment == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Debug().Msg("logger is ready")

	if noDotenv == "yes" {
		log.Debug().Str("DONT_LOAD_DOTENV", noDotenv).Msg("will not load .env")
	}

	if dotenvErr != nil {
		log.Debug().Err(err).Msg("proceeding without .env")
	}

	return internal.Run(&config)
}
