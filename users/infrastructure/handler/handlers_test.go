package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	mocks_application "your-accounts-api/mocks/users/application"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"
	"your-accounts-api/users/application"
	"your-accounts-api/users/infrastructure/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	email string
	token string
	app   *fiber.App
	mock  *mocks_application.MockIUserApp
}

func (suite *TestSuite) SetupSuite() {
	suite.email = "example@exaple.com"
	suite.token = "<token>"
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks_application.NewMockIUserApp(suite.T())

	injection.UserApp = suite.mock
	suite.app = fiber.New()
	NewRoute(suite.app)
}

func (suite *TestSuite) TestCreate201() {
	require := require.New(suite.T())
	requestBody := &model.CreateRequest{
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.Email).Return(uint(1), nil)
	expectedBody, err := json.Marshal(model.NewCreateResponse(uint(1)))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/user", bytes.NewReader(body))
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

	request := httptest.NewRequest(fiber.MethodPost, "/user", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestCreate409() {
	require := require.New(suite.T())
	requestBody := &model.CreateRequest{
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.Email).Return(uint(0), application.ErrUserAlreadyExists)
	expectedErr := []byte(application.ErrUserAlreadyExists.Error())

	request := httptest.NewRequest(fiber.MethodPost, "/user", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusConflict, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestCreate422() {
	require := require.New(suite.T())
	requestBody := &model.CreateRequest{
		Email: "invalid",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "CreateRequest.email",
			Constraint: "email",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/user", bytes.NewReader(body))
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
	requestBody := &model.CreateRequest{
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.Email).Return(uint(0), gorm.ErrInvalidField)
	expectedErr := []byte("Error sign up user")

	request := httptest.NewRequest(fiber.MethodPost, "/user", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestLogin200() {
	require := require.New(suite.T())
	requestBody := &model.LoginRequest{
		CreateRequest: model.CreateRequest{
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	expiresAt := time.Now()
	suite.mock.On("Login", mock.Anything, suite.email).Return(suite.token, expiresAt, nil)
	expectedBody, err := json.Marshal(model.NewLoginResponse(suite.token, expiresAt))
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestLogin400() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodPost, "/login", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestLogin401() {
	require := require.New(suite.T())
	requestBody := &model.LoginRequest{
		CreateRequest: model.CreateRequest{
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Login", mock.Anything, suite.email).Return("", time.Time{}, gorm.ErrRecordNotFound)
	expectedErr := []byte("Invalid credentials")

	request := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnauthorized, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestLogin422() {
	require := require.New(suite.T())
	requestBody := &model.LoginRequest{
		CreateRequest: model.CreateRequest{
			Email: "invalid",
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "LoginRequest.CreateRequest.email",
			Constraint: "email",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusUnprocessableEntity, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestLogin500() {
	require := require.New(suite.T())
	requestBody := &model.LoginRequest{
		CreateRequest: model.CreateRequest{
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Login", mock.Anything, suite.email).Return("", time.Time{}, gorm.ErrInvalidField)
	expectedErr := []byte("Error authenticate user")

	request := httptest.NewRequest(fiber.MethodPost, "/login", bytes.NewReader(body))
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
