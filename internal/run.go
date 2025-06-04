package internal

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		e.Use(middleware.CORS())
	}

	e.Static("/", FrontendPath)

	fmt.Println("Frontend at", FrontendPath)

	api := e.Group("/api")
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	addAPIRoutes(config, api)

	err = e.Start(":4479")
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}
