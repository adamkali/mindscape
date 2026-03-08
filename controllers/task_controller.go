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

// @Summary Create a new Task
// @Description Create a new Task by Authorization Header
//
// @ID          CreateTask
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       CreateTaskRequest   body         repository.InsertNewTaskParams    true "CreateTaskRequest"
// @Success     200                 {object}     responses.TaskResponse
// @Failure     400                 {object}     responses.TaskResponse
// @Failure     401                 {object}     responses.TaskResponse
// @Failure     403                 {object}     responses.TaskResponse
// @Failure     500                 {object}     responses.TaskResponse
// @Router      /tasks              [post]
func (c TaskController) Create(e echo.Context) error {
	return handlers.CreateHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Update a Task
// @Description Update a Task by Authorization Header
//
// @ID          UpdateTask
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       UpdateTaskRequest   body         repository.UpdateTaskContentParams    true "UpdateTaskRequest"
// @Success     200                 {object}     responses.TaskResponse
// @Failure     400                 {object}     responses.TaskResponse
// @Failure     401                 {object}     responses.TaskResponse
// @Failure     403                 {object}     responses.TaskResponse
// @Failure     500                 {object}     responses.TaskResponse
// @Router      /tasks              [put]
func (c TaskController) Update(e echo.Context) error {
	return handlers.UpdateTaskContentHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Delete a Task
// @Description Delete a Task by Authorization Header
//
// @ID          DeleteTask
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       taskId              path         string                         true "TaskID"
// @Success     200                 {object}     responses.StringResponse
// @Failure     400                 {object}     responses.StringResponse
// @Failure     401                 {object}     responses.StringResponse
// @Failure     403                 {object}     responses.StringResponse
// @Failure     500                 {object}     responses.StringResponse
// @Router      /tasks/{taskId}     [delete]
func (c TaskController) Delete(e echo.Context) error {
	return handlers.DeleteHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Update Task Status
// @Description Update Task Status with a Status Char and optional Due Date
//
// @ID          UpdateTaskStatus 
// @Produce     json
// @Param       Authorization       header       string                 true  "auth header"     default(Bearer token)
// @Param       taskId              path         string                 true  "TaskID"
// @Param       status              query        string                 true  "Status char"     default(a)
// @Param       dueDate             query        string                 false "Due Date"
// @Success     200                 {object}     responses.TaskResponse
// @Failure     400                 {object}     responses.TaskResponse
// @Failure     401                 {object}     responses.TaskResponse
// @Failure     403                 {object}     responses.TaskResponse
// @Failure     500                 {object}     responses.TaskResponse
// @Router      /tasks/{taskId}     [put]
func (c TaskController) UpdateTaskStatus(e echo.Context) error {
	return handlers.UpdateTaskStatusHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Get Tasks By Queue Type
// @Description Get Tasks By Queue Type with a Queue Type Char 
// 
// @ID          GetTasksByQueueType
// @Produce     json 
// @Param       Authorization       header       string                 true  "auth header"     default(Bearer token)
// @Param       queueType           query        string                 true  "Queue Type Char"     default(a)
// @Success     200                 {object}     responses.TasksResponse
// @Failure     400                 {object}     responses.TasksResponse
// @Failure     401                 {object}     responses.TasksResponse
// @Failure     403                 {object}     responses.TasksResponse
// @Failure     500                 {object}     responses.TasksResponse
// @Router      /tasks/queue        [get]
func (c TaskController) GetTasksByQueueType(e echo.Context) error {
	return handlers.GetTasksByQueueTypeHandlerJsonHandler(e, *c.Registrar)
}


// @Summary Get Tasks By Task Type
// @Description Get Tasks By Task Type with a Task Type Char 
// 
// @ID          GetTasksByTaskType
// @Produce     json
// @Param       Authorization       header       string                 true  "auth header"     default(Bearer token)
// @Param       taskType            query        string                 true  "Task Type Char"     default(a)
// @Success     200                 {object}     responses.TasksResponse
// @Failure     400                 {object}     responses.TasksResponse
// @Failure     401                 {object}     responses.TasksResponse
// @Failure     403                 {object}     responses.TasksResponse
// @Failure     500                 {object}     responses.TasksResponse
// @Router      /tasks/status       [get]
func (c TaskController) GetTasksByTaskType(e echo.Context) error {
	return handlers.GetTasksByTaskTypeJsonHandler(e, *c.Registrar)
}

func (c TaskController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api/" + c.Name)
	// many
	api.GET("", c.Read, middlewares...)
	api.GET("/queue", c.GetTasksByQueueType, middlewares...)
	api.GET("/status", c.GetTasksByTaskType, middlewares...)
	// single
	api.GET("/:taskId", c.ReadByID, middlewares...)
	api.POST("", c.Create, middlewares...)
	api.PUT("", c.Update, middlewares...)
	api.PUT("/status", c.UpdateTaskStatus, middlewares...)
	api.DELETE("/:taskId", c.Delete, middlewares...)
}
