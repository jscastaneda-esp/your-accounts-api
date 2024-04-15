package bills

import (
	"net/http"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/infrastructure/injection"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type controller struct {
	app application.IBudgetBillApp
}

// BudgetBillCreateHandler godoc
//
//	@Summary		Create bill for budget
//	@Description	create a new bill for budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			Authorization			header		string					true	"Access token"
//	@Param			request					body		model.CreateBillRequest	true	"Bill data"
//	@Success		201						{object}	model.CreateBillResponse
//	@Failure		400						{string}	string
//	@Failure		401						{string}	string
//	@Failure		422						{string}	string
//	@Failure		500						{string}	string
//	@Router			/api/v1/budget/bill/	[post]
func (ctrl *controller) create(c echo.Context) error {
	request := new(model.CreateBillRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	id, err := ctrl.app.Create(c.Request().Context(), request.Description, request.Category, request.BudgetId)
	if err != nil {
		log.Error("Error creating bill:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating bill")
	}

	return c.JSON(http.StatusCreated, model.NewCreateBillResponse(id))
}

// BudgetBillCreateTransactionHandler godoc
//
//	@Summary		Create transaction for bill of the budget
//	@Description	create a new transaction for bill of the budget
//	@Tags			budget
//	@Accept			json
//	@Produce		json
//	@Param			Authorization					header		string								true	"Access token"
//	@Param			request							body		model.CreateBillTransactionRequest	true	"Bill transaction data"
//	@Success		200								{string}	string
//	@Failure		400								{string}	string
//	@Failure		401								{string}	string
//	@Failure		422								{string}	string
//	@Failure		500								{string}	string
//	@Router			/api/v1/budget/bill/transaction	[put]
func (ctrl *controller) createTransaction(c echo.Context) error {
	request := new(model.CreateBillTransactionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	err := ctrl.app.CreateTransaction(c.Request().Context(), request.Description, request.Amount, request.BillId)
	if err != nil {
		log.Error("Error creating bill transaction:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating bill transaction")
	}

	return c.NoContent(http.StatusOK)
}

func NewRoute(api *echo.Group) {
	controller := &controller{injection.BudgetBillApp}

	group := api.Group("/bill")
	group.POST("/", controller.create)
	group.PUT("/transaction", controller.createTransaction)
}
