package handler

import (
	"api-your-accounts/project/application"
	"api-your-accounts/project/domain"
	"api-your-accounts/project/infrastructure/model"
	"api-your-accounts/project/infrastructure/repository/project"
	"api-your-accounts/project/infrastructure/repository/project_log"
	"api-your-accounts/shared/infrastructure/db"
	"api-your-accounts/shared/infrastructure/db/persistent"
	"api-your-accounts/shared/infrastructure/validation"
	"errors"
	"log"

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
//	@Failure		409					{string}	string
//	@Failure		422					{string}	string
//	@Failure		500					{string}	string
//	@Router			/api/v1/project/	[post]
func (ctrl *controller) create(c *fiber.Ctx) error {
	request := new(model.CreateRequest)
	if ok := validation.Validate(c, request); !ok {
		return nil
	}

	project := &domain.Project{
		Name:   request.Name,
		UserId: request.UserId,
		Type:   request.Type,
	}

	var result *domain.Project
	var err error
	if request.CloneId == nil {
		result, err = ctrl.app.Create(c.UserContext(), project)
	} else {
		result, err = ctrl.app.Clone(c.UserContext(), project, *request.CloneId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Clone ID %d not found\n", *request.CloneId)
			return fiber.NewError(fiber.StatusNotFound, "Error creating project. Clone ID not found")
		}
	}

	if err != nil {
		log.Println("Error creating project:", err)
		if errors.Is(err, application.ErrProjectAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusInternalServerError, "Error creating project")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(model.CreateResponse{
		ID: result.ID,
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
		log.Println("Error getting param 'user':", err)
		return fiber.ErrBadRequest
	}

	projects, err := ctrl.app.FindByUser(c.UserContext(), uint(userId))
	if err != nil {
		log.Println("Error reading projects by user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading projects by user")
	}

	response := []model.ReadResponse{}
	for _, project := range projects {
		response = append(response, model.ReadResponse{
			ID:   project.ID,
			Name: project.Name,
			Type: project.Type,
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
		log.Println("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	logs, err := ctrl.app.FindLogsByProject(c.UserContext(), uint(projectId))
	if err != nil {
		log.Println("Error reading logs by project:", err)
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
		log.Println("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	err = ctrl.app.Delete(c.UserContext(), uint(id))
	if err != nil {
		log.Println("Error deleting project:", err)
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
	controller := &controller{
		app: application.NewProjectApp(tm, projectRepo, projectLogRepo),
	}

	group := router.Group("/project")
	group.Post("/", controller.create)
	group.Get("/:user<min(1)>", controller.readByUser)
	group.Get("/logs/:id<min(1)>", controller.readLogs)
	group.Delete("/:id<min(1)>", controller.delete)
}
