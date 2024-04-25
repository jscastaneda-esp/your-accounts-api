package availables

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/budgets/infrastructure/model"
	mocks_application "your-accounts-api/mocks/budgets/application"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	mock     *mocks_application.MockIBudgetAvailableApp
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.budgetId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks_application.NewMockIBudgetAvailableApp(suite.T())
	injection.BudgetAvailableApp = suite.mock

	token := &jwt.Token{
		Claims: &domain.JwtUserClaims{
			ID: 1,
		},
	}

	suite.app = fiber.New()
	suite.app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", token)
		return c.Next()
	})
	NewRoute(suite.app)
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

	request := httptest.NewRequest(fiber.MethodPost, "/available/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/available/", nil)
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
			Field:      "CreateAvailableRequest.name",
			Constraint: "max=40",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/available/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/available/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
