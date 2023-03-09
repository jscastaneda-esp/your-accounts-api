package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/user/domain"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	mock.Mock
}

func (mock *MockUserRepository) FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*domain.User, error) {
	args := mock.Called(ctx, uuid, email)
	err := args.Error(1)
	obj := args.Get(0)
	if obj == nil {
		return nil, err
	}
	user, ok := obj.(*domain.User)
	if !ok {
		panic(fmt.Sprintf("assert: arguments: *domain.User(0) failed because object wasn't correct type: %v", obj))
	}
	return user, err
}

func (mock *MockUserRepository) ExistsByUUID(ctx context.Context, uuid string) (bool, error) {
	args := mock.Called(ctx, uuid)
	return args.Bool(0), args.Error(1)
}

func (mock *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := mock.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (mock *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := mock.Called(ctx, user)
	err := args.Error(1)
	obj := args.Get(0)
	if obj == nil {
		return nil, err
	}
	user, ok := obj.(*domain.User)
	if !ok {
		panic(fmt.Sprintf("assert: arguments: *domain.User(0) failed because object wasn't correct type: %v", args.Get(0)))
	}
	return user, err
}

func (mock *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

type TestSuite struct {
	suite.Suite
	uuid                string
	email               string
	token               string
	mock                *MockUserRepository
	app                 IUserApp
	originalJwtGenerate func(ctx context.Context, id string, uuid string, email string) (string, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "test"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.originalJwtGenerate = jwtGenerate
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
	suite.mock = new(MockUserRepository)
	suite.app = NewUserApp(suite.mock)
}

func (suite *TestSuite) TearDownTest() {
	require.True(suite.T(), suite.mock.AssertExpectations(suite.T()))
}

func (suite *TestSuite) TestExistsSuccessUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mock.On("ExistsByUUID", ctx, suite.uuid).Return(true, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsSuccessEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mock.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mock.On("ExistsByEmail", ctx, suite.email).Return(true, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsSuccessNot() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mock.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mock.On("ExistsByEmail", ctx, suite.email).Return(false, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsErrorUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mock.On("ExistsByUUID", ctx, suite.uuid).Return(false, errExpected)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestExistsErrorEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mock.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mock.On("ExistsByEmail", ctx, suite.email).Return(false, errExpected)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestSignUpSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	userExpected := &domain.User{
		ID:    999,
		UUID:  user.UUID,
		Email: user.Email,
	}
	suite.mock.On("Create", ctx, user).Return(userExpected, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.NoError(err)
	require.NotEmpty(res.ID)
	require.Equal(userExpected.ID, res.ID)
	require.Equal(userExpected.UUID, res.UUID)
	require.Equal(userExpected.Email, res.Email)
}

func (suite *TestSuite) TestSignUpError() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	errExpected := errors.New("not created")
	suite.mock.On("Create", ctx, user).Return(nil, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestLoginSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	suite.mock.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, error) {
		return suite.token, nil
	}

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestLoginErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mock.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(nil, errExpected)

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestLoginErrorJWTGenerate() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	suite.mock.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, error) {
		return "", jwt.ErrInvalidToken
	}

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
