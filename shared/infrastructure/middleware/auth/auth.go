package auth

import (
	"api-your-accounts/shared/domain"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	CtxAuth = authKey("auth")
	bearer  = "Bearer "
)

var (
	jwtValidate = domain.JwtValidate
)

type authKey string

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		ctx := context.WithValue(c.UserContext(), domain.CtxJWTSecret, cfg.JWTSecret)
		auth := c.Get("Authorization")
		if auth == "" {
			c.SetUserContext(ctx)
			return c.Next()
		}

		token := auth[len(bearer):]
		claims, err := jwtValidate(ctx, token)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(&gqlerror.Error{
				Message: err.Error(),
			})
		}

		ctx = context.WithValue(ctx, CtxAuth, claims)
		c.SetUserContext(ctx)

		// Continue stack
		return c.Next()
	}
}
