// TODO: Pendientes tests

package infrastructure

import (
	"context"
	"log"
	"os"
	"strings"

	"api-your-accounts/shared/domain/jwt"
	user "api-your-accounts/user/infrastructure/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

func NewServer() {
	log.Println("Listening server")

	// Environment Variables
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	app := fiber.New()

	// Middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${locals:requestid} | ${ip} |${status}|${method}| ${latency} | ${path}: ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
		}, ","),
	}))
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), jwt.CtxJWTSecret, jwtSecret)
		c.SetUserContext(ctx)
		return c.Next()
	})

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Status available")
	})

	// # Authentication
	auth := app.Group("/auth")
	auth.Post("/user", user.CreateUserHandler)
	auth.Post("/token", user.LoginHandler)

	// # APIs Secures
	api := app.Group("/api")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
	}))

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Route protected")
	})

	// Listening server
	log.Fatal(app.Listen(":" + port))
}
