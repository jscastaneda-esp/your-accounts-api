package handler

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

const (
	defaultJwtSecret = "aSecret"
)

func NewRoute(app *fiber.App) {
	jwtSecret := os.Getenv("JWT_SECRET")

	api := app.Group("/api/v1")
	// Middleware
	{
		api.Use(jwtware.New(jwtware.Config{
			SigningKey: []byte(jwtSecret),
		}))
	}

	// Routes
	{
		api.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Funciona")
		})
	}
}
