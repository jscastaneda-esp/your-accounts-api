package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/shared/domain/persistent"
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/domain/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	uuid                   string
	email                  string
	token                  string
	mockTransactionManager *mocksShared.TransactionManager
	mockUserRepo           *mocks.UserRepository
	mockUserTokenRepo      *mocks.UserTokenRepository
	app                    IUserApp
	originalJwtGenerate    func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "test"
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
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	userExpected := &domain.User{
		ID:    999,
		UUID:  user.UUID,
		Email: user.Email,
	}
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, nil)
	suite.mockUserRepo.On("Create", ctx, user).Return(userExpected, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.NoError(err)
	require.NotEmpty(res.ID)
	require.Equal(userExpected.ID, res.ID)
	require.Equal(userExpected.UUID, res.UUID)
	require.Equal(userExpected.Email, res.Email)
}

func (suite *TestSuite) TestSignUpErrorExistsByUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpExistsByUUID() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(true, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpErrorExistsByEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(false, errExpected)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpExistsByEmail() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
	suite.mockUserRepo.On("ExistsByEmail", ctx, suite.email).Return(true, nil)

	res, err := suite.app.SignUp(ctx, user)

	require.EqualError(ErrUserAlreadyExists, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestSignUpErrorCreate() {
	require := require.New(suite.T())
	ctx := context.Background()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}
	errExpected := errors.New("not created")
	suite.mockUserRepo.On("ExistsByUUID", ctx, suite.uuid).Return(false, nil)
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

	token, err := suite.app.Auth(ctx, suite.uuid, suite.email)

	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestAuthErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(nil, errExpected)

	token, err := suite.app.Auth(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestAuthErrorJWTGenerate() {
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

	token, err := suite.app.Auth(ctx, suite.uuid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestAuthErrorCreateUserToken() {
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

	token, err := suite.app.Auth(ctx, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	oldUserTokenExpected := &domain.UserToken{
		ID:     999,
		Token:  suite.token,
		UserId: userExpected.ID,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token + "New", expiresAt, nil
	}
	newUserToken := &domain.UserToken{
		Token:     suite.token + "New",
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	newUserTokenExpected := &domain.UserToken{
		ID:        1000,
		Token:     suite.token + "New",
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockUserTokenRepo.On("WithTransaction", nil).Return(suite.mockUserTokenRepo)
	suite.mockUserTokenRepo.On("Create", ctx, newUserToken).Return(newUserTokenExpected, nil)
	suite.mockUserTokenRepo.On("Update", ctx, mock.AnythingOfType("*domain.UserToken")).Return(nil, nil)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.NoError(err)
	require.Equal(suite.token+"New", token)
}

func (suite *TestSuite) TestRefreshTokenErrorUserNotFound() {
	require := require.New(suite.T())
	ctx := context.Background()
	errExpected := errors.New("Not exists user")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(nil, errExpected)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorFindToken() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	errExpected := errors.New("Not exists token")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(nil, errExpected)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorTokenAlreadyRefreshed() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	refreshedBy := uint(1000)
	oldUserTokenExpected := &domain.UserToken{
		ID:          999,
		Token:       suite.token,
		UserId:      userExpected.ID,
		RefreshedBy: &refreshedBy,
	}
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(ErrTokenRefreshed, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorJWTGenerate() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	oldUserTokenExpected := &domain.UserToken{
		ID:     999,
		Token:  suite.token,
		UserId: userExpected.ID,
	}
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return "", time.Time{}, jwt.ErrInvalidToken
	}
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorTransaction() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	oldUserTokenExpected := &domain.UserToken{
		ID:     999,
		Token:  suite.token,
		UserId: userExpected.ID,
	}
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token + "New", time.Time{}, nil
	}
	errExpected := errors.New("Error in transaction")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorCreateNewToken() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	oldUserTokenExpected := &domain.UserToken{
		ID:     999,
		Token:  suite.token,
		UserId: userExpected.ID,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token + "New", expiresAt, nil
	}
	newUserToken := &domain.UserToken{
		Token:     suite.token + "New",
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	errExpected := errors.New("Error in creation")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockUserTokenRepo.On("WithTransaction", nil).Return(suite.mockUserTokenRepo)
	suite.mockUserTokenRepo.On("Create", ctx, newUserToken).Return(nil, errExpected)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestRefreshTokenErrorUpdateExistsToken() {
	require := require.New(suite.T())
	ctx := context.Background()
	userExpected := &domain.User{
		ID:    999,
		UUID:  suite.uuid,
		Email: suite.email,
	}
	oldUserTokenExpected := &domain.UserToken{
		ID:     999,
		Token:  suite.token,
		UserId: userExpected.ID,
	}
	expiresAt := time.Now().Add(1 * time.Hour)
	jwtGenerate = func(ctx context.Context, id string, uuid string, email string) (string, time.Time, error) {
		return suite.token + "New", expiresAt, nil
	}
	newUserToken := &domain.UserToken{
		Token:     suite.token + "New",
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	newUserTokenExpected := &domain.UserToken{
		ID:        1000,
		Token:     suite.token + "New",
		UserId:    userExpected.ID,
		ExpiresAt: expiresAt,
	}
	errExpected := errors.New("Error updating")
	suite.mockUserRepo.On("FindByUUIDAndEmail", ctx, suite.uuid, suite.email).Return(userExpected, nil)
	suite.mockUserTokenRepo.On("FindByTokenAndUserId", ctx, suite.token, userExpected.ID).Return(oldUserTokenExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockUserTokenRepo.On("WithTransaction", nil).Return(suite.mockUserTokenRepo)
	suite.mockUserTokenRepo.On("Create", ctx, newUserToken).Return(newUserTokenExpected, nil)
	suite.mockUserTokenRepo.On("Update", ctx, mock.AnythingOfType("*domain.UserToken")).Return(nil, errExpected)

	token, err := suite.app.RefreshToken(ctx, suite.token, suite.uuid, suite.email)

	require.EqualError(errExpected, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
