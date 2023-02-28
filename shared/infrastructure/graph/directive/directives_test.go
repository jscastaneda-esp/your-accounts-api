package directive

import (
	"api-your-accounts/shared/domain"
	"api-your-accounts/shared/infrastructure/graph"
	middleware "api-your-accounts/shared/infrastructure/middleware/auth"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type TestSuite struct {
	suite.Suite
	directives graph.DirectiveRoot
}

func (suite *TestSuite) SetupSuite() {
	suite.directives = GetDirectives()
}

func (suite *TestSuite) SetupTest() {
}

func (suite *TestSuite) TestAuthSuccess() {
	require := require.New(suite.T())

	body := "Hello world!"

	next := func(ctx context.Context) (interface{}, error) {
		return body, nil
	}
	var claims domain.JwtCustomClaim
	ctx := context.WithValue(context.Background(), middleware.CtxAuth, &claims)

	res, err := suite.directives.Auth(ctx, nil, next)
	require.NoError(err)
	require.NotEmpty(res)
	require.Equal(body, res)
}

func (suite *TestSuite) TestAuthErrorNotFoundClaims() {
	require := require.New(suite.T())

	res, err := suite.directives.Auth(context.Background(), nil, nil)
	require.EqualError(&gqlerror.Error{
		Message: "Access Denied",
	}, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestBindingSuccessRequired() {
	require := require.New(suite.T())

	body := "Test"

	next := func(ctx context.Context) (interface{}, error) {
		return body, nil
	}

	res, err := suite.directives.Binding(context.Background(), nil, next, "required")
	require.NoError(err)
	require.NotEmpty(res)
	require.Equal(body, res)
}

func (suite *TestSuite) TestBindingSuccessEmail() {
	require := require.New(suite.T())

	body := "example@example.com"

	next := func(ctx context.Context) (interface{}, error) {
		return body, nil
	}

	res, err := suite.directives.Binding(context.Background(), nil, next, "email")
	require.NoError(err)
	require.NotEmpty(res)
	require.Equal(body, res)
}

func (suite *TestSuite) TestBindingErrorNext() {
	require := require.New(suite.T())

	expectedError := errors.New("mock error")

	next := func(ctx context.Context) (interface{}, error) {
		return nil, expectedError
	}

	res, err := suite.directives.Binding(context.Background(), nil, next, "required")
	require.EqualError(&gqlerror.Error{
		Message: expectedError.Error(),
	}, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestBindingErrorValidation() {
	require := require.New(suite.T())

	next := func(ctx context.Context) (interface{}, error) {
		return "", nil
	}

	res, err := suite.directives.Binding(context.Background(), nil, next, "required")
	require.EqualError(&gqlerror.Error{
		Message: "Key: '' Error:Field validation for '' failed on the 'required' tag",
	}, err.Error())
	require.Nil(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
