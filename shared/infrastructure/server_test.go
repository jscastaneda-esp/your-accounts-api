package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CustomRuntimeError struct{}

func (c CustomRuntimeError) Error() string {
	return "runtime error"
}

func (c CustomRuntimeError) RuntimeError() {}

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestAddRouteSuccess() {
	require := require.New(suite.T())
	server := NewServer(true)

	server.AddRoute(Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(_ echo.Context) error {
			return nil
		},
	})
	server.AddRoute(Router(func(_ *echo.Echo) {
		log.Info("Test")
	}))
	server.Listen()

	require.NotEmpty(server.routes)
	require.Len(server.routes, 2)
	require.IsType(Route{}, server.routes[0])
	require.IsType(Router(nil), server.routes[1])
}

func (suite *TestSuite) TestHandlerSuccess() {
	require := require.New(suite.T())
	server := NewServer(true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e := server.Listen()

	ctx := e.NewContext(req, rec)
	err := healthCheck(ctx)

	require.NoError(err)
	require.Equal(http.StatusOK, rec.Code)
	require.Equal("Server is up and running", rec.Body.String())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
