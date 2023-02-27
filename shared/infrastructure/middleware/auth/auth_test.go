package auth

import (
	"api-your-accounts/shared/domain"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestNewHandlerSuccess(t *testing.T) {
	require := require.New(t)

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

func TestNewHandlerSuccessConfigNext(t *testing.T) {
	require := require.New(t)

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

func TestNewHandlerSuccessValidToken(t *testing.T) {
	require := require.New(t)

	body := "Hello, World!"

	app := fiber.New()
	app.Use(New(Config{}))
	app.Get("/", func(c *fiber.Ctx) error {
		if value, ok := c.UserContext().Value(CtxAuth).(*domain.JwtCustomClaim); value == nil || !ok {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.SendString(body)
	})
	originalJwtValidate := jwtValidate
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

	jwtValidate = originalJwtValidate
}

func TestNewHandlerErrorInvalidToken(t *testing.T) {
	require := require.New(t)

	app := fiber.New()
	app.Use(New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	originalJwtValidate := jwtValidate
	jwtValidate = func(ctx context.Context, token string) (*domain.JwtCustomClaim, error) {
		return nil, domain.ErrInvalidToken
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(fiber.HeaderAuthorization, "Invalid")

	res, _ := app.Test(req, 100000)

	require.Equal(fiber.StatusForbidden, res.StatusCode)
	bytes, err := io.ReadAll(res.Body)
	gqlError := new(gqlerror.Error)
	json.Unmarshal(bytes, gqlError)

	require.NoError(err)
	require.Equal(domain.ErrInvalidToken.Error(), gqlError.Message)

	jwtValidate = originalJwtValidate
}
