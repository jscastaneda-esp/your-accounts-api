package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
	Data string `json:"data" validate:"required"`
}

type TestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *TestSuite) SetupSuite() {
	suite.app = fiber.New()
	suite.app.Post("/", RequestBodyValid(TestStruct{}), func(c *fiber.Ctx) error {
		request := c.Locals(RequestBody).(*TestStruct)
		if request == nil || request.Data == "" {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})
}

func (suite *TestSuite) TestRequestBodyValidSuccess() {
	require := require.New(suite.T())
	test := &TestStruct{
		Data: "Test",
	}
	body, err := json.Marshal(test)
	require.NoError(err)

	request := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewReader(body))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusOK, response.StatusCode)
}

func (suite *TestSuite) TestRequestBodyValidErrorBodyParser() {
	require := require.New(suite.T())

	request := httptest.NewRequest(fiber.MethodPost, "/", nil)
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := suite.app.Test(request)

	require.NoError(err)
	require.NotNil(response)
	require.Equal(fiber.StatusBadRequest, response.StatusCode)
}

func (suite *TestSuite) TestRequestBodyValidErrorValidate() {
	require := require.New(suite.T())
	test := new(TestStruct)
	body, err := json.Marshal(test)
	require.NoError(err)
	validationErrors := []*ErrorResponse{
		{
			Field:      "TestStruct.data",
			Constraint: "required",
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

func TestBodyTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
