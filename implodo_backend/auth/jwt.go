package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"implodo_backend/config"
)

type Claims struct {
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func issueJWT(cfg *config.Config, googleID, email, name string) (string, error) {
	now := time.Now()
	claims := Claims{
		GoogleID: googleID,
		Email:    email,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "implodo_backend",
			Subject:   googleID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.JWTDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func parseJWT(cfg *config.Config, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
