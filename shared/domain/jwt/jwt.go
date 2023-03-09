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
	UUID  string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func JwtGenerate(ctx context.Context, id string, uuid string, email string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		UUID:  uuid,
		Email: email,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString(jwtSecret(ctx))
	if err != nil {
		return "", err
	}

	return token, nil
}

func getJwtSecret(ctx context.Context) interface{} {
	jwtSecret, _ := ctx.Value(CtxJWTSecret).(string)
	return []byte(jwtSecret)
}