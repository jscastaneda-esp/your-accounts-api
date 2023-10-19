package infrastructure

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
	"your-accounts-api/shared/infrastructure/config"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"

	_ "your-accounts-api/docs"
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
	log.Info("Listening server")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			} else {
				switch err.(type) {
				case runtime.Error:
					err = errors.New("internal server error")
				}
			}

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			// Return status code with error message
			return c.Status(code).SendString(err.Error())
		},
	})

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
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
		app.Use(requestid.New())
	}

	// Routes
	{
		//# Root
		app.Get("/", healthCheck)
		app.Get("/swagger/*", swagger.HandlerDefault)

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
				log.Fatalf("use: invalid route %v\n", reflect.TypeOf(r))
			}
		}
	}

	if s.testing {
		log.Info("Listen server on port", config.PORT)
		return app
	}

	// Listening server
	log.Fatal(app.Listen(":" + config.PORT))
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
		routes:  make([]any, 0),
	}
}
