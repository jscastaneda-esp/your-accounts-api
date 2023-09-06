package application

import (
	"context"
	"errors"
	"strings"
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
	Create(ctx context.Context, uid, email string) (uint, error)
	Login(ctx context.Context, uid, email string) (string, error)
}

type userApp struct {
	tm            persistent.TransactionManager
	userRepo      domain.UserRepository
	userTokenRepo domain.UserTokenRepository
}

func (app *userApp) Create(ctx context.Context, uid, email string) (uint, error) {
	example := domain.User{
		UID: uid,
	}
	exists, err := app.userRepo.ExistsByExample(ctx, example)
	if err != nil {
		return 0, err
	} else if exists {
		return 0, ErrUserAlreadyExists
	}

	email = strings.ToLower(email)
	example = domain.User{
		Email: email,
	}
	exists, err = app.userRepo.ExistsByExample(ctx, example)
	if err != nil {
		return 0, err
	} else if exists {
		return 0, ErrUserAlreadyExists
	}

	user := domain.User{
		UID:   uid,
		Email: email,
	}
	return app.userRepo.Save(ctx, user)
}

func (app *userApp) Login(ctx context.Context, uid, email string) (string, error) {
	example := domain.User{
		UID:   uid,
		Email: strings.ToLower(email),
	}
	user, err := app.userRepo.SearchByExample(ctx, example)
	if err != nil {
		return "", err
	}

	token, expiresAt, err := jwtGenerate(user.ID, user.UID, strings.ToLower(user.Email))
	if err != nil {
		return "", err
	}

	userToken := domain.UserToken{
		Token:     token,
		UserId:    user.ID,
		ExpiresAt: expiresAt,
	}
	_, err = app.userTokenRepo.Save(ctx, userToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserApp(tm persistent.TransactionManager, userRepo domain.UserRepository, userTokenRepo domain.UserTokenRepository) IUserApp {
	return &userApp{tm, userRepo, userTokenRepo}
}
