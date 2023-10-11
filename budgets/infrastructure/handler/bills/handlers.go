package bills

import (
	"fmt"
	"log/slog"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"

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
	request := new(model.CreateBillRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	id, err := ctrl.app.Create(c.UserContext(), request.Description, request.Category, request.BudgetId)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating bill: %v\n", err))
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating bill")
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewCreateBillResponse(id))
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.BudgetBillApp}

	group := router.Group("/bill")
	group.Post("/", controller.create)
}
