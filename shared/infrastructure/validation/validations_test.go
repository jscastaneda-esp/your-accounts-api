package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
	ID   string `json:"id" validate:"len=10"`
	Name string `json:"-" validate:"required"`
}

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestValidateStructSuccess() {
	require := require.New(suite.T())
	object := TestStruct{
		ID:   "1234567890",
		Name: "Test",
	}

	errors := ValidateStruct(object)

	require.Empty(errors)
}

func (suite *TestSuite) TestValidateStructErrorRequired() {
	require := require.New(suite.T())
	object := TestStruct{
		ID: "1234567890",
	}

	errors := ValidateStruct(object)

	require.NotEmpty(errors)
	require.Len(errors, 1)

	err := errors[0]
	require.Equal("Name", err.Field)
	require.Equal("required", err.Constraint)
}

func (suite *TestSuite) TestValidateStructErrorLen() {
	require := require.New(suite.T())
	object := TestStruct{
		ID:   "invalid",
		Name: "Test",
	}

	errors := ValidateStruct(object)

	require.NotEmpty(errors)
	require.Len(errors, 1)

	err := errors[0]
	require.Equal("id", err.Field)
	require.Equal("len=10", err.Constraint)
}

func (suite *TestSuite) TestValidateStructErrorMultiple() {
	require := require.New(suite.T())
	object := TestStruct{
		ID: "invalid",
	}

	errors := ValidateStruct(object)

	require.NotEmpty(errors)
	require.Len(errors, 2)

	err1 := errors[0]
	require.Equal("id", err1.Field)
	require.Equal("len=10", err1.Constraint)

	err2 := errors[1]
	require.Equal("Name", err2.Field)
	require.Equal("required", err2.Constraint)
}

func (suite *TestSuite) TestValidateVariableSuccess() {
	require := require.New(suite.T())

	errors := ValidateVariable("valid", "required")

	require.Empty(errors)
}

func (suite *TestSuite) TestValidateVariableErrorRequired() {
	require := require.New(suite.T())

	errors := ValidateVariable("", "required")

	err := errors[0]
	require.Empty(err.Field)
	require.Equal("required", err.Constraint)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
