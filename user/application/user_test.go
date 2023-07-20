package application

import (
	"context"
	"errors"
	"testing"
	"time"
	"your-accounts-api/shared/domain/jwt"
	mocksShared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/user/domain"
	"your-accounts-api/user/domain/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	uid                    string
	email                  string
	token                  string
	mockTransactionManager *mocksShared.TransactionManager
	mockUserRepo           *mocks.UserRepository
	mockUserTokenRepo      *mocks.UserTokenRepository
	app                    IUserApp
	ctx                    context.Context
	originalJwtGenerate    func(id uint, uid string, email string) (string, time.Time, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uid = "test"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.ctx = context.Background()
	suite.originalJwtGenerate = jwtGenerate
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockUserRepo = mocks.NewUserRepository(suite.T())
	suite.mockUserTokenRepo = mocks.NewUserTokenRepository(suite.T())
	instance = nil
	suite.app = NewUserApp(suite.mockTransactionManager, suite.mockUserRepo, suite.mockUserTokenRepo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, suite.email).Return(false, nil)
	suite.mockUserRepo.On("Create", suite.ctx, user).Return(uint(999), nil)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.NoError(err)
	require.Equal(res, uint(999))
}

func (suite *TestSuite) TestCreateErrorExistsByUID() {
	require := require.New(suite.T())
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(false, errExpected)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateExistsByUID() {
	require := require.New(suite.T())
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(true, nil)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorExistsByEmail() {
	require := require.New(suite.T())
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, suite.email).Return(false, errExpected)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateExistsByEmail() {
	require := require.New(suite.T())
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, suite.email).Return(true, nil)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreate() {
	require := require.New(suite.T())
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	errExpected := errors.New("not created")
	suite.mockUserRepo.On("ExistsByUID", suite.ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, suite.email).Return(false, nil)
	suite.mockUserRepo.On("Create", suite.ctx, user).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestLoginSuccess() {
	require := require.New(suite.T())
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(id uint, uid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	suite.mockUserRepo.On("FindByUIDAndEmail", suite.ctx, suite.uid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", suite.ctx, mock.Anything).Return(uint(0), nil)

	token, err := suite.app.Login(suite.ctx, suite.uid, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestLoginErrorFind() {
	require := require.New(suite.T())
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("FindByUIDAndEmail", suite.ctx, suite.uid, suite.email).Return(nil, errExpected)

	token, err := suite.app.Login(suite.ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestLoginErrorJWTGenerate() {
	require := require.New(suite.T())
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	jwtGenerate = func(id uint, uid string, email string) (string, time.Time, error) {
		return "", time.Time{}, jwt.ErrInvalidToken
	}
	suite.mockUserRepo.On("FindByUIDAndEmail", suite.ctx, suite.uid, suite.email).Return(userExpected, nil)

	token, err := suite.app.Login(suite.ctx, suite.uid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestLoginErrorCreateUserToken() {
	require := require.New(suite.T())
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(id uint, uid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	errExpected := errors.New("Error constraint")
	suite.mockUserRepo.On("FindByUIDAndEmail", suite.ctx, suite.uid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	token, err := suite.app.Login(suite.ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
