package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"syscall"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/billy4479/mc-runner/internal/driver"
	"github.com/billy4479/mc-runner/internal/web"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run(frontend fs.FS) error {
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

	d, err := driver.NewDriver(conf)
	if err != nil {
		err = fmt.Errorf("driver: %w", err)
		log.Fatal().Err(err)
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func(d *driver.Driver) {
		// TODO: stop web server gracefully
		<-sigs
		log.Info().Msg("received signal, terminating gracefully")
		d.Stop()
		os.Exit(0)
	}(d)

	err = web.RunWebServer(conf, d, frontend)
	if err != nil {
		err = fmt.Errorf("web: %w", err)
		log.Fatal().Err(err)
		return err
	}

	return nil
}
