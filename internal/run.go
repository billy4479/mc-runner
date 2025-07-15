package internal

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/billy4479/mc-runner/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var FrontendPath string = "frontend/dist"
var BuildMode string = "debug"

var Debug bool = true

func Run(config *Config) error {
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	e := echo.New()

	if BuildMode == "release" {
		Debug = false
	}

	e.HideBanner = true
	e.HidePort = true
	e.Debug = Debug

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			delta := time.Since(start)

			logger := log.Info()
			if err != nil {
				logger = log.Warn()
				if c.Response().Status == http.StatusInternalServerError {
					logger = log.Error()
				}
			}

			logger.
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Str("remote_ip", c.RealIP()).
				Int("status", c.Response().Status).
				Dur("duration", delta).
				Err(err).
				Send()

			return err
		}
	})

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			log.Warn().
				Str("where", "HTTPErrorHandler").
				Err(err).
				Msg("An error occurred but the response is already committed")
			return
		}

		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		msg := echo.Map{"status": code}
		if Debug {
			msg["error"] = err.Error()
		}
		err = c.JSON(code, msg)
		if err != nil {
			log.Error().Err(err).Str("where", "HTTPErrorHandler").Msg("Failed to return JSON response")
		}
	}

	if Debug {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{fmt.Sprintf("http://localhost:%d", config.VitePort)},
		}))
	}

	log.Info().Str("frontend_path", FrontendPath)
	e.Static("/", FrontendPath)

	repo := repository.New(db)
	ctx := context.TODO()

	api := e.Group("/api")
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("repo", repo)
			c.Set("db", db)
			c.Set("db_ctx", ctx)
			return next(c)
		}
	})

	addAPIRoutes(config, api)

	err = ensureAdminToken(repo, ctx)
	if err != nil {
		return fmt.Errorf("admin token: %w", err)
	}

	log.Info().Int("port", config.Port).Msg("Setup completed, starting the application")
	err = e.Start(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

func ensureAdminToken(repo *repository.Queries, ctx context.Context) error {
	tokens, err := repo.GetAllTokensForUser(ctx, 0)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, token := range tokens {
		if token.Expires.Unix() != 0 || token.Expires.After(time.Now()) {
			continue
		}

		if token.Type == "auth_token" || token.Type == "invitation_token" {
			return nil
		}
	}

	log.Warn().
		Str("where", "ensureAdminToken").
		Msg("admin has no auth or invitation token, generating a new one")

	tokenB64, hash, err := generateRandomToken()
	if err != nil {
		return fmt.Errorf("generateRandomToken: %w", err)
	}

	log.Info().
		Str("admin_token", tokenB64).Send()

	return repo.SetToken(ctx, repository.SetTokenParams{
		Token:   hash[:],
		Expires: time.Unix(0, 0),
		Type:    "invitation_token",
		UserID:  0,
	})
}
