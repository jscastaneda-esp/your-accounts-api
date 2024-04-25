package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"
	"time"
	mocks_application "your-accounts-api/mocks/shared/application"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	code       domain.CodeLog
	resourceId uint
	app        *fiber.App
	mock       *mocks_application.MockILogApp
}

func (suite *TestSuite) SetupSuite() {
	suite.code = domain.Budget
	suite.resourceId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mock = mocks_application.NewMockILogApp(suite.T())
	injection.LogApp = suite.mock

	suite.app = fiber.New()
	NewRoute(suite.app)
}

func (suite *TestSuite) TestReadLogs200() {
	require := require.New(suite.T())
	result := domain.Log{
		ID:          1,
		Description: "Test",
		ResourceId:  suite.resourceId,
		CreatedAt:   time.Now(),
	}
	suite.mock.On("FindByProject", mock.Anything, suite.code, suite.resourceId).Return([]domain.Log{result}, nil)
	expectedBody, err := json.Marshal([]model.ReadLogsResponse{model.NewReadLogsResponse(result)})
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/log/%d/code/%s", suite.resourceId, suite.code), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
	resp, err := io.ReadAll(response.Body)
	require.NoError(err)
	require.Equal(expectedBody, resp)
}

func (suite *TestSuite) TestReadLogs404_1() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/logs/%d/code/%s", 0, suite.code), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadLogs404_2() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/log/%d/code/%s", suite.resourceId, ""), nil)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusNotFound, response.StatusCode)
}

func (suite *TestSuite) TestReadLogs500() {
	require := require.New(suite.T())
	suite.mock.On("FindByProject", mock.Anything, suite.code, suite.resourceId).Return(nil, gorm.ErrInvalidField)
	expectedErr := []byte("Error reading logs by resource and code")

	request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/log/%d/code/%s", suite.resourceId, suite.code), nil)
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
