package application

import (
	"context"
	"errors"
	"testing"
	"time"
	mocks_persistent "your-accounts-api/mocks/shared/domain/persistent"
	mocks_domain "your-accounts-api/mocks/users/domain"
	"your-accounts-api/shared/infrastructure/config"
	"your-accounts-api/users/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	email                  string
	token                  string
	mockTransactionManager *mocks_persistent.MockTransactionManager
	mockUserRepo           *mocks_domain.MockUserRepository
	mockUserTokenRepo      *mocks_domain.MockUserTokenRepository
	app                    IUserApp
	ctx                    context.Context
	originalJwtGenerate    func(id uint) (string, time.Time, error)
	originalJwtSecret      []byte
}

func (suite *TestSuite) SetupSuite() {
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.ctx = context.Background()
	suite.originalJwtGenerate = jwtGenerate
	suite.originalJwtSecret = config.JWT_SECRET
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
	config.JWT_SECRET = suite.originalJwtSecret
	suite.mockTransactionManager = mocks_persistent.NewMockTransactionManager(suite.T())
	suite.mockUserRepo = mocks_domain.NewMockUserRepository(suite.T())
	suite.mockUserTokenRepo = mocks_domain.NewMockUserTokenRepository(suite.T())
	suite.app = NewUserApp(suite.mockTransactionManager, suite.mockUserRepo, suite.mockUserTokenRepo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	user := domain.User{
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(false, nil)
	suite.mockUserRepo.On("Save", suite.ctx, user).Return(uint(999), nil)

	res, err := suite.app.Create(suite.ctx, suite.email)

	require.NoError(err)
	require.Equal(res, uint(999))
}

func (suite *TestSuite) TestCreateErrorExists() {
	require := require.New(suite.T())
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(false, errExpected)

	res, err := suite.app.Create(suite.ctx, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateExistsTrue() {
	require := require.New(suite.T())
	suite.mockUserRepo.On("ExistsByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(true, nil)

	res, err := suite.app.Create(suite.ctx, suite.email)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreate() {
	require := require.New(suite.T())
	user := domain.User{
		Email: suite.email,
	}
	errExpected := errors.New("not created")
	suite.mockUserRepo.On("ExistsByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(false, nil)
	suite.mockUserRepo.On("Save", suite.ctx, user).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestLoginSuccessOri() {
	require := require.New(suite.T())
	userExpected := domain.User{
		ID:    999,
		Email: suite.email,
	}
	suite.mockUserRepo.On("SearchByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)

	token, expires, err := suite.app.Login(suite.ctx, suite.email)

	require.NoError(err)
	require.NotEmpty(token)
	require.NotNil(expires)
}

func (suite *TestSuite) TestLoginSuccessMock() {
	require := require.New(suite.T())
	userExpected := domain.User{
		ID:    999,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(id uint) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	suite.mockUserRepo.On("SearchByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)

	token, expires, err := suite.app.Login(suite.ctx, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
	require.Equal(expiresAt, expires)
}

func (suite *TestSuite) TestLoginErrorFind() {
	require := require.New(suite.T())
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("SearchByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(domain.User{}, errExpected)

	token, expires, err := suite.app.Login(suite.ctx, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
	require.Empty(expires)
}

func (suite *TestSuite) TestLoginErrorJWTGenerateMock() {
	require := require.New(suite.T())
	userExpected := domain.User{
		ID:    999,
		Email: suite.email,
	}
	jwtGenerate = func(id uint) (string, time.Time, error) {
		return "", time.Time{}, jwt.ErrInvalidKey
	}
	suite.mockUserRepo.On("SearchByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(userExpected, nil)

	token, expires, err := suite.app.Login(suite.ctx, suite.email)

	require.EqualError(jwt.ErrInvalidKey, err.Error())
	require.Empty(token)
	require.Empty(expires)
}

func (suite *TestSuite) TestLoginErrorCreateUserToken() {
	require := require.New(suite.T())
	userExpected := domain.User{
		ID:    999,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(id uint) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	errExpected := errors.New("Error constraint")
	suite.mockUserRepo.On("SearchByExample", suite.ctx, domain.User{
		Email: suite.email,
	}).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	token, expires, err := suite.app.Login(suite.ctx, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
	require.Empty(expires)
}

func (suite *TestSuite) TestDeleteExpiredSuccess() {
	require := require.New(suite.T())
	suite.mockUserTokenRepo.On("DeleteByExpiresAtGreaterThanNow", suite.ctx).Return(nil)

	err := suite.app.DeleteExpired(suite.ctx)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteExpiredError() {
	require := require.New(suite.T())
	suite.mockUserTokenRepo.On("DeleteByExpiresAtGreaterThanNow", suite.ctx).Return(gorm.ErrRecordNotFound)

	err := suite.app.DeleteExpired(suite.ctx)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
