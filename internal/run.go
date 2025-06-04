package internal

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var FrontendPath string = "frontend/dist"
var BuildMode string = "debug"

var Debug bool = true

func Run(config *Config) error {
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
	addAPIRoutes(config, api)

	err := e.Start(":4479")
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}
