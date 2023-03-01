package application

import (
	jwt "api-your-accounts/shared/domain"
	"api-your-accounts/user/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	errFindByUUIDAndEmail error
	errCreate             error
}

func (mock *MockUserRepository) FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*domain.User, error) {
	return &domain.User{
		ID: 999,
	}, mock.errFindByUUIDAndEmail
}

func (mock *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.ID = 999
	return user, mock.errCreate
}

func (mock *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}

type TestSuite struct {
	suite.Suite
	uuid                string
	email               string
	token               string
	originalJwtGenerate func(ctx context.Context, userId string) (string, error)
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "test"
	suite.email = "example@exaple.com"
	suite.token = "<token>"
	suite.originalJwtGenerate = jwtGenerate
}

func (suite *TestSuite) SetupTest() {
	jwtGenerate = suite.originalJwtGenerate
}

func (suite *TestSuite) TestExistsSuccess() {
	require := require.New(suite.T())

	exists, err := Exists(&MockUserRepository{}, context.Background(), suite.uuid, suite.email)
	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsError() {
	require := require.New(suite.T())

	repo := &MockUserRepository{
		errFindByUUIDAndEmail: errors.New("Not exists"),
	}

	exists, err := Exists(repo, context.Background(), suite.uuid, suite.email)
	require.EqualError(repo.errFindByUUIDAndEmail, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestSignUpSuccess() {
	require := require.New(suite.T())

	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}

	res, err := SignUp(&MockUserRepository{}, context.Background(), user)
	require.NoError(err)
	require.NotEmpty(res.ID)
	require.Equal(user.UUID, res.UUID)
	require.Equal(user.Email, res.Email)
}

func (suite *TestSuite) TestSignUpError() {
	require := require.New(suite.T())

	repo := &MockUserRepository{
		errCreate: errors.New("not created"),
	}
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}

	_, err := SignUp(repo, context.Background(), user)
	require.EqualError(repo.errCreate, err.Error())
}

/*
func Login(repo domain.UserRepository, ctx context.Context, uuid string, email string) (string, error) {
	user, err := repo.FindByUUIDAndEmail(ctx, uuid, email)
	if err != nil {
		return "", err
	}

	token, err := sharedD.JwtGenerate(ctx, fmt.Sprint(user.Id))
	if err != nil {
		return "", err
	}

	return token, nil
}
*/

func (suite *TestSuite) TestLoginSuccess() {
	require := require.New(suite.T())

	jwtGenerate = func(ctx context.Context, userId string) (string, error) {
		return suite.token, nil
	}

	token, err := Login(&MockUserRepository{}, context.Background(), suite.uuid, suite.email)
	require.NoError(err)
	require.Equal(suite.token, token)
}

func (suite *TestSuite) TestLoginErrorFind() {
	require := require.New(suite.T())

	repo := &MockUserRepository{
		errFindByUUIDAndEmail: errors.New("Not exists"),
	}

	token, err := Login(repo, context.Background(), suite.uuid, suite.email)
	require.EqualError(repo.errFindByUUIDAndEmail, err.Error())
	require.Empty(token)
}

func (suite *TestSuite) TestLoginErrorJWTGenerate() {
	require := require.New(suite.T())

	jwtGenerate = func(ctx context.Context, userId string) (string, error) {
		return "", jwt.ErrInvalidToken
	}

	token, err := Login(&MockUserRepository{}, context.Background(), suite.uuid, suite.email)
	require.EqualError(jwt.ErrInvalidToken, err.Error())
	require.Empty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
