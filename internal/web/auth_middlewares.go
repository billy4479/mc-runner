package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/billy4479/mc-runner/internal/repository"
	"github.com/labstack/echo/v4"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		repo := c.Get("repo").(*repository.Queries)
		ctx := c.Get("db_ctx").(context.Context)

		user, err := getUserFromTokenChecked(repo, ctx, cookie.Value, "auth_token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		// Set again the cookie expiration time to the maximum possible
		c.SetCookie(&http.Cookie{
			Name:     "auth",
			Secure:   true,
			HttpOnly: true,
			Value:    cookie.Value,
			Expires:  time.Now().Add(400 * 24 * time.Hour),
			Path:     "/api",
		})

		c.Set("user", user)
		return next(c)
	}
}

func adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return authMiddleware(
		func(c echo.Context) error {
			user := c.Get("user").(*repository.User)
			if user.ID != 0 {
				return echo.NewHTTPError(http.StatusForbidden,
					fmt.Errorf("user_id %d is not allowed", user.ID))
			}
			return next(c)
		})
}
