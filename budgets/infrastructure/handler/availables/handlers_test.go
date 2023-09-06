package availables

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/budgets/application/mocks"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2"
	go_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	name     string
	budgetId uint
	app      *fiber.App
	mock     *mocks.IBudgetAvailableApp
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.budgetId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewIBudgetAvailableApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	token := &go_jwt.Token{
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
}

func (suite *TestSuite) TestCreate201() {
	require := require.New(suite.T())
	requestBody := model.CreateAvailableRequest{
		Name:     suite.name,
		BudgetId: suite.budgetId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, suite.name, suite.budgetId).Return(uint(1), nil)
	expectedBody, err := json.Marshal(model.NewCreateResponse(uint(1)))
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
	requestBody := model.CreateAvailableRequest{
		Name:     "Cupidatat ullamco voluptate non aute consequat fugiat.",
		BudgetId: suite.budgetId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "name",
			Constraint: "max=40",
		},
	}
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

func (suite *TestSuite) TestCreate500() {
	require := require.New(suite.T())
	requestBody := model.CreateAvailableRequest{
		Name:     suite.name,
		BudgetId: suite.budgetId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(uint(0), gorm.ErrInvalidField)
	expectedErr := []byte("Error creating available")

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

func (suite *TestSuite) TestNewRoute() {
	require := require.New(suite.T())
	app := fiber.New()

	NewRoute(app)

	routes := app.GetRoutes()
	require.Len(routes, 1)

	route1 := routes[0]
	require.Equal(fiber.MethodPost, route1.Method)
	require.Equal("/available/", route1.Path)
	require.Len(route1.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
