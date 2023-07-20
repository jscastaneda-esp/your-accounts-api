package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	"your-accounts-api/budget/application/mocks"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/model"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2"
	goJwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	name     string
	year     uint16
	month    uint8
	budgetId uint
	cloneId  uint
	app      *fiber.App
	mock     *mocks.IBudgetApp
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.year = 2023
	suite.month = 1
	suite.budgetId = 1
	suite.cloneId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewIBudgetApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	token := &goJwt.Token{
		Claims: &jwt.JwtUserClaims{
			ID: 1,
		},
	}

	suite.app = fiber.New()
	suite.app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", token)
		return c.Next()
	})
	suite.app.Post("/", ctrl.create)
	suite.app.Get("/", ctrl.read)
	suite.app.Get("/:id<min(1)>", ctrl.readById)
	suite.app.Delete("/:id<min(1)>", ctrl.delete)
}

func (suite *TestSuite) TestCreate201() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name: "Test",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(suite.budgetId, nil)
	expectedBody, err := json.Marshal(model.NewCreateResponse(suite.budgetId))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusCreated, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestCreate400() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodPost, "/", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestCreate422() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name: "Cupidatat ullamco voluptate non aute consequat fugiat.",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "name",
		Constraint: "max=40",
	})
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnprocessableEntity, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestCreate404() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		CloneId: &suite.cloneId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Clone", mock.Anything, mock.Anything, *requestBody.CloneId).Return(uint(0), gorm.ErrRecordNotFound)
	expectedErr := []byte("Error creating budget. Clone ID not found")

	request := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestCreate500() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name: "Test",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(uint(0), gorm.ErrInvalidField)
	expectedErr := []byte("Error creating budget")

	request := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request, 6000000)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestRead200() {
	require := require.New(suite.T())
	result := []*domain.Budget{
		{
			ID:    suite.budgetId,
			Name:  "Test",
			Year:  1,
			Month: 1,
		},
		{
			ID:    suite.budgetId + 1,
			Name:  "Test 2",
			Year:  2,
			Month: 1,
		},
	}
	suite.mock.On("FindByUserId", mock.Anything, mock.Anything).Return(result, nil)
	expectedBody, err := json.Marshal([]model.ReadResponse{
		model.NewReadResponse(result[0]),
		model.NewReadResponse(result[1]),
	})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, "/", nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestRead500() {
	require := require.New(suite.T())
	suite.mock.On("FindByUserId", mock.Anything, mock.Anything).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading budgets by user")

	request := httptest.NewRequest(fiber.MethodGet, "/", nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestReadByID200() {
	require := require.New(suite.T())
	result := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.budgetId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(result, nil)
	expectedBody, err := json.Marshal(model.NewReadByIDResponse(result))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestReadByID404() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", 0), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadByID404Find() {
	require := require.New(suite.T())
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(nil, gorm.ErrRecordNotFound)
	expectedErr := []byte("Budget ID not found")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestReadByID500() {
	require := require.New(suite.T())
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading projects by user")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestDelete200() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.budgetId).Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestDelete404() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.budgetId).Return(gorm.ErrRecordNotFound)
	expectedErr := []byte("Budget ID not found")

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestDelete500() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.budgetId).Return(gorm.ErrInvalidField)
	expectedErr := []byte("Error deleting budget")

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.budgetId), nil)
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
	require.Len(routes, 6)

	route1 := routes[0]
	require.Equal(fiber.MethodGet, route1.Method)
	require.Equal("/budget/", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodGet, route2.Method)
	require.Equal("/budget/:id<min(1)>", route2.Path)
	require.Len(route2.Handlers, 1)

	route3 := routes[2]
	require.Equal(fiber.MethodHead, route3.Method)
	require.Equal("/budget/", route3.Path)
	require.Len(route3.Handlers, 1)

	route4 := routes[3]
	require.Equal(fiber.MethodHead, route4.Method)
	require.Equal("/budget/:id<min(1)>", route4.Path)
	require.Len(route4.Handlers, 1)

	route5 := routes[4]
	require.Equal(fiber.MethodPost, route5.Method)
	require.Equal("/budget/", route5.Path)
	require.Len(route5.Handlers, 1)

	route6 := routes[5]
	require.Equal(fiber.MethodDelete, route6.Method)
	require.Equal("/budget/:id<min(1)>", route6.Path)
	require.Len(route6.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
