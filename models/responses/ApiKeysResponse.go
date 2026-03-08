package responses

import (
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type ApiKeysResponse struct {
	Data    []services.ApiKeyDTO `json:"data"`
	Message string               `json:"message"`
	Success bool                 `json:"success"`
}

func NewApiKeysResponse() *ApiKeysResponse {
	return &ApiKeysResponse{
		Data:    []services.ApiKeyDTO{},
		Success: true,
		Message: "Ok",
	}
}

func (r *ApiKeysResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *ApiKeysResponse) Successful(ctx echo.Context, keys []services.ApiKeyDTO) error {
	r.Data = keys
	return ctx.JSON(200, r)
}
