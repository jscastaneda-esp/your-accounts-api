package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	"your-accounts-api/project/application"
	"your-accounts-api/project/application/mocks"
	"your-accounts-api/project/domain"
	"your-accounts-api/project/infrastructure/model"
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	userId     uint
	typeBudget domain.ProjectType
	cloneId    uint
	projectId  uint
	app        *fiber.App
	mock       *mocks.IProjectApp
}

func (suite *TestSuite) SetupSuite() {
	suite.userId = 1
	suite.typeBudget = domain.TypeBudget
	suite.cloneId = 1
	suite.projectId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks.NewIProjectApp(suite.T())
	ctrl := &controller{
		app: suite.mock,
	}

	suite.app = fiber.New()
	suite.app.Post("/", ctrl.create)
	suite.app.Get("/:user<min(1)>", ctrl.readByUser)
	suite.app.Get("/logs/:id<min(1)>", ctrl.readLogs)
	suite.app.Delete("/:id<min(1)>", ctrl.delete)
}

func (suite *TestSuite) TestCreate201() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name:   "Test",
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything).Return(suite.projectId, nil)
	expectedBody, err := json.Marshal(model.CreateResponse{
		ID: suite.projectId,
	})
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

func (suite *TestSuite) TestCreateClone201() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		CloneId: &suite.cloneId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Clone", mock.Anything, *requestBody.CloneId).Return(suite.projectId, nil)
	expectedBody, err := json.Marshal(model.CreateResponse{
		ID: suite.projectId,
	})
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

func (suite *TestSuite) TestCreate404() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		CloneId: &suite.cloneId,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Clone", mock.Anything, *requestBody.CloneId).Return(uint(0), gorm.ErrRecordNotFound)
	expectedErr := []byte("Error creating project. Clone ID not found")

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

func (suite *TestSuite) TestCreate422() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name:   "Cupidatat ullamco voluptate non aute consequat fugiat.",
		UserId: suite.userId,
		Type:   suite.typeBudget,
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

func (suite *TestSuite) TestCreate500() {
	require := require.New(suite.T())
	requestBody := model.CreateRequest{
		Name:   "Test",
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(err)
	suite.mock.On("Create", mock.Anything, mock.Anything).Return(uint(0), gorm.ErrInvalidField)
	expectedErr := []byte("Error creating project")

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

func (suite *TestSuite) TestReadByUser200() {
	require := require.New(suite.T())
	result := &application.FindByUserRecord{
		ID:   suite.projectId,
		Name: "Test",
		Type: suite.typeBudget,
		Data: make(map[string]any),
	}
	suite.mock.On("FindByUser", mock.Anything, suite.userId).Return([]*application.FindByUserRecord{result}, nil)
	expectedBody, err := json.Marshal([]model.ReadResponse{
		{
			ID:   result.ID,
			Name: result.Name,
			Type: result.Type,
			Data: result.Data,
		},
	})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.userId), nil)
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

func (suite *TestSuite) TestReadByUser500() {
	require := require.New(suite.T())
	suite.mock.On("FindByUser", mock.Anything, suite.userId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading projects by user")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%d", suite.userId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusInternalServerError, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedErr, resp)
}

func (suite *TestSuite) TestReadLogs200() {
	require := require.New(suite.T())
	result := &domain.ProjectLog{
		ID:          1,
		Description: "Test",
		ProjectId:   suite.projectId,
		CreatedAt:   time.Now(),
	}
	suite.mock.On("FindLogsByProject", mock.Anything, suite.projectId).Return([]*domain.ProjectLog{result}, nil)
	expectedBody, err := json.Marshal([]model.ReadLogsResponse{
		{
			ID:          result.ID,
			Description: result.Description,
			CreatedAt:   result.CreatedAt,
		},
	})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d", suite.projectId), nil)
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
	suite.mock.On("FindLogsByProject", mock.Anything, suite.projectId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading logs by project")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d", suite.projectId), nil)
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
	suite.mock.On("Delete", mock.Anything, suite.projectId).Return(nil)

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.projectId), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestDelete404() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", 0), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestDelete404ErrorDelete() {
	require := require.New(suite.T())
	suite.mock.On("Delete", mock.Anything, suite.projectId).Return(gorm.ErrRecordNotFound)
	expectedErr := []byte("Project ID not found")

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.projectId), nil)
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
	suite.mock.On("Delete", mock.Anything, suite.projectId).Return(gorm.ErrInvalidField)
	expectedErr := []byte("Error deleting project")

	request := httptest.NewRequest(fiber.MethodDelete, fmt.Sprintf("/%d", suite.projectId), nil)
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
	require.Equal("/project/:user<min(1)>", route1.Path)
	require.Len(route1.Handlers, 1)

	route2 := routes[1]
	require.Equal(fiber.MethodGet, route2.Method)
	require.Equal("/project/logs/:id<min(1)>", route2.Path)
	require.Len(route2.Handlers, 1)

	route3 := routes[2]
	require.Equal(fiber.MethodHead, route3.Method)
	require.Equal("/project/:user<min(1)>", route3.Path)
	require.Len(route3.Handlers, 1)

	route4 := routes[3]
	require.Equal(fiber.MethodHead, route4.Method)
	require.Equal("/project/logs/:id<min(1)>", route4.Path)
	require.Len(route4.Handlers, 1)

	route5 := routes[4]
	require.Equal(fiber.MethodPost, route5.Method)
	require.Equal("/project/", route5.Path)
	require.Len(route5.Handlers, 1)

	route6 := routes[5]
	require.Equal(fiber.MethodDelete, route6.Method)
	require.Equal("/project/:id<min(1)>", route6.Path)
	require.Len(route6.Handlers, 1)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
