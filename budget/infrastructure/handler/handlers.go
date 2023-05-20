package handler

import (
	"api-your-accounts/budget/application"
	"api-your-accounts/budget/infrastructure/model"
	"api-your-accounts/budget/infrastructure/repository/budget"
	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/db/persistent"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type controller struct {
	app application.IBudgetApp
}

// BudgetReadByIdHandler godoc
//
//	@Summary		Read budget by ID
//	@Description	read budget by ID
//	@Tags			budget
//	@Produce		json
//	@Param			Authorization		header		string	true	"Access token"
//	@Param			id					path		uint	true	"Budget ID"
//	@Success		200					{array}		model.ReadResponse
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		404					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/budget/{id}	[get]
func (ctrl *controller) readById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Error getting param 'id': %v\n", err)
		return fiber.ErrBadRequest
	}

	budget, err := ctrl.app.FindById(c.UserContext(), uint(id))
	if err != nil {
		log.Printf("Error reading budget by id: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Budget ID not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error reading projects by user")
	}

	return c.JSON(model.ReadResponse{
		ID:               budget.ID,
		Name:             budget.Name,
		Year:             budget.Year,
		Month:            budget.Month,
		FixedIncome:      budget.FixedIncome,
		AdditionalIncome: budget.AdditionalIncome,
		TotalBalance:     budget.TotalBalance,
		Total:            budget.Total,
		EstimatedBalance: budget.EstimatedBalance,
		ProjectId:        budget.ProjectId,
	})
}

func NewRoute(router fiber.Router) {
	tm := persistent.NewTransactionManager(db.DB)
	budgetRepo := budget.NewRepository(db.DB)
	controller := &controller{
		app: application.NewBudgetApp(tm, budgetRepo),
	}

	group := router.Group("/budget")
	group.Get("/:id<min(1)>", controller.readById)
}
