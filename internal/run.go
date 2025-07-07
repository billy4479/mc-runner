package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/billy4479/mc-runner/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
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
	e.Debug = Debug

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())

	if Debug {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:5173"},
		}))
	}

	e.Static("/", FrontendPath)

	fmt.Println("Frontend at", FrontendPath)

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

	err = e.Start(":4479")
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

	fmt.Println("warn: admin has no auth or invitation token, generating a new one")
	tokenB64, hash, err := generateRandomToken()
	if err != nil {
		return fmt.Errorf("generateRandomToken: %w", err)
	}
	fmt.Println(tokenB64)

	return repo.SetToken(ctx, repository.SetTokenParams{
		Token:   hash[:],
		Expires: time.Unix(0, 0),
		Type:    "invitation_token",
		UserID:  0,
	})
}
