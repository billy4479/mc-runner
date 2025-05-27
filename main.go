package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.HideBanner = true
	e.Debug = true

	e.GET("/api/mc-hook", func(c echo.Context) error {
		for k, v := range c.QueryParams() {
			fmt.Printf("%s = %s\n", k, v[0])
		}

		return c.NoContent(http.StatusOK)
	})

	e.Start(":4479")
}
