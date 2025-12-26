package web

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/billy4479/mc-runner/internal/repository"
	"github.com/labstack/echo/v4"
)

func addAuthRoutes(g *echo.Group) {
	g.POST("/invite", func(c echo.Context) error {
		repo := c.Get("repo").(*repository.Queries)
		db := c.Get("db").(*sql.DB)
		ctx := c.Get("db_ctx").(context.Context)

		tokenB64, hash, err := generateRandomToken()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer tx.Rollback() // TODO: check this

		repoTx := repo.WithTx(tx)
		id, err := repoTx.CreateUser(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = repoTx.SetToken(ctx, repository.SetTokenParams{
			Token:   hash[:],
			Expires: time.Now().Add(time.Hour * 3 * 24),
			Type:    "invitation_token",
			UserID:  id,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = tx.Commit()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, echo.Map{"invitation_token": tokenB64})

	}, adminMiddleware)

	g.POST("/register", func(c echo.Context) error {
		repo := c.Get("repo").(*repository.Queries)
		ctx := c.Get("db_ctx").(context.Context)

		type ReqBody struct {
			InvitationToken string `json:"invitation_token"`
			Name            string `json:"name"`
		}

		var body ReqBody
		err := c.Bind(&body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		user, err := getUserFromTokenChecked(repo, ctx, body.InvitationToken, "invitation_token")
		if err != nil {
			return err
		}

		err = repo.SetUserName(ctx, repository.SetUserNameParams{
			ID: user.ID,
			Name: sql.NullString{
				String: body.Name,
				Valid:  true,
			},
		})

		if err != nil {
			if err.Error()[:6] == "UNIQUE" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user already exists: %w", err))
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		cookie, err := makeNewAuthToken(repo, ctx, user, "invitation_token")
		if err != nil {
			return err
		}

		c.SetCookie(cookie)

		return c.NoContent(http.StatusCreated)
	})

	g.POST("/addDevice", func(c echo.Context) error {
		repo := c.Get("repo").(*repository.Queries)
		ctx := c.Get("db_ctx").(context.Context)
		user := c.Get("user").(*repository.User)

		tokenB64, hash, err := generateRandomToken()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = repo.SetToken(ctx, repository.SetTokenParams{
			Token:   hash[:],
			Expires: time.Now().Add(1 * time.Hour),
			Type:    "login_token",
			UserID:  user.ID,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, echo.Map{"login_token": tokenB64})
	}, authMiddleware)

	g.POST("/login", func(c echo.Context) error {
		repo := c.Get("repo").(*repository.Queries)
		ctx := c.Get("db_ctx").(context.Context)

		type ReqBody struct {
			LoginToken string `json:"login_token"`
		}

		var body ReqBody
		err := c.Bind(&body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		user, err := getUserFromTokenChecked(repo, ctx, body.LoginToken, "login_token")
		if err != nil {
			return err
		}

		cookie, err := makeNewAuthToken(repo, ctx, user, "login_token")
		if err != nil {
			return err
		}

		c.SetCookie(cookie)

		return c.NoContent(http.StatusCreated)
	})

	g.GET("/me", func(c echo.Context) error {
		user := c.Get("user").(*repository.User)

		type UserForFrontend struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}

		name := "<null>"
		if user.Name.Valid {
			name = user.Name.String
		}

		return c.JSON(http.StatusOK, UserForFrontend{
			ID:   user.ID,
			Name: name,
		})
	}, authMiddleware)

	g.DELETE("/logout", func(c echo.Context) error {
		repo := c.Get("repo").(*repository.Queries)
		ctx := c.Get("db_ctx").(context.Context)
		user := c.Get("user").(*repository.User)

		cookie, err := c.Cookie("auth")
		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Errorf("at this point the cookie must exist: %w", err),
			)
		}

		hash, err := getHashFromB64(cookie.Value)
		if err != nil {
			return err
		}

		err = repo.RemoveTokenExact(ctx, repository.RemoveTokenExactParams{
			Token:  hash,
			UserID: user.ID,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     "auth",
			Secure:   true,
			HttpOnly: true,
			Value:    "",
			Path:     "/api",
			Expires:  time.Now(),
		})

		return c.JSON(http.StatusOK, echo.Map{})
	}, authMiddleware)
}
