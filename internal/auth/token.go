package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/turanoo/bitebattle/pkg/config"
)

var jwtKey []byte

func InitJWTKey(cfg *config.Config) {
	jwtKey = []byte(cfg.Application.JWTSecret)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	expiration := time.Now().Add(72 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
