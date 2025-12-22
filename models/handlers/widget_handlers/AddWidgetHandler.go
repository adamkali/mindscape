package widget_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AddWidgetHandler struct {
	ctx           echo.Context
	code          int
	err           error
	data          *repository.UserWidget
	validator     services.ValidatorService
	widgetService services.IWidgetService
	authService   services.IAuthService
}

func NewAddWidgetHandler (
	ctx echo.Context,
	validator services.ValidatorService,
	ws services.IWidgetService,
	as services.IAuthService,
) *AddWidgetHandler{
	return &AddWidgetHandler{
		ctx:           ctx,
		validator:     validator,
		widgetService: ws,
		authService:   as,
	}
}

func (ah *AddWidgetHandler) Handle() handlers.IHandler {
	jwt_token := ah.ctx.Get("user").(*jwt.Token)
	err := ah.authService.CheckToken(jwt_token.Raw)
	userId := jwt_token.Claims.(*services.CustomJwt).UserId
	if err != nil {
		return handlers.Lock(ah, 401, err)
	}

	request, err := ah.validator.ValidateAddUserWidgetRequest(ah.ctx)
	if err != nil {
		return handlers.Lock(ah, 400, err)
	}

	ah.data, err = ah.widgetService.CreateWidget(request.IntoRepositoryParams(userId))
	if err != nil {
		return handlers.Lock(ah, 500, err)
	}
	ah.code = 200
	return ah
}

func (ah *AddWidgetHandler) JSON() error {
	if ah.err == nil {
		return responses.NewUserWidgetResponse().Successful(ah.ctx, ah.data)
	} else {
		return responses.NewUserWidgetResponse().Fail(ah.ctx, ah.code, ah.err)
	}
}

func (ah *AddWidgetHandler) SetCode(code int) handlers.IHandler {
	ah.code = code
	return ah
}

func (ah *AddWidgetHandler) SetError(err error) handlers.IHandler {
	ah.err = err
	return ah
}

func (ah *AddWidgetHandler) Code() int {
	return ah.code
}

func (ah *AddWidgetHandler) Error() error {
	return ah.err
}

func (ah *AddWidgetHandler) Data() any {
	return ah.data
}

