package application

import (
	"context"
	"errors"
	"strings"
	"time"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/users/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTokenRefreshed    = errors.New("token already refreshed")

	jwtGenerate = jwt.JwtGenerate
)

//go:generate mockery --name IUserApp --filename user-app.go
type IUserApp interface {
	Create(ctx context.Context, email string) (uint, error)
	Login(ctx context.Context, email string) (string, time.Time, error)
	DeleteExpired(ctx context.Context) error
}

type userApp struct {
	tm            persistent.TransactionManager
	userRepo      domain.UserRepository
	userTokenRepo domain.UserTokenRepository
}

func (app *userApp) Create(ctx context.Context, email string) (uint, error) {
	email = strings.ToLower(email)
	example := domain.User{
		Email: email,
	}
	exists, err := app.userRepo.ExistsByExample(ctx, example)
	if err != nil {
		return 0, err
	} else if exists {
		return 0, ErrUserAlreadyExists
	}

	user := domain.User{
		Email: email,
	}
	return app.userRepo.Save(ctx, user)
}

func (app *userApp) Login(ctx context.Context, email string) (string, time.Time, error) {
	example := domain.User{
		Email: strings.ToLower(email),
	}
	user, err := app.userRepo.SearchByExample(ctx, example)
	if err != nil {
		return "", time.Time{}, err
	}

	token, expiresAt, err := jwtGenerate(user.ID, strings.ToLower(user.Email))
	if err != nil {
		return "", time.Time{}, err
	}

	userToken := domain.UserToken{
		Token:     token,
		UserId:    user.ID,
		ExpiresAt: expiresAt,
	}
	_, err = app.userTokenRepo.Save(ctx, userToken)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (app *userApp) DeleteExpired(ctx context.Context) error {
	err := app.userTokenRepo.DeleteByExpiresAtGreaterThanNow(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewUserApp(tm persistent.TransactionManager, userRepo domain.UserRepository, userTokenRepo domain.UserTokenRepository) IUserApp {
	return &userApp{tm, userRepo, userTokenRepo}
}
