package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/user/domain"
	"context"
	"fmt"
)

var (
	jwtGenerate = jwt.JwtGenerate
)

func Exists(repo domain.UserRepository, ctx context.Context, uuid string, email string) (bool, error) {
	exists, err := repo.ExistsByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}
	if exists {
		return exists, nil
	}

	exists, err = repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func SignUp(repo domain.UserRepository, ctx context.Context, user *domain.User) (*domain.User, error) {
	return repo.Create(ctx, user)
}

func Login(repo domain.UserRepository, ctx context.Context, uuid string, email string) (string, error) {
	user, err := repo.FindByUUIDAndEmail(ctx, uuid, email)
	if err != nil {
		return "", err
	}

	token, err := jwtGenerate(ctx, fmt.Sprint(user.ID), user.UUID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
