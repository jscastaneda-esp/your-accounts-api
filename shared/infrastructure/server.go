package infrastructure

import (
	"net/http"
	"reflect"
	"strings"
	"time"
	"your-accounts-api/shared/infrastructure/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "your-accounts-api/docs"
)

type Route struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
}

type Router func(e *echo.Echo)

type Server struct {
	testing bool
	routes  []any
}

func (s *Server) AddRoute(route any) {
	s.routes = append(s.routes, route)
}

func (s *Server) Listen() *echo.Echo {
	log.Info("Listening server")

	e := echo.New()
	e.Validator = NewCustomValidator()
	e.Binder = NewCustomBinder()

	// Middlewares
	{
		e.Use(middleware.Logger())
		e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Skipper: func(c echo.Context) bool {
				return c.RealIP() == "127.0.0.1" || strings.HasPrefix(c.RealIP(), "172")
			},
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: 10, ExpiresIn: 1 * time.Minute},
			),
		}))
		e.Use(middleware.Recover())
		e.Use(middleware.RequestID())
	}

	// Routes
	{
		//# Root
		e.GET("/", healthCheck)
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		// # Additional
		for _, route := range s.routes {
			switch r := route.(type) {
			case Route:
				e.Add(r.Method, r.Path, r.Handler)
			case Router:
				r(e)
			default:
				log.Fatalf("use: invalid route %v\n", reflect.TypeOf(r))
			}
		}
	}

	if s.testing {
		log.Info("Listen server on port", config.PORT)
		return e
	}

	// Listening server
	log.Fatal(e.Start(":" + config.PORT))
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
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Server is up and running")
}

func NewServer(testing bool) *Server {
	return &Server{
		testing: testing,
		routes:  make([]any, 0),
	}
}
