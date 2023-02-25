package auth

import "github.com/gofiber/fiber/v2"

type Config struct {
	Next      func(c *fiber.Ctx) bool
	JWTSecret string
}

var ConfigDefault = Config{
	Next:      nil,
	JWTSecret: "aSecret",
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = ConfigDefault.JWTSecret
	}
	return cfg
}
