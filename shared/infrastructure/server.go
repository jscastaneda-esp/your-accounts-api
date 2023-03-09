// TODO: Pendientes tests

package infrastructure

import (
	"context"
	"log"
	"os"
	"strings"

	"api-your-accounts/shared/domain/jwt"
	user "api-your-accounts/user/infrastructure/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

// HealckCheck godoc
//
//	@Summary		Show the status of server
//	@Description	get the status of server
//	@Tags			main
//	@Produce		plain
//	@Success		200	{string}	string	"Status available"
//	@Failure		500
//	@Router			/ [get]
func healthCheck(c *fiber.Ctx) error {
	return c.SendString("Server is up and running")
}

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
	{
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
	}

	// Routes
	{
		//# Root
		app.Get("/", healthCheck)
		app.Get("/swagger/*", swagger.New(swagger.Config{
			Title: "Doc API",
		}))

		//# /user
		user.NewGroup(app)

		//# /api
		{
			api := app.Group("/api")
			api.Use(jwtware.New(jwtware.Config{
				SigningKey: []byte(jwtSecret),
			}))

			//## /v1
			{
				v1 := api.Group("/v1")
				v1.Get("/", func(c *fiber.Ctx) error {
					return c.SendString("Funciona")
				})
			}
		}
	}

	// Listening server
	log.Fatal(app.Listen(":" + port))
}
