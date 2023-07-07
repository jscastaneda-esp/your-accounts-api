package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	"your-accounts-api/budget/application/mocks"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	name      string
	year      uint16
	month     uint8
	projectId uint
	app       *fiber.App
	mock      *mocks.IBudgetApp
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.year = 2023
	suite.month = 1
	suite.projectId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewIBudgetApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	suite.app = fiber.New()
	suite.app.Get("/:id<min(1)>", ctrl.readById)
}

func (suite *TestSuite) TestReadByUser200() {
	require := require.New(suite.T())
	result := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.On("FindById", mock.Anything, suite.projectId).Return(result, nil)
	expectedBody, err := json.Marshal(model.ReadResponse{
		ID:        result.ID,
		Name:      result.Name,
		Year:      result.Year,
		Month:     result.Month,
		ProjectId: result.ProjectId,
	})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.projectId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestReadByUser404() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", 0), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadByUser404Find() {
	require := require.New(suite.T())
	suite.mock.On("FindById", mock.Anything, suite.projectId).Return(nil, gorm.ErrRecordNotFound)
	expectedErr := []byte("Budget ID not found")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.projectId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestReadByUser500() {
	require := require.New(suite.T())
	suite.mock.On("FindById", mock.Anything, suite.projectId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading projects by user")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.projectId), nil)
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
	require.Equal("/budget/:id<min(1)>", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodHead, route2.Method)
	require.Equal("/budget/:id<min(1)>", route2.Path)
	require.Len(route2.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
