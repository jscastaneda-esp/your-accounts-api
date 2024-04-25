package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"

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
	request := c.Locals(validation.RequestBody).(*model.CreateRequest)
	id, err := ctrl.app.Create(c.UserContext(), request.Email)
	if err != nil {
		log.Error("Error sign up user:", err)
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
	request := c.Locals(validation.RequestBody).(*model.LoginRequest)
	token, expiresAt, err := ctrl.app.Login(c.UserContext(), request.Email)
	if err != nil {
		log.Error("Error authenticate user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error authenticate user")
	}

	return c.JSON(model.NewLoginResponse(token, expiresAt))
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.UserApp}

	router.Post("/user", validation.RequestBodyValid(model.CreateRequest{}), controller.create)
	router.Post("/login", validation.RequestBodyValid(model.LoginRequest{}), controller.login)
}
