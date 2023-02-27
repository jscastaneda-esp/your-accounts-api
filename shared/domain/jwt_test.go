package domain

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InRlc3QiLCJpYXQiOjE2Nzc1MTY1NDB9.-8dzHSaH9ZORt_RlGL5-dMkIKWaOjhq09lUHNQOda7w"

func TestJwtGenerateSuccess(t *testing.T) {
	require := require.New(t)

	token, err := JwtGenerate(context.Background(), "test")
	require.NoError(err)
	require.NotEmpty(token)
}

func TestJwtGenerateErrorKeyInvalid(t *testing.T) {
	require := require.New(t)

	originalJwtSecret := jwtSecret
	jwtSecret = func(ctx context.Context) interface{} {
		return ""
	}

	token, err := JwtGenerate(context.Background(), "test")
	require.EqualError(jwt.ErrInvalidKeyType, err.Error())
	require.Empty(token)

	jwtSecret = originalJwtSecret
}

func TestJwtValidateSuccess(t *testing.T) {
	require := require.New(t)

	claims, err := JwtValidate(context.Background(), token)
	require.NoError(err)
	require.NotNil(claims)
	require.Equal("test", claims.ID)
}

func TestJwtValidateErrorInvalidToken(t *testing.T) {
	require := require.New(t)

	originalJwtParseWithClaims := jwtParseWithClaims
	jwtParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{
			Claims: claims,
		}, nil
	}

	claims, err := JwtValidate(context.Background(), token)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(claims)

	jwtParseWithClaims = originalJwtParseWithClaims
}

func TestJwtValidateErrorParse(t *testing.T) {
	require := require.New(t)

	originalJwtSecret := jwtSecret
	jwtSecret = func(ctx context.Context) interface{} {
		return ""
	}

	claims, err := JwtValidate(context.Background(), token)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(claims)

	jwtSecret = originalJwtSecret
}

func TestJwtValidateErrorInvalidTokenClaims(t *testing.T) {
	require := require.New(t)

	originalJwtParseWithClaims := jwtParseWithClaims
	jwtParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{
			Valid:  true,
			Claims: jwt.MapClaims{},
		}, nil
	}

	claims, err := JwtValidate(context.Background(), token)
	require.EqualError(err, ErrInvalidTokenClaims.Error())
	require.Nil(claims)

	jwtParseWithClaims = originalJwtParseWithClaims
}
