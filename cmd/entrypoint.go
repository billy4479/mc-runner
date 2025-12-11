package cmd

import (
	"fmt"
	"os"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/billy4479/mc-runner/internal/driver"
	"github.com/billy4479/mc-runner/internal/web"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() error {
	conf, dotenvErr, err := config.NewConfig()

	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(conf.GetLogLevel())
	if conf.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Debug().Msg("logger is ready")

	if dotenvErr != nil {
		log.Warn().Err(err).Msg("proceeding without .env")
	}

	driver, err := driver.NewDriver(conf)
	if err != nil {
		err = fmt.Errorf("driver: %w", err)
		log.Fatal().Err(err)
		return err
	}

	err = web.RunWebServer(conf, driver)
	if err != nil {
		err = fmt.Errorf("web: %w", err)
		log.Fatal().Err(err)
		return err
	}

	return nil
}
