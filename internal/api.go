package internal

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func addAPIRoutes(config *Config, api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	api.GET("/hooks/mc", func(c echo.Context) error {
		for k, v := range c.QueryParams() {
			fmt.Printf("%s = %s\n", k, v[0])
		}

		return c.NoContent(http.StatusOK)
	})
}
