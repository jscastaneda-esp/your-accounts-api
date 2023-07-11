package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CtxJWTSecret = jwtKey("jwtSecret")

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")

	jwtParseWithClaims = jwt.ParseWithClaims
	jwtSecret          = getJwtSecret
)

type jwtKey string

type JwtCustomClaim struct {
	UID   string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func JwtGenerate(ctx context.Context, id string, uid string, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		UID:   uid,
		Email: email,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString(jwtSecret(ctx))
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func getJwtSecret(ctx context.Context) any {
	jwtSecret, _ := ctx.Value(CtxJWTSecret).(string)
	return []byte(jwtSecret)
}
