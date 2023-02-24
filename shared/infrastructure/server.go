package infrastructure

import (
	"api-your-accounts/shared/infrastructure/graph"
	"api-your-accounts/shared/infrastructure/graph/resolver"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const defaultPort = "8080"

type FiberResponseWriter struct {
	ctx *fiber.Ctx
}

func (FiberResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w FiberResponseWriter) Write(bytes []byte) (int, error) {
	return w.ctx.Write(bytes)
}

func (w FiberResponseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}

// Defining the Playground handler
func getPlaygroundHandler() fiber.Handler {
	handler := playground.Handler("GraphQL", "/query")

	return func(c *fiber.Ctx) error {
		w := FiberResponseWriter{c}
		r := new(http.Request)

		c.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
		handler(w, r)
		return nil
	}
}

// Defining the GraphQL handler
func postGraphqlHandler() fiber.Handler {
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}))

	return func(c *fiber.Ctx) error {
		headers := make(http.Header)
		for key, value := range c.GetReqHeaders() {
			headers[key] = []string{value}
		}

		w := FiberResponseWriter{c}
		r := &http.Request{
			Method: c.Method(),
			Header: headers,
			Body:   io.NopCloser(bytes.NewReader(c.Body())),
		}
		r = r.WithContext(c.UserContext())

		c.Set("Content-Type", fiber.MIMEApplicationJSON)
		server.ServeHTTP(w, r)
		return nil
	}
}

func NewServer() {
	// Environment Variables
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New()

	// Middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${locals:requestid} | ${ip} |${status}|${method}| ${latency} | ${path} | ${error}]\n",
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

	// Routes
	app.Get("/", getPlaygroundHandler())
	app.Post("/query", postGraphqlHandler())

	// Listening server
	log.Fatal(app.Listen(":" + port))
}
