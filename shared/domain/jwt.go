package domain

import (
	"context"
	"errors"
	"log"
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
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func JwtGenerate(ctx context.Context, userId string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString(jwtSecret(ctx))
	if err != nil {
		return "", err
	}

	return token, nil
}

func JwtValidate(ctx context.Context, token string) (*JwtCustomClaim, error) {
	tokenParse, err := jwtParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there's a problem with the signing method")
		}
		return jwtSecret(ctx), nil
	})
	if err != nil || !tokenParse.Valid {
		if err != nil {
			log.Println("Error parse token:", err)
		}

		return nil, ErrInvalidToken
	}

	claims, ok := tokenParse.Claims.(*JwtCustomClaim)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}

func getJwtSecret(ctx context.Context) interface{} {
	jwtSecret, _ := ctx.Value(CtxJWTSecret).(string)
	return []byte(jwtSecret)
}
