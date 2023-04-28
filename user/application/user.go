package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/user/domain"
	"context"
	"fmt"
	"strings"
)

var (
	jwtGenerate = jwt.JwtGenerate
)

//go:generate mockery --name IUserApp --filename user-app.go
type IUserApp interface {
	Exists(ctx context.Context, uuid string, email string) (bool, error)
	SignUp(ctx context.Context, user *domain.User) (*domain.User, error)
	Login(ctx context.Context, uuid string, email string) (string, error)
}

type userApp struct {
	userRepo      domain.UserRepository
	userTokenRepo domain.UserTokenRepository
}

func (app *userApp) Exists(ctx context.Context, uuid string, email string) (bool, error) {
	exists, err := app.userRepo.ExistsByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}
	if exists {
		return exists, nil
	}

	exists, err = app.userRepo.ExistsByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (app *userApp) SignUp(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.Email = strings.ToLower(user.Email)
	return app.userRepo.Create(ctx, user)
}

func (app *userApp) Login(ctx context.Context, uuid string, email string) (string, error) {
	user, err := app.userRepo.FindByUUIDAndEmail(ctx, uuid, strings.ToLower(email))
	if err != nil {
		return "", err
	}

	token, expiresAt, err := jwtGenerate(ctx, fmt.Sprint(user.ID), user.UUID, strings.ToLower(user.Email))
	if err != nil {
		return "", err
	}

	userToken := &domain.UserToken{
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

func NewUserApp(userRepo domain.UserRepository, userTokenRepo domain.UserTokenRepository) IUserApp {
	return &userApp{userRepo, userTokenRepo}
}
