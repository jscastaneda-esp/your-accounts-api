package infrastructure

import (
	"api-your-accounts/infrastructure/graph"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
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
func playgroundHandler() fiber.Handler {
	handler := playground.Handler("GraphQL", "/query")

	return func(c *fiber.Ctx) error {
		w := FiberResponseWriter{c}
		r := new(http.Request)

		c.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
		handler(w, r)
		return nil
	}
}

func graphqlHandler() fiber.Handler {
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

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
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := fiber.New()

	app.Get("/", playgroundHandler())
	app.Post("/query", graphqlHandler())

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
