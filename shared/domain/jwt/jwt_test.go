package jwt

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestJwtGenerateSuccess() {
	require := require.New(suite.T())

	token, _, err := JwtGenerate(1, "test", "test")

	require.NoError(err)
	require.NotEmpty(token)
}

func (suite *TestSuite) TestGetUserDataSuccess() {
	require := require.New(suite.T())
	token := &jwt.Token{
		Claims: &JwtUserClaims{
			ID: 1,
		},
	}
	app := fiber.New()
	c := app.AcquireCtx(new(fasthttp.RequestCtx))
	c.Locals("user", token)

	userData := GetUserData(c)

	require.NotNil(userData)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
