package handler

import (
	"errors"
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/handler/availables"
	"your-accounts-api/budgets/infrastructure/handler/bills"
	"your-accounts-api/budgets/infrastructure/model"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"

	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/utils/convert"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2"
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
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := c.Locals(validation.RequestBody).(*model.CreateRequest)
	userData := getUserData(c)

	var id uint
	var err error
	if request.CloneId == nil {
		id, err = ctrl.app.Create(c.UserContext(), userData.ID, request.Name)
	} else {
		id, err = ctrl.app.Clone(c.UserContext(), userData.ID, *request.CloneId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Clone ID %d not found", *request.CloneId)
			return fiber.NewError(fiber.StatusNotFound, "Error creating budget. Clone ID not found")
		}
	}

	if err != nil {
		log.Error("Error creating budget:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating budget")
	}

	return c.Status(fiber.StatusCreated).JSON(model.NewCreateResponse(id))
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
func (ctrl *controller) read(c *fiber.Ctx) error {
	userData := getUserData(c)

	budgets, err := ctrl.app.FindByUserId(c.UserContext(), userData.ID)
	if err != nil {
		log.Error("Error reading budgets by user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading budgets by user")
	}

	response := make([]model.ReadResponse, 0)
	for _, budget := range budgets {
		response = append(response, model.NewReadResponse(budget))
	}

	return c.JSON(response)
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
func (ctrl *controller) readById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	budget, err := ctrl.app.FindById(c.UserContext(), uint(id))
	if err != nil {
		log.Error("Error reading budget by id:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Budget ID not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error reading projects by user")
	}

	return c.JSON(model.NewReadByIDResponse(budget))
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
//	@Param			request				body		model.ChangesRequest	true	"Changes data"
//	@Success		200					{string}	string
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		404					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{object}	model.ChangesResponse
//	@Router			/api/v1/budget/{id}	[put]
func (ctrl *controller) changes(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	request := c.Locals(validation.RequestBody).(*model.ChangesRequest)
	changes := []application.Change{}
	for _, item := range request.Changes {
		changes = append(changes, application.Change(item))
	}

	results := ctrl.app.Changes(c.UserContext(), uint(id), changes)
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
		return c.Status(fiber.StatusInternalServerError).JSON(model.ChangesResponse{
			Changes: errs,
		})
	}

	return c.SendStatus(fiber.StatusOK)
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
func (ctrl *controller) delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	err = ctrl.app.Delete(c.UserContext(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Error deleting budget:", err)
			return fiber.NewError(fiber.StatusNotFound, "Budget ID not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting budget")
	}

	return c.SendStatus(fiber.StatusOK)
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.BudgetApp}

	group := router.Group("/budget")
	group.Post("/", validation.RequestBodyValid(model.CreateRequest{}), controller.create)
	group.Get("/", controller.read)
	group.Get("/:id<min(1)>", controller.readById)
	group.Put("/:id<min(1)>/changes", validation.RequestBodyValid(model.ChangesRequest{}), controller.changes)
	group.Delete("/:id<min(1)>", controller.delete)

	// Additional routes
	availables.NewRoute(group)
	bills.NewRoute(group)
}

func getUserData(c *fiber.Ctx) *shared.JwtUserClaims {
	token := c.Locals("user").(*jwt.Token)
	return token.Claims.(*shared.JwtUserClaims)
}
