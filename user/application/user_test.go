package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/domain/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	uuid                string
	email               string
	token               string
	mockUserRepo        *mocks.UserRepository
	mockUserTokenRepo   *mocks.UserTokenRepository
	app                 IUserApp
	originalJwtGenerate func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "test"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.originalJwtGenerate = jwtGenerate
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
	suite.mockUserRepo = mocks.NewUserRepository(suite.T())
	suite.mockUserTokenRepo = mocks.NewUserTokenRepository(suite.T())
	suite.app = NewUserApp(suite.mockUserRepo, suite.mockUserTokenRepo)
}

func (suite *TestSuite) TestExistsSuccessUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(true, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsSuccessEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(true, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsSuccessNot() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, nil)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsErrorUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, errExpected)

	exists, err := suite.app.Exists(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestExistsErrorEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, errExpected)

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
	suite.mockUserRepo.On("Create", ctx, user).Return(userExpected, nil)

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
	suite.mockUserRepo.On("Create", ctx, user).Return(nil, errExpected)

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
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	userToken := &domain.UserToken{
		Token:     suite.token,
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", ctx, userToken).Return(nil, nil)

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestLoginErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(nil, errExpected)

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
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return "", time.Time{}, jwt.ErrInvalidToken
	}
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestLoginErrorCreateUserToken() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	userToken := &domain.UserToken{
		Token:     suite.token,
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	errExpected := errors.New("Error constraint")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", ctx, userToken).Return(nil, errExpected)

	token, err := suite.app.Login(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
