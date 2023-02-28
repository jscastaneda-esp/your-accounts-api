package domain

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	token                      string
	originalJwtSecret          func(ctx context.Context) interface{}
	originalJwtParseWithClaims func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InRlc3QiLCJpYXQiOjE2Nzc1MTY1NDB9.-8dzHSaH9ZORt_RlGL5-dMkIKWaOjhq09lUHNQOda7w"
	suite.originalJwtSecret = jwtSecret
	suite.originalJwtParseWithClaims = jwtParseWithClaims
}

func (suite *TestSuite) SetupTest() {
	jwtSecret = suite.originalJwtSecret
	jwtParseWithClaims = suite.originalJwtParseWithClaims
}

func (suite *TestSuite) TestJwtGenerateSuccess() {
	require := require.New(suite.T())

	token, err := JwtGenerate(context.Background(), "test")
	require.NoError(err)
	require.NotEmpty(token)
}

func (suite *TestSuite) TestJwtGenerateErrorKeyInvalid() {
	require := require.New(suite.T())

	jwtSecret = func(ctx context.Context) interface{} {
		return ""
	}

	token, err := JwtGenerate(context.Background(), "test")
	require.EqualError(jwt.ErrInvalidKeyType, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestJwtValidateSuccess() {
	require := require.New(suite.T())

	claims, err := JwtValidate(context.Background(), suite.token)
	require.NoError(err)
	require.NotNil(claims)
	require.Equal("test", claims.ID)
}

func (suite *TestSuite) TestJwtValidateErrorInvalidToken() {
	require := require.New(suite.T())

	jwtParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{
			Claims: claims,
		}, nil
	}

	claims, err := JwtValidate(context.Background(), suite.token)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(claims)
}

func (suite *TestSuite) TestJwtValidateErrorParse() {
	require := require.New(suite.T())

	jwtSecret = func(ctx context.Context) interface{} {
		return ""
	}

	claims, err := JwtValidate(context.Background(), suite.token)
	require.EqualError(err, ErrInvalidToken.Error())
	require.Nil(claims)
}

func (suite *TestSuite) TestJwtValidateErrorInvalidTokenClaims() {
	require := require.New(suite.T())

	jwtParseWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{
			Valid:  true,
			Claims: jwt.MapClaims{},
		}, nil
	}

	claims, err := JwtValidate(context.Background(), suite.token)
	require.EqualError(err, ErrInvalidTokenClaims.Error())
	require.Nil(claims)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
