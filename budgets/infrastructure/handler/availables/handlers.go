package availables

import (
	"log"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/model"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2"
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
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateAvailableRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	id, err := ctrl.app.Create(c.UserContext(), request.Name, request.BudgetId)
	if err != nil {
		log.Printf("Error creating available: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating available")
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewCreateAvailableResponse(id))
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.BudgetAvailableApp}

	group := router.Group("/available")
	group.Post("/", controller.create)
}
