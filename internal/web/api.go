package web

import (
	"fmt"
	"net/http"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/billy4479/mc-runner/internal/driver"
	"github.com/labstack/echo/v4"
)

func addAPIRoutes(conf *config.Config, api *echo.Group, driver *driver.Driver) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	auth := api.Group("/auth")
	addAuthRoutes(auth)

	addWebsocket(api, conf, driver)

	api.GET("/hooks/mc", func(c echo.Context) error {
		for k, v := range c.QueryParams() {
			fmt.Printf("%s = %s\n", k, v[0])
		}

		return c.NoContent(http.StatusOK)
	})
}
