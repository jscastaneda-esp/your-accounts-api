package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"your-accounts-api/shared/domain/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

type Router func(app fiber.Router)

type Server struct {
	testing bool
	routes  []any
}

func (s *Server) AddRoute(route any) {
	s.routes = append(s.routes, route)
}

func (s *Server) Listen() *fiber.App {
	log.Println("Listening server")

	// Environment Variables
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = defaultJwtSecret
		os.Setenv("JWT_SECRET", jwtSecret)
	}

	app := fiber.New()

	// Middleware
	{
		const patternLog string = `${%s} |${%s} ${%srequestid} ${%s}| ${%s} |${%s}|${%s}| ${%s} | ${%s}
Headers: ${%s}
Params: ${%s}
Body: ${%s}
Response: ${%s}

`
		app.Use(logger.New(logger.Config{
			Format: fmt.Sprintf(patternLog,
				logger.TagTime, logger.TagMagenta, logger.TagLocals, logger.TagReset, logger.TagIP, logger.TagStatus, logger.TagMethod, logger.TagLatency, logger.TagPath, logger.TagReqHeaders, logger.TagQueryStringParams, logger.TagBody, logger.TagResBody),
			TimeFormat: "2006/01/02 15:04:05",
		}))
		app.Use(limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1" || strings.HasPrefix(c.IP(), "172")
			},
			Max:        10,
			Expiration: 1 * time.Minute,
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

		// # Additional
		for _, route := range s.routes {
			switch r := route.(type) {
			case Route:
				if fiber.MethodGet == r.Method {
					app.Get(r.Path, r.Handler)
				} else {
					app.Add(r.Method, r.Path, r.Handler)
				}
			case Router:
				r(app)
			default:
				log.Panicf("use: invalid route %v\n", reflect.TypeOf(r))
			}
		}
	}

	if s.testing {
		log.Printf("Listen server on port %s\n", port)
		return app
	}

	// Listening server
	log.Panic(app.Listen(":" + port))
	return nil
}

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

func NewServer(testing bool) *Server {
	return &Server{
		testing: testing,
		routes:  []any{},
	}
}
