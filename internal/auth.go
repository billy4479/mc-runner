package internal

import (
	"context"
	"crypto/rand"
	"crypto/sha3"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/billy4479/mc-runner/repository"
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

		token, err := base64.RawURLEncoding.DecodeString(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		hash := sha3.Sum256(token)
		user, err := getUserFromTokenChecked(repo, ctx, hash[:], "auth_token")
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

func generateRandomToken() (tokenB64 string, hash [32]byte, err error) {
	token := make([]byte, 32) // 256 bits
	_, err = rand.Read(token)
	if err != nil {
		return
	}

	hash = sha3.Sum256(token)
	tokenB64 = base64.RawURLEncoding.EncodeToString(token)
	fmt.Println(hex.EncodeToString(token))
	return
}

var (
	ErrExpiredToken = errors.New("token expired")
)

func getUserFromTokenChecked(repo *repository.Queries, ctx context.Context, hash []byte, tokenType string) (*repository.User, error) {
	userAndExpiration, err := repo.GetUserFromToken(ctx, repository.GetUserFromTokenParams{
		Token: hash,
		Type:  tokenType,
	})

	if err != nil {
		return nil, err
	}

	if userAndExpiration.Expires.Unix() == 0 {
		return &userAndExpiration.User, nil
	}

	if userAndExpiration.Expires.After(time.Now()) {
		return nil, fmt.Errorf("%w: %s", ErrExpiredToken, tokenType)
	}

	return &userAndExpiration.User, nil
}

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

		token, err := base64.RawURLEncoding.DecodeString(body.InvitationToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		hash := sha3.Sum256(token)
		user, err := getUserFromTokenChecked(repo, ctx, hash[:], "invitation_token")

		if err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if errors.Is(err, ErrExpiredToken) {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = repo.SetUserName(ctx, repository.SetUserNameParams{
			ID: user.ID,
			Name: sql.NullString{
				String: body.Name,
				Valid:  true,
			},
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tokenB64, hash, err := generateRandomToken()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = repo.SetToken(ctx, repository.SetTokenParams{
			Token:   hash[:],
			Expires: time.Unix(0, 0),
			Type:    "auth_token",
			UserID:  user.ID,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		c.SetCookie(&http.Cookie{
			Name:     "auth",
			Secure:   true,
			HttpOnly: true,
			Value:    tokenB64,
			Path:     "/api",
			Expires:  time.Now().Add(400 * 24 * time.Hour), // This is the maximum expiration time
		})

		err = repo.RemoveTokenById(ctx, repository.RemoveTokenByIdParams{
			UserID: user.ID,
			Type:   "invitation_token",
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusCreated, echo.Map{"auth_token": tokenB64})
	})

	g.GET("/me", func(c echo.Context) error {
		user := c.Get("user").(*repository.User)
		return c.JSON(http.StatusOK, user)
	}, authMiddleware)
}
