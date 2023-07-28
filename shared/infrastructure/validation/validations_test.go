package validation

import (
	"encoding/json"
	"testing"
	"your-accounts-api/shared/domain/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TestStruct struct {
	Data string `json:"data" validate:"required"`
}

type TestSuite struct {
	suite.Suite
	fastCtx *fasthttp.RequestCtx
	ctx     *fiber.Ctx
}

func (suite *TestSuite) SetupSuite() {
	suite.fastCtx = new(fasthttp.RequestCtx)

	app := fiber.New()
	suite.ctx = app.AcquireCtx(suite.fastCtx)
}

func (suite *TestSuite) SetupTest() {
	suite.fastCtx.Request.Reset()
	suite.fastCtx.Response.Reset()
}

func (suite *TestSuite) TestValidateSuccess() {
	require := require.New(suite.T())
	test := &TestStruct{
		Data: "Test",
	}
	body, err := json.Marshal(test)
	require.NoError(err)
	suite.fastCtx.Request.Header.SetContentType(fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)
	mapper := new(TestStruct)

	result := Validate(suite.ctx, mapper)

	require.True(result)
	require.Equal(test.Data, mapper.Data)
	require.Equal(fiber.StatusOK, suite.fastCtx.Response.StatusCode())
}

func (suite *TestSuite) TestValidateErrorBodyParser() {
	require := require.New(suite.T())
	suite.fastCtx.Request.Header.SetContentType(fiber.MIMEApplicationJSON)
	mapper := new(TestStruct)

	result := Validate(suite.ctx, mapper)

	require.False(result)
	require.Empty(mapper.Data)
	require.Equal(fiber.StatusBadRequest, suite.fastCtx.Response.StatusCode())
}

func (suite *TestSuite) TestValidateErrorValidate() {
	require := require.New(suite.T())
	test := new(TestStruct)
	body, err := json.Marshal(test)
	require.NoError(err)
	suite.fastCtx.Request.Header.SetContentType(fiber.MIMEApplicationJSON)
	suite.fastCtx.Request.SetBody(body)
	mapper := new(TestStruct)
	validationErrors := []*validation.ErrorResponse{
		{
			Field:      "data",
			Constraint: "required",
		},
	}
	expectedBody, err := json.Marshal(validationErrors)
	require.NoError(err)

	result := Validate(suite.ctx, mapper)

	require.False(result)
	require.Empty(mapper.Data)
	require.Equal(fiber.StatusUnprocessableEntity, suite.fastCtx.Response.StatusCode())
	require.Equal(expectedBody, suite.fastCtx.Response.Body())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
