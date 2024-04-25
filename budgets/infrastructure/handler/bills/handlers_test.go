package bills

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/infrastructure/model"
	mocks_application "your-accounts-api/mocks/budgets/application"
	shared "your-accounts-api/shared/domain"
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
	description string
	category    domain.BudgetBillCategory
	budgetId    uint
	billId      uint
	app         *fiber.App
	mock        *mocks_application.MockIBudgetBillApp
}

func (suite *TestSuite) SetupSuite() {
	suite.description = "Test"
	suite.category = domain.Entertainment
	suite.budgetId = 1
	suite.billId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks_application.NewMockIBudgetBillApp(suite.T())
	injection.BudgetBillApp = suite.mock

	token := &jwt.Token{
		Claims: &shared.JwtUserClaims{
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
	requestBody := model.CreateBillRequest{
		Description: suite.description,
		BudgetId:    suite.budgetId,
		Category:    suite.category,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, suite.description, suite.category, suite.budgetId).Return(uint(1), nil)
	expectedBody, err := json.Marshal(model.NewCreateResponse(uint(1)))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/bill/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/bill/", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestCreate422() {
	require := require.New(suite.T())
	requestBody := model.CreateBillRequest{
		Description: suite.description,
		BudgetId:    suite.budgetId,
		Category:    "Test",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "CreateBillRequest.category",
			Constraint: "oneof='house' 'entertainment' 'personal' 'vehicle_transportation' 'education' 'services' 'financial' 'saving' 'others'",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/bill/", bytes.NewReader(body))
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
	requestBody := model.CreateBillRequest{
		Description: suite.description,
		BudgetId:    suite.budgetId,
		Category:    suite.category,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint(0), gorm.ErrInvalidField)
	expectedErr := []byte("Error creating bill")

	request := httptest.NewRequest(fiber.MethodPost, "/bill/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestCreateTransaction200() {
	require := require.New(suite.T())
	requestBody := model.CreateBillTransactionRequest{
		Description: suite.description,
		Amount:      float64(1000),
		BillId:      suite.billId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("CreateTransaction", mock.Anything, suite.description, requestBody.Amount, suite.billId).Return(nil)

	request := httptest.NewRequest(fiber.MethodPut, "/bill/transaction", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestCreateTransaction400() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodPut, "/bill/transaction", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestCreateTransaction422() {
	require := require.New(suite.T())
	requestBody := model.CreateBillTransactionRequest{
		Description: suite.description,
		BillId:      suite.billId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "CreateBillTransactionRequest.amount",
			Constraint: "required",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPut, "/bill/transaction", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnprocessableEntity, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestCreateTransaction500() {
	require := require.New(suite.T())
	requestBody := model.CreateBillTransactionRequest{
		Description: suite.description,
		Amount:      float64(1000),
		BillId:      suite.billId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gorm.ErrInvalidField)
	expectedErr := []byte("Error creating bill transaction")

	request := httptest.NewRequest(fiber.MethodPut, "/bill/transaction", bytes.NewReader(body))
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
