package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/users/application"
	"your-accounts-api/users/infrastructure/model"

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
func (ctrl *controller) create(c echo.Context) error {
	request := new(model.CreateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	id, err := ctrl.app.Create(c.Request().Context(), request.Email)
	if err != nil {
		log.Error("Error sign up user:", err)
		if errors.Is(err, application.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error sign up user")
		}
	}

	return c.JSON(http.StatusCreated, model.NewCreateResponse(id))
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
func (ctrl *controller) login(c echo.Context) error {
	request := new(model.LoginRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	token, expiresAt, err := ctrl.app.Login(c.Request().Context(), request.Email)
	if err != nil {
		log.Error("Error authenticate user:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Error authenticate user")
	}

	return c.JSON(http.StatusOK, model.NewLoginResponse(token, expiresAt))
}

func NewRoute(api *echo.Echo) {
	controller := &controller{injection.UserApp}

	api.POST("/user", controller.create)
	api.POST("/login", controller.login)
}
