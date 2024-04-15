package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"your-accounts-api/shared/infrastructure/config"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	})

	logsRouter = func(router *echo.Echo) {
		router.GET("/project/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Project")
		})
	}

	budgetsRouter = func(router *echo.Echo) {
		router.GET("/budget/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Budget")
		})
	}
}

func (suite *TestSuite) TestNewRouteSuccessData() {
	require := require.New(suite.T())
	e := echo.New()

	NewRoute(e)

	routes := e.GetRoutes(true)
	require.Len(routes, 4)

	route1 := routes[0]
	require.Equal(http.MethodGet, route1.Method)
	require.Equal("/api/v1/project/", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(http.MethodGet, route2.Method)
	require.Equal("/api/v1/budget/", route2.Path)
	require.Len(route2.Handlers, 1)

	route3 := routes[2]
	require.Equal(http.MethodHead, route3.Method)
	require.Equal("/api/v1/project/", route3.Path)
	require.Len(route3.Handlers, 1)

	route4 := routes[3]
	require.Equal(http.MethodHead, route4.Method)
	require.Equal("/api/v1/budget/", route4.Path)
	require.Len(route4.Handlers, 1)

	middleware := e.GetRoutes()
	useFilter := make([]fiber.Route, 0)
	for _, m := range middleware {
		if m.Path == "/api/v1" {
			useFilter = append(useFilter, m)
		}
	}
	require.Len(useFilter, 9)

	handler := useFilter[0].Handlers
	require.Len(useFilter[0].Handlers, 1)
	for i := 1; i < len(useFilter); i++ {
		require.Len(useFilter[i].Handlers, 1)
		require.Equal(handler, useFilter[i].Handlers)
	}
}

func (suite *TestSuite) TestNewRouteSuccessRequest() {
	require := require.New(suite.T())
	request := httptest.NewRequest(fiber.MethodGet, "/api/v1/project", nil)
	request.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.UWfppjZwV_hQ4PU5O9ds7s9jxK4l6u6PDmAHkoVuFpg"))
	app := fiber.New()
	config.JWT_SECRET = "aSecret"

	NewRoute(app)
	response, err := app.Test(request, 10)

	require.NoError(err)
	require.NotNil(response)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Project"), resp)
}

func (suite *TestSuite) TestNewRouteErrorUnauthorized() {
	require := require.New(suite.T())
	request := httptest.NewRequest(fiber.MethodGet, "/api/v1/budget", nil)
	request.Header.Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.UWfppjZwV_hQ4PU5O9ds7s9jxK4l6u6PDmAHkoVuFpg"))
	app := fiber.New()
	config.JWT_SECRET = "other"

	NewRoute(app)
	response, err := app.Test(request, 10)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnauthorized, response.StatusCode)

	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal([]byte("Invalid or expired JWT"), resp)
}

// func (suite *TestSuite) TestNewRouteErrorBadRequest() {
// 	require := require.New(suite.T())
// 	request := httptest.NewRequest(fiber.MethodGet, "/api/v1", nil)
// 	app := fiber.New()

// 	NewRoute(app)
// 	response, err := app.Test(request, 10000)

// 	require.NoError(err)
// 	require.NotNil(response)
// 	require.Equal(fiber.StatusBadRequest, response.StatusCode)

// 	resp, err := io.ReadAll(response.Body)
// 	require.NoError(err)
// 	require.Equal([]byte("Missing or malformed JWT"), resp)
// }

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
