package controllers

import (
	handlers "github.com/adamkali/mindscape/models/handlers/task_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type TaskController struct {
	Name      string
	Registrar *services.Registrar
}

func (c TaskController) ControllerName() string {
	return c.Name
}

func BuildTaskController(registrar *services.Registrar) TaskController {
	return TaskController{
		Name:      "tasks",
		Registrar: registrar,
	}
}

// @Summary Read all tasks
// @Description Read all tasks that are available to the user
//
// @ID          ReadTasks
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Success     200                 {object}     responses.TasksResponse
// @Failure     401                 {object}     responses.TasksResponse
// @Failure     500                 {object}     responses.TasksResponse
// @Router      /tasks              [get]
func (c TaskController) Read(e echo.Context) error {
	return handlers.ReadHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Read by a TaskID
// @Description Read by a TaskID
//
// @ID          ReadTask
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       taskId       		path       string                         true "TaskID"
// @Success     200                 {object}     responses.TaskResponse
// @Failure     400                 {object}     responses.TaskResponse
// @Failure     401                 {object}     responses.TaskResponse
// @Failure     403                 {object}     responses.TaskResponse
// @Failure     500                 {object}     responses.TaskResponse
// @Router      /tasks/{taskId}     [get]
func (c TaskController) ReadByID(e echo.Context) error {
	return handlers.ReadByIDHandlerJsonHandler(e, *c.Registrar)
}

func (c TaskController) Attatch(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	api := e.Group("/api" + c.Name)
	api.GET("", c.Read, authMiddleware)
}
