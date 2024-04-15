package handler

import (
	"errors"
	"net/http"
	"strconv"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/handler/availables"
	"your-accounts-api/budgets/infrastructure/handler/bills"
	"your-accounts-api/budgets/infrastructure/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/utils/convert"
	"your-accounts-api/shared/infrastructure"
	"your-accounts-api/shared/infrastructure/injection"

	"gorm.io/gorm"
)

type controller struct {
	app application.IBudgetApp
}

// BudgetCreateHandler godoc
//
//	@Summary		Create budget
//	@Description	create a new budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Access token"
//	@Param			request			body		model.CreateRequest	true	"Budget data"
//	@Success		201				{object}	model.CreateResponse
//	@Failure		400				{string}	string
//	@Failure		401				{string}	string
//	@Failure		404				{string}	string
//	@Failure		422				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/budget/	[post]
func (ctrl *controller) create(c echo.Context) error {
	request := new(model.CreateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	userData := getUserData(c)

	var id uint
	var err error
	if request.CloneId == nil {
		id, err = ctrl.app.Create(c.Request().Context(), userData.ID, request.Name)
	} else {
		id, err = ctrl.app.Clone(c.Request().Context(), userData.ID, *request.CloneId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Clone ID %d not found", *request.CloneId)
			return echo.NewHTTPError(http.StatusNotFound, "Error creating budget. Clone ID not found")
		}
	}

	if err != nil {
		log.Error("Error creating budget:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating budget")
	}

	return c.JSON(http.StatusCreated, model.NewCreateResponse(id))
}

// BudgetReadByUserHandler godoc
//
//	@Summary		Read budgets by user
//	@Description	read budgets associated to an user
//	@Tags			budget
//	@Produce		json
//	@Param			Authorization	header		string	true	"Access token"
//	@Param			user			path		uint	true	"User ID"
//	@Success		200				{array}		model.ReadResponse
//	@Failure		400				{string}	string
//	@Failure		401				{string}	string
//	@Failure		404				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/budget/	[get]
func (ctrl *controller) read(c echo.Context) error {
	userData := getUserData(c)

	budgets, err := ctrl.app.FindByUserId(c.Request().Context(), userData.ID)
	if err != nil {
		log.Error("Error reading budgets by user:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading budgets by user")
	}

	response := make([]model.ReadResponse, 0)
	for _, budget := range budgets {
		response = append(response, model.NewReadResponse(budget))
	}

	return c.JSON(http.StatusOK, response)
}

// BudgetReadByIdHandler godoc
//
//	@Summary		Read budget by ID
//	@Description	read budget by ID
//	@Tags			budget
//	@Produce		json
//	@Param			Authorization		header		string	true	"Access token"
//	@Param			id					path		uint	true	"Budget ID"
//	@Success		200					{object}	model.ReadByIDResponse
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		404					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/budget/{id}	[get]
func (ctrl *controller) readById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return echo.ErrBadRequest
	}

	budget, err := ctrl.app.FindById(c.Request().Context(), uint(id))
	if err != nil {
		log.Error("Error reading budget by id:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Budget ID not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading projects by user")
	}

	return c.JSON(http.StatusOK, model.NewReadByIDResponse(budget))
}

// BudgetChangesHandler godoc
//
//	@Summary		Receive changes in budget
//	@Description	receive changes associated to a budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string					true	"Access token"
//	@Param			id					path		uint					true	"Budget ID"
//	@Param			request				body		[]model.ChangeRequest	true	"Changes data"
//	@Success		200					{string}	string
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		404					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{array}		model.ChangeResponse
//	@Router			/api/v1/budget/{id}	[put]
func (ctrl *controller) changes(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return echo.ErrBadRequest
	}

	request := []model.ChangeRequest{}
	validator := c.Echo().Validator.(*infrastructure.CustomValidator)
	if err := validator.ValidateVariable(&request, "min=1,dive,required"); err != nil {
		return err
	}

	changes := []application.Change{}
	for _, item := range request {
		changes = append(changes, application.Change(item))
	}

	results := ctrl.app.Changes(c.Request().Context(), uint(id), changes)
	errs := []model.ChangeResponse{}
	for _, result := range results {
		if result.Err != nil {
			log.Errorf("Error processing change %v in budget: %s", result.Change, result.Err.Error())

			if errors.Is(result.Err, application.ErrIncompleteData) {
				errs = append(errs, model.NewChangeResponse(result.Change, "Incomplete data"))
			} else if errors.Is(result.Err, convert.ErrValueIncompatibleType) {
				errs = append(errs, model.NewChangeResponse(result.Change, "Incompatible data type"))
			} else {
				errs = append(errs, model.NewChangeResponse(result.Change, "Error processing change"))
			}
		}
	}

	if len(errs) > 0 {
		log.Error("Error processing changes in budget")
		return c.JSON(http.StatusInternalServerError, errs)
	}

	return c.NoContent(http.StatusOK)
}

// BudgetDeleteHandler godoc
//
//	@Summary		Delete budget
//	@Description	Delete an budget by ID
//	@Tags			budget
//	@Produce		json
//	@Param			Authorization		header		string	true	"Access token"
//	@Param			id					path		uint	true	"Budget ID"
//	@Success		200					{string}	string
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		404					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/budget/{id}	[delete]
func (ctrl *controller) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return echo.ErrBadRequest
	}

	err = ctrl.app.Delete(c.Request().Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error deleting budget:", err)
			return echo.NewHTTPError(http.StatusNotFound, "Budget ID not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting budget")
	}

	return c.NoContent(http.StatusOK)
}

func getUserData(c echo.Context) *shared.JwtUserClaims {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*shared.JwtUserClaims)
}

func NewRoute(api *echo.Group) {
	controller := &controller{injection.BudgetApp}

	group := api.Group("/budget")
	group.POST("", controller.create)
	group.GET("", controller.read)
	group.GET("/:id", controller.readById)
	group.PUT("/:id", controller.changes)
	group.DELETE("/:id", controller.delete)

	// Additional routes
	availables.NewRoute(group)
	bills.NewRoute(group)
}
