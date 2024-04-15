package logs

import (
	"net/http"
	"strconv"
	"your-accounts-api/shared/application"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/injection"
	"your-accounts-api/shared/infrastructure/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
func (ctrl *controller) readLogs(c echo.Context) error {
	resourceId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error("Error getting param 'id':", err)
		return echo.ErrBadRequest
	}

	code := c.Param("code")
	if code == "" {
		log.Error("Error getting param 'code'")
		return echo.ErrBadRequest
	}

	logs, err := ctrl.app.FindByProject(c.Request().Context(), domain.CodeLog(code), uint(resourceId))
	if err != nil {
		log.Error("Error reading logs by resource and code:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading logs by resource and code")
	}

	response := make([]*model.ReadLogsResponse, 0)
	for _, logItem := range logs {
		response = append(response, model.NewReadLogsResponse(logItem))
	}

	return c.JSON(http.StatusOK, response)
}

func NewRoute(api *echo.Group) {
	controller := &controller{injection.LogApp}

	group := api.Group("/log")
	group.GET("/:id/code/:code", controller.readLogs)
}
