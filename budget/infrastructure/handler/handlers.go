package handler

import (
	"errors"
	"log"
	"your-accounts-api/budget/application"
	"your-accounts-api/budget/infrastructure/model"
	"your-accounts-api/budget/infrastructure/repository/budget"
	projectApp "your-accounts-api/project/application"
	"your-accounts-api/project/infrastructure/repository/project"
	"your-accounts-api/project/infrastructure/repository/project_log"

	// "your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/domain/jwt"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2"
	goJwt "github.com/golang-jwt/jwt/v5"
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
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	userData := getUserData(c)

	var id uint
	var err error
	if request.CloneId == nil {
		id, err = ctrl.app.Create(c.UserContext(), userData.ID, request.Name)
	} else {
		id, err = ctrl.app.Clone(c.UserContext(), userData.ID, *request.CloneId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Clone ID %d not found\n", *request.CloneId)
			return fiber.NewError(fiber.StatusNotFound, "Error creating budget. Clone ID not found")
		}
	}

	if err != nil {
		log.Printf("Error creating budget: %v\n", err)
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
		log.Printf("Error reading budgets by user: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading budgets by user")
	}

	response := []model.ReadResponse{}
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

	return c.JSON(model.NewReadByIDResponse(budget))
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
		log.Printf("Error getting param 'id': %v\n", err)
		return fiber.ErrBadRequest
	}

	err = ctrl.app.Delete(c.UserContext(), uint(id))
	if err != nil {
		log.Printf("Error deleting budget: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Budget ID not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting budget")
	}

	return c.SendStatus(fiber.StatusOK)
}

func getUserData(c *fiber.Ctx) *jwt.JwtUserClaims {
	token := c.Locals("user").(*goJwt.Token)
	return token.Claims.(*jwt.JwtUserClaims)
}

func NewRoute(router fiber.Router) {
	budgetRepo := budget.DefaultRepository()
	projectRepo := project.DefaultRepository()
	projectLogRepo := project_log.DefaultRepository()
	projectApp := projectApp.NewProjectApp(db.Tm, projectRepo, projectLogRepo)
	app := application.NewBudgetApp(db.Tm, budgetRepo, projectApp)

	controller := &controller{app}

	group := router.Group("/budget")
	group.Post("/", controller.create)
	group.Get("/", controller.read)
	group.Get("/:id<min(1)>", controller.readById)
	group.Delete("/:id<min(1)>", controller.delete)
}
