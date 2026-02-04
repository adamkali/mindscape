package responses

import (
	"github.com/labstack/echo/v4"
)

type CoolifyActionData struct {
	Action       string `json:"action"`
	AppUUID      string `json:"app_uuid"`
	DeploymentId string `json:"deployment_id,omitempty"`
}

type CoolifyActionResponse struct {
	Data    *CoolifyActionData `json:"data"`
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
} // @name CoolifyActionResponse

func NewCoolifyActionResponse() *CoolifyActionResponse {
	return &CoolifyActionResponse{
		Data:    nil,
		Success: false,
		Message: "",
	}
}

func (r *CoolifyActionResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *CoolifyActionResponse) Successful(ctx echo.Context, data *CoolifyActionData) error {
	r.Success = true
	r.Data = data
	return ctx.JSON(200, r)
}
