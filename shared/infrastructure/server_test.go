package infrastructure

import (
	"io"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestAddRouteSuccess() {
	require := require.New(suite.T())
	server := NewServer(true)

	server.AddRoute(Route{
		Method: fiber.MethodGet,
		Path:   "/",
		Handler: func(_ *fiber.Ctx) error {
			return nil
		},
	})
	server.AddRoute(Router(func(_ *fiber.App) {
		log.Println("Test")
	}))

	require.NotEmpty(server.routes)
	require.Len(server.routes, 2)
	require.IsType(Route{}, server.routes[0])
	require.IsType(Router(nil), server.routes[1])
}

func (suite *TestSuite) TestListenSuccessCustomPort() {
	require := require.New(suite.T())
	server := NewServer(true)
	request := httptest.NewRequest(fiber.MethodGet, "/", nil)
	os.Setenv("PORT", "999")

	app := server.Listen()
	response, err := app.Test(request, 1)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Server is up and running"), resp)
}

func (suite *TestSuite) TestListenSuccessDefaultRoutes() {
	require := require.New(suite.T())
	server := NewServer(true)
	request := httptest.NewRequest(fiber.MethodGet, "/", nil)

	app := server.Listen()
	response, err := app.Test(request, 1)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Server is up and running"), resp)
}

func (suite *TestSuite) TestListenSuccessCustomRouteGet() {
	require := require.New(suite.T())
	server := NewServer(true)
	server.AddRoute(Route{
		Method: fiber.MethodGet,
		Path:   "/route-get",
		Handler: func(c *fiber.Ctx) error {
			return c.SendString("Route Get")
		},
	})
	request := httptest.NewRequest(fiber.MethodGet, "/route-get", nil)

	app := server.Listen()
	response, err := app.Test(request, 1)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Route Get"), resp)
}

func (suite *TestSuite) TestListenSuccessCustomRoutePost() {
	require := require.New(suite.T())
	server := NewServer(true)
	server.AddRoute(Route{
		Method: fiber.MethodPost,
		Path:   "/route-post",
		Handler: func(c *fiber.Ctx) error {
			return c.SendString("Route Post")
		},
	})
	request := httptest.NewRequest(fiber.MethodPost, "/route-post", nil)

	app := server.Listen()
	response, err := app.Test(request, 1)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Route Post"), resp)
}

func (suite *TestSuite) TestListenSuccessCustomRouteRouter() {
	require := require.New(suite.T())
	server := NewServer(true)
	server.AddRoute(Router(func(app *fiber.App) {
		app.Get("/router", func(c *fiber.Ctx) error {
			return c.SendString("Router")
		})
	}))
	request := httptest.NewRequest(fiber.MethodGet, "/router", nil)

	app := server.Listen()
	response, err := app.Test(request, 1)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Router"), resp)
}

func (suite *TestSuite) TestListenErrorCustomRoutePanic() {
	require := require.New(suite.T())
	server := NewServer(true)
	server.AddRoute("panic")

	require.Panics(func() {
		server.Listen()
	})
}

func (suite *TestSuite) TestListenErrorPanic() {
	require := require.New(suite.T())
	server := NewServer(false)
	os.Setenv("PORT", "9999999")

	require.Panics(func() {
		server.Listen()
	})
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
