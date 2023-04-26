package jwt

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	originalJwtSecret          func(ctx context.Context) any
	originalJwtParseWithClaims func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc, options ...jwt.ParserOption) (*jwt.Token, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.originalJwtSecret = jwtSecret
	suite.originalJwtParseWithClaims = jwtParseWithClaims
}

func (suite *TestSuite) SetupTest() {
	jwtSecret = suite.originalJwtSecret
	jwtParseWithClaims = suite.originalJwtParseWithClaims
}

func (suite *TestSuite) TestJwtGenerateSuccess() {
	require := require.New(suite.T())

	token, err := JwtGenerate(context.Background(), "test", "test", "test")

	require.NoError(err)
	require.NotEmpty(token)
}

func (suite *TestSuite) TestJwtGenerateErrorKeyInvalid() {
	require := require.New(suite.T())

	jwtSecret = func(ctx context.Context) any {
		return ""
	}

	token, err := JwtGenerate(context.Background(), "test", "test", "test")

	require.EqualError(jwt.ErrInvalidKeyType, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
