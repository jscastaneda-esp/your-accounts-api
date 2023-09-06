package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	"your-accounts-api/shared/application/mocks"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	resourceId uint
	app        *fiber.App
	mock       *mocks.ILogApp
}

func (suite *TestSuite) SetupSuite() {
	suite.resourceId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewILogApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	suite.app = fiber.New()
	suite.app.Get("/logs/:id<min(1)>", ctrl.readLogs)
}

func (suite *TestSuite) TestReadLogs200() {
	require := require.New(suite.T())
	result := &domain.Log{
		ID:          1,
		Description: "Test",
		ResourceId:  suite.resourceId,
		CreatedAt:   time.Now(),
	}
	suite.mock.On("FindLogsByProject", mock.Anything, suite.resourceId).Return([]*domain.Log{result}, nil)
	expectedBody, err := json.Marshal([]*model.ReadLogsResponse{model.NewReadLogsResponse(result)})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d", suite.resourceId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestReadLogs404() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d", 0), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadLogs500() {
	require := require.New(suite.T())
	suite.mock.On("FindLogsByProject", mock.Anything, suite.resourceId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading logs by project")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d", suite.resourceId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestNewRoute() {
	require := require.New(suite.T())
	app := fiber.New()

	NewRoute(app)

	routes := app.GetRoutes()
	require.Len(routes, 2)

	route1 := routes[0]
	require.Equal(fiber.MethodGet, route1.Method)
	require.Equal("/project/logs/:id<min(1)>", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodHead, route2.Method)
	require.Equal("/project/logs/:id<min(1)>", route2.Path)
	require.Len(route2.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
