package auth

import (
	"api-your-accounts/shared/domain"
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	originalJwtValidate func(ctx context.Context, token string) (*domain.JwtCustomClaim, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.originalJwtValidate = jwtValidate
}

func (suite *TestSuite) SetupTest() {
	jwtValidate = suite.originalJwtValidate
}

func (suite *TestSuite) TestNewHandlerSuccess() {
	require := require.New(suite.T())

	body := "Hello, World!"

	app := fiber.New()
	app.Use(New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(body)
	})

	req := httptest.NewRequest("GET", "/", nil)

	res, _ := app.Test(req, 1)

	require.Equal(fiber.StatusOK, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(body, string(bytes))
}

func (suite *TestSuite) TestNewHandlerSuccessConfigNext() {
	require := require.New(suite.T())

	body := "Hello, World!"

	app := fiber.New()
	app.Use(New(Config{
		Next: func(c *fiber.Ctx) bool {
			return true
		},
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(body)
	})

	req := httptest.NewRequest("GET", "/", nil)

	res, _ := app.Test(req, 1)

	require.Equal(fiber.StatusOK, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(body, string(bytes))
}

func (suite *TestSuite) TestNewHandlerSuccessValidToken() {
	require := require.New(suite.T())

	body := "Hello, World!"

	app := fiber.New()
	app.Use(New(Config{}))
	app.Get("/", func(c *fiber.Ctx) error {
		if value, ok := c.UserContext().Value(CtxAuth).(*domain.JwtCustomClaim); value == nil || !ok {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.SendString(body)
	})
	jwtValidate = func(ctx context.Context, token string) (*domain.JwtCustomClaim, error) {
		return &domain.JwtCustomClaim{
			ID: "test",
		}, nil
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(fiber.HeaderAuthorization, "Bearer token")

	res, _ := app.Test(req, 100000)

	require.Equal(fiber.StatusOK, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	require.NoError(err)
	require.Equal(body, string(bytes))
}

func (suite *TestSuite) TestNewHandlerErrorInvalidToken() {
	require := require.New(suite.T())

	app := fiber.New()
	app.Use(New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	jwtValidate = func(ctx context.Context, token string) (*domain.JwtCustomClaim, error) {
		return nil, domain.ErrInvalidToken
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(fiber.HeaderAuthorization, "Invalid")

	res, _ := app.Test(req, 100000)

	require.Equal(fiber.StatusForbidden, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)

	require.NoError(err)
	require.Equal(domain.ErrInvalidToken.Error(), string(bytes))
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
