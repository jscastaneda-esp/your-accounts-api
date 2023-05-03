package application

import (
	"api-your-accounts/shared/domain/jwt"
	"api-your-accounts/shared/domain/transaction"
	"api-your-accounts/user/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTokenRefreshed = errors.New("token already refreshed")

	jwtGenerate = jwt.JwtGenerate
)

//go:generate mockery --name IUserApp --filename user-app.go
type IUserApp interface {
	Exists(ctx context.Context, uuid, email string) (bool, error)
	SignUp(ctx context.Context, user *domain.User) (*domain.User, error)
	Login(ctx context.Context, uuid, email string) (string, error)
	RefreshToken(ctx context.Context, token, uuid, email string) (string, error)
}

type userApp struct {
	tm            transaction.TransactionManager
	userRepo      domain.UserRepository
	userTokenRepo domain.UserTokenRepository
}

func (app *userApp) Exists(ctx context.Context, uuid, email string) (bool, error) {
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

func (app *userApp) Login(ctx context.Context, uuid, email string) (string, error) {
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

func (app *userApp) RefreshToken(ctx context.Context, token, uuid, email string) (string, error) {
	user, err := app.userRepo.FindByUUIDAndEmail(ctx, uuid, strings.ToLower(email))
	if err != nil {
		return "", err
	}

	oldUserToken, err := app.userTokenRepo.FindByTokenAndUserId(ctx, token, user.ID)
	if err != nil {
		return "", err
	}
	if oldUserToken.RefreshedBy != nil {
		return "", ErrTokenRefreshed
	}

	token, expiresAt, err := jwtGenerate(ctx, fmt.Sprint(user.ID), user.UUID, strings.ToLower(user.Email))
	if err != nil {
		return "", err
	}

	err = app.tm.Transaction(func(tx transaction.Transaction) error {
		repo := app.userTokenRepo.WithTransaction(tx)

		newUserToken := &domain.UserToken{
			Token:     token,
			UserId:    user.ID,
			ExpiresAt: expiresAt,
		}
		newUserToken, err := repo.Create(ctx, newUserToken)
		if err != nil {
			return err
		}

		refreshedAt := time.Now()
		oldUserToken.RefreshedBy = &newUserToken.ID
		oldUserToken.RefreshedAt = &refreshedAt
		_, err = repo.Update(ctx, oldUserToken)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserApp(tm transaction.TransactionManager, userRepo domain.UserRepository, userTokenRepo domain.UserTokenRepository) IUserApp {
	return &userApp{tm, userRepo, userTokenRepo}
}
