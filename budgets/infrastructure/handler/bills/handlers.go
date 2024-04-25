package bills

import (
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
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
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := c.Locals(validation.RequestBody).(*model.CreateBillRequest)
	id, err := ctrl.app.Create(c.UserContext(), request.Description, request.Category, request.BudgetId)
	if err != nil {
		log.Error("Error creating bill:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating bill")
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewCreateBillResponse(id))
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
func (ctrl *controller) createTransaction(c *fiber.Ctx) error {
	request := c.Locals(validation.RequestBody).(*model.CreateBillTransactionRequest)
	err := ctrl.app.CreateTransaction(c.UserContext(), request.Description, request.Amount, request.BillId)
	if err != nil {
		log.Error("Error creating bill transaction:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating bill transaction")
	}

	return c.SendStatus(fiber.StatusOK)
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.BudgetBillApp}

	group := router.Group("/bill")
	group.Post("/", validation.RequestBodyValid(model.CreateBillRequest{}), controller.create)
	group.Put("/transaction", validation.RequestBodyValid(model.CreateBillTransactionRequest{}), controller.createTransaction)
}
