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
	originalJwtGenerate    func(ctx context.Context, id string, uid string, email string) (string, time.Time, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uid = "test"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.originalJwtGenerate = jwtGenerate
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockUserRepo = mocks.NewUserRepository(suite.T())
	suite.mockUserTokenRepo = mocks.NewUserTokenRepository(suite.T())
	suite.app = NewUserApp(suite.mockTransactionManager, suite.mockUserRepo, suite.mockUserTokenRepo)
}

func (suite *TestSuite) TestSignUpSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	userExpected := &domain.User{
		ID:    999,
		UID:   user.UID,
		Email: user.Email,
	}
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, nil)
	suite.mockUserRepo.On("Create", ctx, user).Return(userExpected, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.NoError(err)
	require.NotEmpty(res.ID)
	require.Equal(userExpected.ID, res.ID)
	require.Equal(userExpected.UID, res.UID)
	require.Equal(userExpected.Email, res.Email)
}

func (suite *TestSuite) TestSignUpErrorExistsByUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(false, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpExistsByUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(true, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpErrorExistsByEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpExistsByEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(true, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpErrorCreate() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}
	errExpected := errors.New("not created")
	suite.mockUserRepo.On("ExistsByUID", ctx, suite.uid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, nil)
	suite.mockUserRepo.On("Create", ctx, user).Return(nil, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestAuthSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	suite.mockUserRepo.On("FindByUIDAndEmail", ctx, suite.uid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", ctx, mock.Anything).Return(nil, nil)

	token, err := suite.app.Auth(ctx, suite.uid, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestAuthErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("FindByUIDAndEmail", ctx, suite.uid, suite.email).Return(nil, errExpected)

	token, err := suite.app.Auth(ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestAuthErrorJWTGenerate() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	jwtGenerate = func(ctx context.Context, id string, uid string, email string) (string, time.Time, error) {
		return "", time.Time{}, jwt.ErrInvalidToken
	}
	suite.mockUserRepo.On("FindByUIDAndEmail", ctx, suite.uid, suite.email).Return(userExpected, nil)

	token, err := suite.app.Auth(ctx, suite.uid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestAuthErrorCreateUserToken() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UID:   suite.uid,
		Email: suite.email,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uid string, email string) (string, time.Time, error) {
		return suite.token, expiresAt, nil
	}
	errExpected := errors.New("Error constraint")
	suite.mockUserRepo.On("FindByUIDAndEmail", ctx, suite.uid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("Create", ctx, mock.Anything).Return(nil, errExpected)

	token, err := suite.app.Auth(ctx, suite.uid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
