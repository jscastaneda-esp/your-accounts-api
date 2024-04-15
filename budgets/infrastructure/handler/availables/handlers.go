package availables

import (
	"net/http"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/infrastructure/injection"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type controller struct {
	app application.IBudgetAvailableApp
}

// BudgetAvailableCreateHandler godoc
//
//	@Summary		Create available for budget
//	@Description	create a new available for budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			Authorization				header		string							true	"Access token"
//	@Param			request						body		model.CreateAvailableRequest	true	"Available data"
//	@Success		201							{object}	model.CreateAvailableResponse
//	@Failure		400							{string}	string
//	@Failure		401							{string}	string
//	@Failure		422							{string}	string
//	@Failure		500							{string}	string
//	@Router			/api/v1/budget/available/	[post]
func (ctrl *controller) create(c echo.Context) error {
	request := new(model.CreateAvailableRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	id, err := ctrl.app.Create(c.Request().Context(), request.Name, request.BudgetId)
	if err != nil {
		log.Error("Error creating available:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating available")
	}

	return c.JSON(http.StatusCreated, model.NewCreateAvailableResponse(id))
}

func NewRoute(api *echo.Group) {
	controller := &controller{injection.BudgetAvailableApp}

	group := api.Group("/available")
	group.POST("/", controller.create)
}
