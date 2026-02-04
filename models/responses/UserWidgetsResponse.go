package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type UserWidgetsResponse struct {
	Data    []UserWidgetData `json:"data"`
	Success bool            `json:"success"`
	Message string          `json:"message"`
}	


func NewUserWidgetsResponse() *UserWidgetsResponse {
	return &UserWidgetsResponse{Success: false, Message: ""}
}

func (u *UserWidgetsResponse) Successful(
	ctx echo.Context,
	data []repository.UserWidget,
) error {
	u.Data = make([]UserWidgetData, len(data))
	for i, val := range data {
		u.Data[i] = UserWidgetData{}
		u.Data[i].UserWidgetFromData(&val)
	}
	u.Success = true
	u.Message = "OK"
	return ctx.JSON(200, u)
}

func (u *UserWidgetsResponse) Fail(
	ctx echo.Context,
	code int,
	err error,
) error  {
	u.Message = err.Error()
	return ctx.JSON(code, u)
}
