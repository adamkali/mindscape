package responses

import (
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type TasksResponse struct {
	Data    []services.TaskDTO `json:"data"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

func NewTasksResponse() *TasksResponse {
	return &TasksResponse{
		Data:    []services.TaskDTO {},
		Success: true,
		Message: "Ok",
	}
}

func (r *TasksResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *TasksResponse) Successful(ctx echo.Context, tasks []services.TaskDTO) error {
	r.Success = true
	r.Data = tasks
	return ctx.JSON(200, r)
}
