package handler

import (
	"errors"
	"fmt"
	"log/slog"

	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"
	"your-accounts-api/users/application"
	"your-accounts-api/users/infrastructure/model"

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
//	@Failure		400		{string}	string
//	@Failure		409		{string}	string
//	@Failure		422		{string}	string
//	@Failure		500		{string}	string
//	@Router			/user	[post]
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	id, err := ctrl.app.Create(c.UserContext(), request.UID, request.Email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sign up user: %v\n", err))
		if errors.Is(err, application.ErrUserAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error sign up user")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewCreateResponse(id))
}

// UserLoginHandler godoc
//
//	@Summary		Authenticate user
//	@Description	create token for access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoginRequest	true	"Authentication data"
//	@Success		200		{object}	model.LoginResponse
//	@Failure		400		{string}	string
//	@Failure		401		{string}	string
//	@Failure		422		{string}	string
//	@Failure		500		{string}	string
//	@Router			/login	[post]
func (ctrl *controller) login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	token, err := ctrl.app.Login(c.UserContext(), request.UID, request.Email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error authenticate user: %v\n", err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error authenticate user")
	}

	return c.JSON(model.NewLoginResponse(token))
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.UserApp}

	router.Post("/user", controller.create)
	router.Post("/login", controller.login)
}
