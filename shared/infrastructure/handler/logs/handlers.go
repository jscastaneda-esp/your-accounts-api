package logs

import (
	"log"
	"your-accounts-api/shared/application"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/model"

	"github.com/gofiber/fiber/v2"
)

type controller struct {
	app application.ILogApp
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

	response := make([]*model.ReadLogsResponse, 0)
	for _, log := range logs {
		response = append(response, model.NewReadLogsResponse(log))
	}

	return c.JSON(response)
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.LogApp}

	group := router.Group("/project")
	group.Get("/logs/:id<min(1)>", controller.readLogs)
}
