package handler

import (
	"errors"
	"log"

	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/db/tx"
	"api-your-accounts/shared/infrastructure/validation"
	"api-your-accounts/user/application"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/model"
	"api-your-accounts/user/infrastructure/repository/user"
	"api-your-accounts/user/infrastructure/repository/user_token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type controller struct {
	app application.IUserApp
}

// UserCreateHandler godoc
//
//	@Summary		Create user
//	@Description	Create user in the system
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateRequest	true	"User data"
//	@Success		201		{object}	model.CreateResponse
//	@Failure		409		{string}	string
//	@Failure		422		{string}	string
//	@Failure		500		{string}	string
//	@Router			/user	[post]
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	user := &domain.User{
		UUID:  request.UUID,
		Email: request.Email,
	}

	exists, err := ctrl.app.Exists(c.UserContext(), user.UUID, user.Email)
	if exists {
		return fiber.NewError(fiber.StatusConflict, "User already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	result, err := ctrl.app.SignUp(c.UserContext(), user)
	if err != nil {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreateResponse{
		ID:        result.ID,
		UUID:      result.UUID,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
}

// UserAuthHandler godoc
//
//	@Summary		Authenticate user
//	@Description	create token for access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request		body		model.AuthRequest	true	"Authentication data"
//	@Success		200			{object}	model.AuthResponse
//	@Failure		401			{string}	string
//	@Failure		422			{string}	string
//	@Failure		500			{string}	string
//	@Router			/user/auth	[post]
func (ctrl *controller) auth(c *fiber.Ctx) error {
	request := new(model.AuthRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	token, err := ctrl.app.Auth(c.UserContext(), request.UUID, request.Email)
	if err != nil {
		log.Println("Error authenticate user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error authenticate user")
	}

	return c.JSON(model.AuthResponse{
		Token: token,
	})
}

// UserRefreshTokenHandler godoc
//
//	@Summary		Refresh token of user
//	@Description	refresh token for access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request				body		model.RefreshTokenRequest	true	"Refresh token data"
//	@Success		200					{object}	model.RefreshTokenResponse
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{string}	string
//	@Router			/user/refresh-token	[put]
func (ctrl *controller) refreshToken(c *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	token, err := ctrl.app.RefreshToken(c.UserContext(), request.Token, request.UUID, request.Email)
	if err != nil {
		log.Println("Error refresh token user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid data")
		} else if errors.Is(err, application.ErrTokenRefreshed) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error refresh token user")
	}

	return c.JSON(model.RefreshTokenResponse{
		AuthResponse: model.AuthResponse{
			Token: token,
		},
	})
}

func NewRoute(router fiber.Router) {
	tm := tx.NewTransactionManager(db.DB)
	userRepo := user.NewRepository(db.DB)
	userTokenRepo := user_token.NewRepository(db.DB)
	controller := &controller{
		app: application.NewUserApp(tm, userRepo, userTokenRepo),
	}

	group := router.Group("/user")
	group.Post("/", controller.create)
	group.Post("/auth", controller.auth)
	group.Put("/refresh-token", controller.refreshToken)
}
