package jwt

import (
	"errors"
	"time"
	"your-accounts-api/shared/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

const CtxJWTSecret = jwtKey("jwtSecret")

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)

type jwtKey string

type JwtUserClaims struct {
	ID    uint   `json:"id"`
	UID   string `json:"uid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func JwtGenerate(id uint, uid string, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtUserClaims{
		ID:    id,
		UID:   uid,
		Email: email,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString([]byte(config.JWT_SECRET))
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}
