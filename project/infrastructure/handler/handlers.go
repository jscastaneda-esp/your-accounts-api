package handler

import (
	"api-your-accounts/project/application"
	"api-your-accounts/project/domain"
	"api-your-accounts/project/infrastructure/model"
	"api-your-accounts/shared/infrastructure/validation"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type controller struct {
	app application.IProjectApp
}

// ProjectCreateHandler godoc
//
//	@Summary		Create project
//	@Description	create a new project
//	@Tags			project
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string				true	"Access token"
//	@Param			request				body		model.CreateRequest	true	"Project data"
//	@Success		201					{object}	model.CreateResponse
//	@Failure		400					{string}	string
//	@Failure		401					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/project/	[post]
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreateResponse{
		ID: uint(1),
	})
}

// ProjectReadByUserHandler godoc
//
//	@Summary		Read projects by user
//	@Description	read projects associated to an user
//	@Tags			project
//	@Produce		json
//	@Param			Authorization			header		string	true	"Access token"
//	@Param			user					path		uint	true	"User ID"
//	@Success		200						{array}		model.ReadResponse
//	@Failure		400						{string}	string
//	@Failure		401						{string}	string
//	@Failure		404						{string}	string
//	@Failure		500						{string}	string
//	@Router			/api/v1/project/{user}	[get]
func (ctrl *controller) readByUser(c *fiber.Ctx) error {
	log.Println(c.Params("user"))

	return c.JSON([]model.ReadResponse{
		{
			ID:   uint(1),
			Name: "Test",
			Type: domain.Budget,
		},
	})
}

// ProjectReadTransactionsHandler godoc
//
//	@Summary		Read transactions by project
//	@Description	read transactions associated to a project
//	@Tags			project
//	@Produce		json
//	@Param			Authorization						header		string	true	"Access token"
//	@Param			id									path		uint	true	"Project ID"
//	@Success		200									{array}		model.ReadTransactionResponse
//	@Failure		400									{string}	string
//	@Failure		401									{string}	string
//	@Failure		404									{string}	string
//	@Failure		500									{string}	string
//	@Router			/api/v1/project/transactions/{id}	[get]
func (ctrl *controller) readTransactions(c *fiber.Ctx) error {
	log.Println(c.Params("id"))

	return c.JSON([]model.ReadTransactionResponse{
		{
			ID:          1,
			Description: "Test",
			CreatedAt:   time.Now(),
		},
	})
}

// ProjectDeleteHandler godoc
//
//	@Summary		Delete project
//	@Description	Delete an project by ID
//	@Tags			project
//	@Produce		json
//	@Param			Authorization	header	string	true	"Access token"
//	@Param			id				path	uint	true	"Project ID"
//	@Success		200
//	@Failure		400						{string}	string
//	@Failure		401						{string}	string
//	@Failure		404						{string}	string
//	@Failure		500						{string}	string
//	@Router			/api/v1/project/{id}	[delete]
func (ctrl *controller) delete(c *fiber.Ctx) error {
	log.Println(c.Params("id"))

	return nil
}

func NewRoute(router fiber.Router) {
	controller := &controller{
		app: application.NewProjectApp(),
	}

	group := router.Group("/project")
	group.Post("/", controller.create)
	group.Get("/:user<min(1)>", controller.readByUser)
	group.Get("/transactions/:id<min(1)>", controller.readTransactions)
	group.Delete("/:id<min(1)>", controller.delete)
}
