package responses

import (
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type ApiKeyResponse struct {
	Data    *services.ApiKeyDTO `json:"data"`
	Message string              `json:"message"`
	Success bool                `json:"success"`
}

func NewApiKeyResponse() *ApiKeyResponse {
	return &ApiKeyResponse{
		Data:    nil,
		Success: true,
		Message: "Ok",
	}
}

func (r *ApiKeyResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *ApiKeyResponse) Successful(ctx echo.Context, key *services.ApiKeyDTO) error {
	r.Data = key
	return ctx.JSON(200, r)
}
