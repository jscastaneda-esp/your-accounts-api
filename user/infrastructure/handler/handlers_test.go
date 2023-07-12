package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"your-accounts-api/shared/domain/validation"
	"your-accounts-api/user/application"
	"your-accounts-api/user/application/mocks"
	"your-accounts-api/user/infrastructure/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	uid   string
	email string
	token string
	app   *fiber.App
	mock  *mocks.IUserApp
}

func (suite *TestSuite) SetupSuite() {
	suite.uid = "01234567890123456789012345678901"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewIUserApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	suite.app = fiber.New()
	suite.app.Post("/user", ctrl.create)
	suite.app.Post("/login", ctrl.login)
}

func (suite *TestSuite) TestCreate201() {
	require := require.New(suite.T())
	requestBody := &model.CreateRequest{
		UID:   suite.uid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.UID, requestBody.Email).Return(uint(1), nil)
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
		UID:   suite.uid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.UID, requestBody.Email).Return(uint(0), application.ErrUserAlreadyExists)
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
		UID:   "invalid",
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "uid",
		Constraint: "len=32",
	})
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
		UID:   suite.uid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, requestBody.UID, requestBody.Email).Return(uint(0), gorm.ErrInvalidField)
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
			UID:   suite.uid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Login", mock.Anything, suite.uid, suite.email).Return(suite.token, nil)
	expectedBody, err := json.Marshal(model.NewLoginResponse(suite.token))
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
			UID:   suite.uid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Login", mock.Anything, suite.uid, suite.email).Return("", gorm.ErrRecordNotFound)
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
			UID:   "invalid",
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "uid",
		Constraint: "len=32",
	})
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
			UID:   suite.uid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Login", mock.Anything, suite.uid, suite.email).Return("", gorm.ErrInvalidField)
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

func (suite *TestSuite) TestNewRoute() {
	require := require.New(suite.T())
	app := fiber.New()

	NewRoute(app)

	routes := app.GetRoutes()
	require.Len(routes, 2)

	route1 := routes[0]
	require.Equal(fiber.MethodPost, route1.Method)
	require.Equal("/user", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodPost, route2.Method)
	require.Equal("/login", route2.Path)
	require.Len(route2.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
