package widget_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/schemas"
	"github.com/adamkali/mindscape/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetByIdHandler struct {
	ctx           echo.Context
	code          int
	err           error
	data          *schemas.WidgetSchema
	widgetService services.IWidgetService
}

func NewGetByIdHandler(ctx echo.Context, widgetService services.IWidgetService) *GetByIdHandler {
	return &GetByIdHandler{
		ctx:           ctx,
		widgetService: widgetService,
	}
}

func (gh *GetByIdHandler) Handle() handlers.IHandler {
	schemaId, err := uuid.Parse(gh.ctx.Param("schema_id"))
	if err != nil {
		return handlers.Lock(gh, 400, err)
	}
	gh.data, gh.err = gh.widgetService.GetWidgetSchema(schemaId)
	if gh.err != nil {
		return handlers.Lock(gh, 404, gh.err)
	}
	gh.code = 200
	return gh
}

func (gh *GetByIdHandler) JSON() error {
	if gh.err == nil {
		return responses.NewWidgetResponse().Successful(gh.ctx, *gh.data)
	} else {
		return responses.NewWidgetResponse().Fail(gh.ctx, gh.code, gh.err)
	}
}

func (gh *GetByIdHandler) SetError(err error) handlers.IHandler {
	gh.err = err
	return gh
}

func (gh *GetByIdHandler) SetCode(code int) handlers.IHandler {
	gh.code = code
	return gh
}

func (gh *GetByIdHandler) Code() int {
	return gh.code
}

func (gh *GetByIdHandler) Data() any {
	return gh.data
}

func (gh *GetByIdHandler) Error() error {
	return gh.err
}
