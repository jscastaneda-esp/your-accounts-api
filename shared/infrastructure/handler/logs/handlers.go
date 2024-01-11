package logs

import (
	"fmt"
	"your-accounts-api/shared/application"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/model"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

type controller struct {
	app application.ILogApp
}

// ReadLogsHandler godoc
//
//	@Summary		Read logs by resource and code
//	@Description	read logs associated to a resource and code
//	@Tags			log
//	@Produce		json
//	@Param			Authorization					header		string			true	"Access token"
//	@Param			id								path		uint			true	"Resource ID"
//	@Param			code							path		domain.CodeLog	true	"Code"
//	@Success		200								{array}		model.ReadLogsResponse
//	@Failure		400								{string}	string
//	@Failure		401								{string}	string
//	@Failure		404								{string}	string
//	@Failure		500								{string}	string
//	@Router			/api/v1/log/{id}/code/{code}	[get]
func (ctrl *controller) readLogs(c *fiber.Ctx) error {
	resourceId, err := c.ParamsInt("id")
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return fiber.ErrBadRequest
	}

	code := c.Params("code")
	if code == "" {
		log.Error("Error getting param 'code'")
		return fiber.ErrBadRequest
	}

	logs, err := ctrl.app.FindByProject(c.UserContext(), domain.CodeLog(code), uint(resourceId))
	if err != nil {
		log.Error("Error reading logs by resource and code:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading logs by resource and code")
	}

	response := make([]*model.ReadLogsResponse, 0)
	for _, logItem := range logs {
		response = append(response, model.NewReadLogsResponse(logItem))
	}

	return c.JSON(response)
}

func NewRoute(router fiber.Router) {
	controller := &controller{injection.LogApp}

	group := router.Group("/log")
	group.Get(fmt.Sprintf("/:id<min(1)>/code/:code<regex(^(%s|%s)$)>", domain.Budget, domain.BudgetBill), controller.readLogs)
}
