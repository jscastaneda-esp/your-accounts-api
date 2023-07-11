package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/domain/persistent"
	"your-accounts-api/user/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTokenRefreshed    = errors.New("token already refreshed")

	jwtGenerate = jwt.JwtGenerate
)

//go:generate mockery --name IUserApp --filename user-app.go
type IUserApp interface {
	SignUp(ctx context.Context, user domain.User) (*domain.User, error)
	Auth(ctx context.Context, uid, email string) (string, error)
}

type userApp struct {
	tm            persistent.TransactionManager
	userRepo      domain.UserRepository
	userTokenRepo domain.UserTokenRepository
}

func (app *userApp) SignUp(ctx context.Context, user domain.User) (*domain.User, error) {
	exists, err := app.userRepo.ExistsByUID(ctx, user.UID)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrUserAlreadyExists
	}

	user.Email = strings.ToLower(user.Email)
	exists, err = app.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrUserAlreadyExists
	}

	return app.userRepo.Create(ctx, user)
}

func (app *userApp) Auth(ctx context.Context, uid, email string) (string, error) {
	user, err := app.userRepo.FindByUIDAndEmail(ctx, uid, strings.ToLower(email))
	if err != nil {
		return "", err
	}

	token, expiresAt, err := jwtGenerate(ctx, fmt.Sprint(user.ID), user.UID, strings.ToLower(user.Email))
	if err != nil {
		return "", err
	}

	userToken := domain.UserToken{
		Token:     token,
		UserId:    user.ID,
		ExpiresAt: expiresAt,
	}
	_, err = app.userTokenRepo.Create(ctx, userToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserApp(tm persistent.TransactionManager, userRepo domain.UserRepository, userTokenRepo domain.UserTokenRepository) IUserApp {
	return &userApp{tm, userRepo, userTokenRepo}
}
