package handler

import (
	"api-your-accounts/shared/domain/validation"
	"api-your-accounts/user/application"
	"api-your-accounts/user/application/mocks"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/model"
	"encoding/json"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	uuid       string
	email      string
	token      string
	fastCtx    *fasthttp.RequestCtx
	ctx        *fiber.Ctx
	mock       *mocks.IUserApp
	controller *userController
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "01234567890123456789012345678901"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.fastCtx = &fasthttp.RequestCtx{}

	app := fiber.New()
	suite.ctx = app.AcquireCtx(suite.fastCtx)
}

func (suite *TestSuite) SetupTest() {
	suite.fastCtx.Request.Reset()
	suite.fastCtx.Response.Reset()
	suite.mock = mocks.NewIUserApp(suite.T())
	suite.controller = &userController{
		app: suite.mock,
	}
}

func (suite *TestSuite) TestCreateUserSuccess() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	result := &domain.User{
		ID:        1,
		UUID:      requestBody.UUID,
		Email:     requestBody.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.mock.On("Exists", suite.ctx.UserContext(), suite.uuid, suite.email).Return(false, nil)
	suite.mock.On("SignUp", suite.ctx.UserContext(), mock.Anything).Return(result, nil)

	expectedBody, err := json.Marshal(model.CreateResponse{
		ID:        result.ID,
		UUID:      result.UUID,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
	require.NoError(err)

	err = suite.controller.createUser(suite.ctx)

	require.NoError(err)
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestCreateUserSuccessRecordNotFound() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	result := &domain.User{
		ID:        1,
		UUID:      requestBody.UUID,
		Email:     requestBody.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.mock.On("Exists", suite.ctx.UserContext(), suite.uuid, suite.email).Return(false, gorm.ErrRecordNotFound)
	suite.mock.On("SignUp", suite.ctx.UserContext(), mock.Anything).Return(result, nil)

	expectedBody, err := json.Marshal(model.CreateResponse{
		ID:        result.ID,
		UUID:      result.UUID,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
	require.NoError(err)

	err = suite.controller.createUser(suite.ctx)

	require.NoError(err)
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestCreateUserErrorBodyParser() {
	require := require.New(suite.T())
	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnprocessableEntity,
		Message: "unexpected end of JSON input",
	}

	err := suite.controller.createUser(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestCreateUserErrorInvalidStruct() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  "invalid",
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "uuid",
		Constraint: "len=32",
	})
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	err = suite.controller.createUser(suite.ctx)

	require.NoError(err)
	require.Equal(fiber.StatusUnprocessableEntity, suite.fastCtx.Response.StatusCode())
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestCreateUserErrorConflict() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Exists", suite.ctx.UserContext(), suite.uuid, suite.email).Return(true, nil)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnprocessableEntity,
		Message: "User already exists",
	}

	err = suite.controller.createUser(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestCreateUserErrorExists() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Exists", suite.ctx.UserContext(), suite.uuid, suite.email).Return(false, gorm.ErrInvalidField)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: "Error sign up user",
	}

	err = suite.controller.createUser(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestCreateUserErrorSignUp() {
	require := require.New(suite.T())

	requestBody := &model.CreateRequest{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Exists", suite.ctx.UserContext(), suite.uuid, suite.email).Return(false, nil)
	suite.mock.On("SignUp", suite.ctx.UserContext(), mock.Anything).Return(nil, gorm.ErrInvalidField)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: "Error sign up user",
	}

	err = suite.controller.createUser(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestAuthSuccess() {
	require := require.New(suite.T())

	requestBody := &model.AuthRequest{
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Auth", suite.ctx.UserContext(), suite.uuid, suite.email).Return(suite.token, nil)

	expectedBody, err := json.Marshal(fiber.Map{
		"token": suite.token,
	})
	require.NoError(err)

	err = suite.controller.auth(suite.ctx)

	require.NoError(err)
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestAuthErrorBodyParser() {
	require := require.New(suite.T())
	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnprocessableEntity,
		Message: "unexpected end of JSON input",
	}

	err := suite.controller.auth(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestAuthErrorInvalidStruct() {
	require := require.New(suite.T())

	requestBody := &model.AuthRequest{
		CreateRequest: model.CreateRequest{
			UUID:  "invalid",
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "uuid",
		Constraint: "len=32",
	})
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	err = suite.controller.auth(suite.ctx)

	require.NoError(err)
	require.Equal(fiber.StatusUnprocessableEntity, suite.fastCtx.Response.StatusCode())
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestAuthErrorUnauthorized() {
	require := require.New(suite.T())

	requestBody := &model.AuthRequest{
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Auth", suite.ctx.UserContext(), suite.uuid, suite.email).Return("", gorm.ErrRecordNotFound)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnauthorized,
		Message: "Invalid credentials",
	}

	err = suite.controller.auth(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestAuthError() {
	require := require.New(suite.T())

	requestBody := &model.AuthRequest{
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("Auth", suite.ctx.UserContext(), suite.uuid, suite.email).Return("", gorm.ErrInvalidField)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: "Error authenticate user",
	}

	err = suite.controller.auth(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestRefreshTokenSuccess() {
	require := require.New(suite.T())

	requestBody := &model.RefreshTokenRequest{
		Token: suite.token,
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("RefreshToken", suite.ctx.UserContext(), suite.token, suite.uuid, suite.email).Return(suite.token+"New", nil)

	expectedBody, err := json.Marshal(fiber.Map{
		"token": suite.token + "New",
	})
	require.NoError(err)

	err = suite.controller.refreshToken(suite.ctx)

	require.NoError(err)
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestRefreshTokenErrorBodyParser() {
	require := require.New(suite.T())
	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnprocessableEntity,
		Message: "unexpected end of JSON input",
	}

	err := suite.controller.refreshToken(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestRefreshTokenErrorInvalidStruct() {
	require := require.New(suite.T())

	requestBody := &model.RefreshTokenRequest{
		Token: "",
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	validationErrors := []*validation.ErrorResponse{}
	validationErrors = append(validationErrors, &validation.ErrorResponse{
		Field:      "token",
		Constraint: "required",
	})
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	err = suite.controller.refreshToken(suite.ctx)

	require.NoError(err)
	require.Equal(fiber.StatusUnprocessableEntity, suite.fastCtx.Response.StatusCode())
	require.Equal([]byte(fiber.MIMEApplicationJSON), suite.fastCtx.Response.Header.ContentType())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func (suite *TestSuite) TestRefreshTokenErrorUnauthorized() {
	require := require.New(suite.T())

	requestBody := &model.RefreshTokenRequest{
		Token: suite.token,
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("RefreshToken", suite.ctx.UserContext(), suite.token, suite.uuid, suite.email).Return("", gorm.ErrRecordNotFound)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusUnauthorized,
		Message: "Invalid data",
	}

	err = suite.controller.refreshToken(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestRefreshTokenErrorBadRequest() {
	require := require.New(suite.T())

	requestBody := &model.RefreshTokenRequest{
		Token: suite.token,
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("RefreshToken", suite.ctx.UserContext(), suite.token, suite.uuid, suite.email).Return("", application.ErrTokenRefreshed)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: application.ErrTokenRefreshed.Error(),
	}

	err = suite.controller.refreshToken(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestRefreshTokenError() {
	require := require.New(suite.T())

	requestBody := &model.RefreshTokenRequest{
		Token: suite.token,
		CreateRequest: model.CreateRequest{
			UUID:  suite.uuid,
			Email: suite.email,
		},
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)

	suite.fastCtx.Request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)

	suite.mock.On("RefreshToken", suite.ctx.UserContext(), suite.token, suite.uuid, suite.email).Return("", gorm.ErrInvalidField)

	expectedErr := &fiber.Error{
		Code:    fiber.StatusInternalServerError,
		Message: "Error refresh token user",
	}

	err = suite.controller.refreshToken(suite.ctx)

	require.EqualError(expectedErr, err.Error())
}

func (suite *TestSuite) TestNewRoute() {
	require := require.New(suite.T())
	app := fiber.New()

	NewRoute(app)

	routes := app.GetRoutes()
	require.Len(routes, 3)

	route1 := routes[0]
	require.Equal(fiber.MethodPost, route1.Method)
	require.Equal("/user/", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodPost, route2.Method)
	require.Equal("/user/auth", route2.Path)
	require.Len(route2.Handlers, 1)

	route3 := routes[2]
	require.Equal(fiber.MethodPut, route3.Method)
	require.Equal("/user/refresh-token", route3.Path)
	require.Len(route3.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
