package handler

import (
	"errors"
	"log"

	"api-your-accounts/shared/domain/validation"
	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/db/tx"
	"api-your-accounts/user/application"
	"api-your-accounts/user/domain"
	"api-your-accounts/user/infrastructure/model"
	"api-your-accounts/user/infrastructure/repository/user"
	"api-your-accounts/user/infrastructure/repository/user_token"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type userController struct {
	app application.IUserApp
}

// CreateUserHandler godoc
//
//	@Summary		Create user
//	@Description	Create user in the system
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.CreateRequest	true	"User data"
//	@Success		200		{object}	model.CreateResponse
//	@Failure		409		{string}	string	"Conflict"
//	@Failure		422		{string}	string	"Unprocessable Entity"
//	@Failure		500		{string}	string	"Internal Server Error"
//	@Router			/user	[post]
func (controller *userController) createUser(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	user := &domain.User{
		UUID:  request.UUID,
		Email: request.Email,
	}

	exists, err := controller.app.Exists(c.UserContext(), user.UUID, user.Email)
	if exists {
		return fiber.NewError(fiber.StatusConflict, "User already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	result, err := controller.app.SignUp(c.UserContext(), user)
	if err != nil {
		log.Println("Error sign up user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
	}

	return c.JSON(model.CreateResponse{
		ID:        result.ID,
		UUID:      result.UUID,
		Email:     result.Email,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
}

// AuthHandler godoc
//
//	@Summary		Authenticate user
//	@Description	create token for access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request		body		model.AuthRequest	true	"Authentication data"
//	@Success		200			{object}	map[string]any
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		422			{string}	string	"Unprocessable Entity"
//	@Failure		500			{string}	string	"Internal Server Error"
//	@Router			/user/auth	[post]
func (controller *userController) auth(c *fiber.Ctx) error {
	request := new(model.AuthRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	token, err := controller.app.Auth(c.UserContext(), request.UUID, request.Email)
	if err != nil {
		log.Println("Error authenticate user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error authenticate user")
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// RefreshTokenHandler godoc
//
//	@Summary		Refresh token of user
//	@Description	refresh token for access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request				body		model.RefreshTokenRequest	true	"Refresh token data"
//	@Success		200					{object}	map[string]any
//	@Failure		400					{string}	string	"BadRequest"
//	@Failure		401					{string}	string	"Unauthorized"
//	@Failure		422					{string}	string	"Unprocessable Entity"
//	@Failure		500					{string}	string	"Internal Server Error"
//	@Router			/user/refresh-token	[put]
func (controller *userController) refreshToken(c *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	if err := c.BodyParser(request); err != nil {
		log.Println("Error request body parser:", err)
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if errors := validation.ValidateStruct(request); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	token, err := controller.app.RefreshToken(c.UserContext(), request.Token, request.UUID, request.Email)
	if err != nil {
		log.Println("Error refresh token user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid data")
		} else if errors.Is(err, application.ErrTokenRefreshed) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error refresh token user")
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func NewRoute(app *fiber.App) {
	tm := tx.NewTransactionManager(db.DB)
	userRepo := user.NewRepository(db.DB)
	userTokenRepo := user_token.NewRepository(db.DB)
	controller := &userController{
		app: application.NewUserApp(tm, userRepo, userTokenRepo),
	}

	group := app.Group("/user")
	group.Post("/", controller.createUser)
	group.Post("/auth", controller.auth)
	group.Put("/refresh-token", controller.refreshToken)
}
