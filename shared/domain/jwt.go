package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CtxJWTSecret = jwtKey("jwtSecret")

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

	token, err := t.SignedString(getJwtSecret(ctx))
	if err != nil {
		return "", err
	}

	return token, nil
}

func JwtValidate(ctx context.Context, token string) (*JwtCustomClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}

		return getJwtSecret(ctx), nil
	})
	if err != nil || !tokenParse.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := tokenParse.Claims.(*JwtCustomClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func getJwtSecret(ctx context.Context) []byte {
	jwtSecret, _ := ctx.Value(CtxJWTSecret).(string)
	return []byte(jwtSecret)
}
