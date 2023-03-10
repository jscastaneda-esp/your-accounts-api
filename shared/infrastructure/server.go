// TODO: Pendientes tests

package infrastructure

import (
	"api-your-accounts/shared/domain/jwt"
	user "api-your-accounts/user/infrastructure/handler"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const (
	defaultPort      = "8080"
	defaultJwtSecret = "aSecret"
)

var (
	newUserRoute = user.NewRoute
)

type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

type Router func(app *fiber.App)

type Server struct {
	testing bool
	routes  []interface{}
}

func (s *Server) AddRoute(routes ...interface{}) {
	s.routes = append(s.routes, routes...)
}

func (s *Server) Listen() *fiber.App {
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
				panic(fmt.Sprintf("use: invalid route %v\n", reflect.TypeOf(r)))
			}
		}
	}

	if s.testing {
		return app
	}

	// Listening server
	log.Fatal(app.Listen(":" + port))
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
	}
}
