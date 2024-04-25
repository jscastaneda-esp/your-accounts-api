package domain

import "github.com/golang-jwt/jwt/v5"

type JwtUserClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}
