package application

import (
	sharedDom "api-your-accounts/shared/domain"
	"api-your-accounts/user/domain"
	"context"
	"fmt"
)

var (
	jwtGenerate = sharedDom.JwtGenerate
)

func Exists(repo domain.UserRepository, ctx context.Context, uuid string, email string) (bool, error) {
	_, err := repo.FindByUUIDAndEmail(ctx, uuid, email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func SignUp(repo domain.UserRepository, ctx context.Context, user *domain.User) (*domain.User, error) {
	return repo.Create(ctx, user)
}

func Login(repo domain.UserRepository, ctx context.Context, uuid string, email string) (string, error) {
	user, err := repo.FindByUUIDAndEmail(ctx, uuid, email)
	if err != nil {
		return "", err
	}

	token, err := jwtGenerate(ctx, fmt.Sprint(user.ID))
	if err != nil {
		return "", err
	}

	return token, nil
}
