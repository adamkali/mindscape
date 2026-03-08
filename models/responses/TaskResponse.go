package responses

import (
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type TaskResponse struct {
	Data *services.TaskDTO `json:"data"`
	Message string `json:"message"`
	Success bool `json:"success"`
}

func NewTaskResponse() *TaskResponse {
	return &TaskResponse{
		Data: nil,
		Success: true,
		Message: "Ok",
	}
}

func (t *TaskResponse) Fail(ctx echo.Context, code int, err error) error {
	t.Success = false
	t.Message = err.Error()
	return ctx.JSON(code, t)
}

func (t *TaskResponse) Successful(ctx echo.Context, task *services.TaskDTO) error {
	t.Data = task
	return ctx.JSON(200, t)
}
