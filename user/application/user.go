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

type IUserApp interface {
	Exists(ctx context.Context, uuid string, email string) (bool, error)
	SignUp(ctx context.Context, user *domain.User) (*domain.User, error)
	Login(ctx context.Context, uuid string, email string) (string, error)
}

type userApp struct {
	repo domain.UserRepository
}

func (app *userApp) Exists(ctx context.Context, uuid string, email string) (bool, error) {
	exists, err := app.repo.ExistsByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}
	if exists {
		return exists, nil
	}

	exists, err = app.repo.ExistsByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (app *userApp) SignUp(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.Email = strings.ToLower(user.Email)
	return app.repo.Create(ctx, user)
}

func (app *userApp) Login(ctx context.Context, uuid string, email string) (string, error) {
	user, err := app.repo.FindByUUIDAndEmail(ctx, uuid, strings.ToLower(email))
	if err != nil {
		return "", err
	}

	token, err := jwtGenerate(ctx, fmt.Sprint(user.ID), user.UUID, strings.ToLower(user.Email))
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserApp(repo domain.UserRepository) IUserApp {
	return &userApp{repo}
}
