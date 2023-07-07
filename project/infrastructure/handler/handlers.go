package handler

import (
	"errors"
	"log"
	"your-accounts-api/budget/infrastructure/repository/budget"
	"your-accounts-api/project/application"
	"your-accounts-api/project/infrastructure/model"
	"your-accounts-api/project/infrastructure/repository/project"
	"your-accounts-api/project/infrastructure/repository/project_log"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/db/persistent"
	"your-accounts-api/shared/infrastructure/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
//	@Failure		404					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/project/	[post]
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	var projectId uint
	var err error
	if request.CloneId == nil {
		create := application.CreateData{
			Name:   request.Name,
			UserId: request.UserId,
			Type:   request.Type,
		}
		projectId, err = ctrl.app.Create(c.UserContext(), create)
	} else {
		projectId, err = ctrl.app.Clone(c.UserContext(), *request.CloneId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Clone ID %d not found\n", *request.CloneId)
			return fiber.NewError(fiber.StatusNotFound, "Error creating project. Clone ID not found")
		}
	}

	if err != nil {
		log.Printf("Error creating project: %v\n", err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error creating project")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreateResponse{
		ID: projectId,
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
	userId, err := c.ParamsInt("user")
	if err != nil {
		log.Printf("Error getting param 'user': %v\n", err)
		return fiber.ErrBadRequest
	}

	projects, err := ctrl.app.FindByUser(c.UserContext(), uint(userId))
	if err != nil {
		log.Printf("Error reading projects by user: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading projects by user")
	}

	response := []model.ReadResponse{}
	for _, record := range projects {
		response = append(response, model.ReadResponse{
			ID:   record.ID,
			Name: record.Name,
			Type: record.Type,
			Data: record.Data,
		})
	}

	return c.JSON(response)
}

// ProjectReadTransactionsHandler godoc
//
//	@Summary		Read logs by project
//	@Description	read logs associated to a project
//	@Tags			project
//	@Produce		json
//	@Param			Authorization				header		string	true	"Access token"
//	@Param			id							path		uint	true	"Project ID"
//	@Success		200							{array}		model.ReadLogsResponse
//	@Failure		400							{string}	string
//	@Failure		401							{string}	string
//	@Failure		404							{string}	string
//	@Failure		500							{string}	string
//	@Router			/api/v1/project/logs/{id}	[get]
func (ctrl *controller) readLogs(c *fiber.Ctx) error {
	projectId, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Error getting param 'id': %v\n", err)
		return fiber.ErrBadRequest
	}

	logs, err := ctrl.app.FindLogsByProject(c.UserContext(), uint(projectId))
	if err != nil {
		log.Printf("Error reading logs by project: %v\n", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading logs by project")
	}

	response := []model.ReadLogsResponse{}
	for _, log := range logs {
		response = append(response, model.ReadLogsResponse{
			ID:          log.ID,
			Description: log.Description,
			CreatedAt:   log.CreatedAt,
		})
	}

	return c.JSON(response)
}

// ProjectDeleteHandler godoc
//
//	@Summary		Delete project
//	@Description	Delete an project by ID
//	@Tags			project
//	@Produce		json
//	@Param			Authorization			header		string	true	"Access token"
//	@Param			id						path		uint	true	"Project ID"
//	@Success		200						{string}	string
//	@Failure		400						{string}	string
//	@Failure		401						{string}	string
//	@Failure		404						{string}	string
//	@Failure		500						{string}	string
//	@Router			/api/v1/project/{id}	[delete]
func (ctrl *controller) delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Error getting param 'id': %v\n", err)
		return fiber.ErrBadRequest
	}

	err = ctrl.app.Delete(c.UserContext(), uint(id))
	if err != nil {
		log.Printf("Error deleting project: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Project ID not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting project")
	}

	return c.SendStatus(fiber.StatusOK)
}

func NewRoute(router fiber.Router) {
	tm := persistent.NewTransactionManager(db.DB)
	projectRepo := project.NewRepository(db.DB)
	projectLogRepo := project_log.NewRepository(db.DB)
	budgetRepo := budget.NewRepository(db.DB)
	controller := &controller{
		app: application.NewProjectApp(tm, projectRepo, projectLogRepo, budgetRepo),
	}

	group := router.Group("/project")
	group.Post("/", controller.create)
	group.Get("/:user<min(1)>", controller.readByUser)
	group.Get("/logs/:id<min(1)>", controller.readLogs)
	group.Delete("/:id<min(1)>", controller.delete)
}
