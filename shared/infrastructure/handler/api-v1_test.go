package handler

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/shared/infrastructure/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestSuite struct {
	suite.Suite
	fastCtx *fasthttp.RequestCtx
	ctx     *fiber.Ctx
}

func (suite *TestSuite) SetupSuite() {
	suite.fastCtx = new(fasthttp.RequestCtx)

	app := fiber.New()
	suite.ctx = app.AcquireCtx(suite.fastCtx)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Test")
	})

	logsRouter = func(router fiber.Router) {
		router.Get("/project/", func(c *fiber.Ctx) error {
			return c.SendString("Project")
		})
	}

	budgetsRouter = func(router fiber.Router) {
		router.Get("/budget/", func(c *fiber.Ctx) error {
			return c.SendString("Budget")
		})
	}
}

func (suite *TestSuite) SetupTest() {
	suite.fastCtx.Request.Reset()
	suite.fastCtx.Response.Reset()
}

func (suite *TestSuite) TestNewRouteSuccessData() {
	require := require.New(suite.T())
	app := fiber.New()

	NewRoute(app)

	routes := app.GetRoutes(true)
	require.Len(routes, 4)

	route1 := routes[0]
	require.Equal(fiber.MethodGet, route1.Method)
	require.Equal("/api/v1/project/", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodGet, route2.Method)
	require.Equal("/api/v1/budget/", route2.Path)
	require.Len(route2.Handlers, 1)

	route3 := routes[2]
	require.Equal(fiber.MethodHead, route3.Method)
	require.Equal("/api/v1/project/", route3.Path)
	require.Len(route3.Handlers, 1)

	route4 := routes[3]
	require.Equal(fiber.MethodHead, route4.Method)
	require.Equal("/api/v1/budget/", route4.Path)
	require.Len(route4.Handlers, 1)

	middleware := app.GetRoutes()
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
