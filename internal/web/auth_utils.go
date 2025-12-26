package web

import (
	"context"
	"crypto/rand"
	"crypto/sha3"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/billy4479/mc-runner/internal/repository"
	"github.com/labstack/echo/v4"
)

func generateRandomToken() (tokenB64 string, hash [32]byte, err error) {
	token := make([]byte, 32) // 256 bits
	_, err = rand.Read(token)
	if err != nil {
		return
	}

	hash = sha3.Sum256(token)
	tokenB64 = base64.RawURLEncoding.EncodeToString(token)
	return
}

var (
	ErrExpiredToken = errors.New("token expired")
)

func getUserFromTokenChecked(repo *repository.Queries, ctx context.Context, tokenB64 string, tokenType string) (*repository.User, error) {
	hash, err := getHashFromB64(tokenB64)
	if err != nil {
		return nil, err
	}

	userAndExpiration, err := repo.GetUserFromToken(ctx, repository.GetUserFromTokenParams{
		Token: hash,
		Type:  tokenType,
	})

	if err == sql.ErrNoRows {
		return nil, echo.NewHTTPError(http.StatusNotFound, err)
	}

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if userAndExpiration.Expires.Unix() == 0 {
		return &userAndExpiration.User, nil
	}

	if userAndExpiration.Expires.Before(time.Now()) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("%w: %s", ErrExpiredToken, tokenType))
	}

	return &userAndExpiration.User, nil
}

func getHashFromB64(b64 string) ([]byte, error) {
	token, err := base64.RawURLEncoding.DecodeString(b64)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	hash := sha3.Sum256(token)
	return hash[:], nil
}

func makeNewAuthToken(repo *repository.Queries, ctx context.Context, user *repository.User, tokenToRemove string) (*http.Cookie, error) {
	tokenB64, hash, err := generateRandomToken()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = repo.SetToken(ctx, repository.SetTokenParams{
		Token:   hash[:],
		Expires: time.Unix(0, 0),
		Type:    "auth_token",
		UserID:  user.ID,
	})
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cookie := &http.Cookie{
		Name:     "auth",
		Secure:   true,
		HttpOnly: true,
		Value:    tokenB64,
		Path:     "/api",
		Expires:  time.Now().Add(400 * 24 * time.Hour), // This is the maximum expiration time
	}

	if tokenToRemove != "" {
		err = repo.RemoveTokenById(ctx, repository.RemoveTokenByIdParams{
			UserID: user.ID,
			Type:   tokenToRemove,
		})
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return cookie, nil
}
