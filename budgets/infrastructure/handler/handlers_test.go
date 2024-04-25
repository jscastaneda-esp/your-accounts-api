package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/infrastructure/model"
	mocks_application "your-accounts-api/mocks/budgets/application"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/utils/convert"
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
	id       uint
	name     string
	year     uint16
	month    uint8
	budgetId uint
	cloneId  uint
	app      *fiber.App
	mock     *mocks_application.MockIBudgetApp
}

func (suite *TestSuite) SetupSuite() {
	suite.id = 999
	suite.name = "Test"
	suite.year = 2023
	suite.month = 1
	suite.budgetId = 1
	suite.cloneId = 1
}

func (suite *TestSuite) SetupTest() {
	token := &jwt.Token{
		Claims: &shared.JwtUserClaims{
			ID: 1,
		},
	}

	suite.mock = mocks_application.NewMockIBudgetApp(suite.T())
	injection.BudgetApp = suite.mock

	suite.app = fiber.New()
	suite.app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", token)
		return c.Next()
	})
	NewRoute(suite.app)
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

	request := httptest.NewRequest(fiber.MethodPost, "/budget/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/budget/", nil)
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
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "CreateRequest.name",
			Constraint: "max=40",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/budget/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/budget/", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/budget/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestRead200() {
	require := require.New(suite.T())
	ids := []uint{suite.budgetId, suite.budgetId + 1}
	names := []string{"Test 1", "Test 2"}
	year := uint16(1)
	month := uint8(1)
	zeroFloat := 0.0
	zeroUInt := uint8(0)
	result := []domain.Budget{
		{
			ID:             &ids[0],
			Name:           &names[0],
			Year:           &year,
			Month:          &month,
			TotalAvailable: &zeroFloat,
			TotalPending:   &zeroFloat,
			PendingBills:   &zeroUInt,
		},
		{
			ID:             &ids[1],
			Name:           &names[1],
			Year:           &year,
			Month:          &month,
			TotalAvailable: &zeroFloat,
			TotalPending:   &zeroFloat,
			PendingBills:   &zeroUInt,
		},
	}
	suite.mock.On("FindByUserId", mock.Anything, mock.Anything).Return(result, nil)
	expectedBody, err := json.Marshal([]model.ReadResponse{
		model.NewReadResponse(result[0]),
		model.NewReadResponse(result[1]),
	})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, "/budget/", nil)
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

	request := httptest.NewRequest(fiber.MethodGet, "/budget/", nil)
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
	zeroFloat := 0.0
	result := domain.Budget{
		ID:               &suite.id,
		Name:             &suite.name,
		Year:             &suite.year,
		Month:            &suite.month,
		FixedIncome:      &zeroFloat,
		AdditionalIncome: &zeroFloat,
	}
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(result, nil)
	expectedBody, err := json.Marshal(model.NewReadByIDResponse(result))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
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

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/budget/%d", 0), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadByID404Find() {
	require := require.New(suite.T())
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(domain.Budget{}, gorm.ErrRecordNotFound)
	expectedErr := []byte("Budget ID not found")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
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
	suite.mock.On("FindById", mock.Anything, suite.budgetId).Return(domain.Budget{}, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading projects by user")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestChanges200() {
	require := require.New(suite.T())
	requestBody := model.ChangesRequest{
		Changes: []model.ChangeRequest{
			{
				ID:      uint(1),
				Section: domain.Main,
				Action:  shared.Update,
				Detail: map[string]any{
					"name": "Test",
				},
			},
			{
				ID:      uint(1),
				Section: domain.Available,
				Action:  shared.Delete,
				Detail: map[string]any{
					"name": "Test",
				},
			},
			{
				ID:      uint(1),
				Section: domain.Bill,
				Action:  shared.Delete,
				Detail: map[string]any{
					"description": "Test",
				},
			},
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	resultChanges := []application.ChangeResult{
		{
			Change: application.Change(requestBody.Changes[0]),
		},
		{
			Change: application.Change(requestBody.Changes[1]),
		},
		{
			Change: application.Change(requestBody.Changes[2]),
		},
	}
	suite.mock.On("Changes", mock.Anything, mock.Anything, mock.Anything).Return(resultChanges)

	request := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/budget/%d", suite.budgetId), bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestChanges400() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestChanges422_1() {
	require := require.New(suite.T())
	requestBody := model.ChangesRequest{}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "ChangesRequest.changes",
			Constraint: "min=1",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/budget/%d", suite.budgetId), bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnprocessableEntity, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestChanges422_2() {
	require := require.New(suite.T())
	requestBody := model.ChangesRequest{
		Changes: []model.ChangeRequest{
			{
				ID:     uint(1),
				Action: shared.Update,
				Detail: map[string]any{
					"name": "Test",
				},
			},
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "ChangesRequest.changes[0].section",
			Constraint: "required",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/budget/%d", suite.budgetId), bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnprocessableEntity, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestChanges500() {
	require := require.New(suite.T())
	requestBody := model.ChangesRequest{
		Changes: []model.ChangeRequest{
			{
				ID:      uint(1),
				Section: domain.Main,
				Action:  shared.Update,
				Detail: map[string]any{
					"name": "Test",
				},
			},
			{
				ID:      uint(1),
				Section: domain.Available,
				Action:  shared.Delete,
				Detail: map[string]any{
					"name": "Test",
				},
			},
			{
				ID:      uint(1),
				Section: domain.Bill,
				Action:  shared.Delete,
				Detail: map[string]any{
					"description": "Test",
				},
			},
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	resultChanges := []application.ChangeResult{
		{
			Change: application.Change(requestBody.Changes[0]),
			Err:    application.ErrIncompleteData,
		},
		{
			Change: application.Change(requestBody.Changes[1]),
			Err:    convert.ErrValueIncompatibleType,
		},
		{
			Change: application.Change(requestBody.Changes[2]),
			Err:    errors.New("error"),
		},
	}
	responseBody := model.ChangesResponse{
		Changes: []model.ChangeResponse{
			model.NewChangeResponse(resultChanges[0].Change, "Incomplete data"),
			model.NewChangeResponse(resultChanges[1].Change, "Incompatible data type"),
			model.NewChangeResponse(resultChanges[2].Change, "Error processing change"),
		},
	}
	expectedBody, err := json.Marshal(responseBody)
	require.NoError(err)
	suite.mock.On("Changes", mock.Anything, mock.Anything, mock.Anything).Return(resultChanges)

	request := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("/budget/%d", suite.budgetId), bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestDelete200() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.budgetId).Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestDelete404() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.budgetId).Return(gorm.ErrRecordNotFound)
	expectedErr := []byte("Budget ID not found")

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
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

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/budget/%d", suite.budgetId), nil)
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
