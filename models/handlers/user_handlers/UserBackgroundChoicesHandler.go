package user_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


type UserBackgroundChoicesHandler struct {
	ctx              echo.Context
	code             int
	err              error
	urls             []string
	ValidatorService services.ValidatorService
	AuthService      services.IAuthService
	MinioService     services.IMinioService
}

func NewUserBackgroundChoicesHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	AuthService services.IAuthService,
	MinioService services.IMinioService,
) *BackgroundHandler {
	return &BackgroundHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		AuthService:      AuthService,
		MinioService:     MinioService,
	}
 }
func (h *UserBackgroundChoicesHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	err := h.AuthService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}

	h.urls, h.err = h.MinioService.GetUserBackgroundChoices(userID)
	if h.err != nil {
		return handlers.Lock(h, 500, h.err)
	}

	return h
}

func (h *UserBackgroundChoicesHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringsResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewStringsResponse().Successful(h.ctx, h.urls)
	}
}

func (h *UserBackgroundChoicesHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *UserBackgroundChoicesHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *UserBackgroundChoicesHandler) Code() int {
	return h.code
}

func (h *UserBackgroundChoicesHandler) Data() any {
	return h.urls
}

func (h *UserBackgroundChoicesHandler) Error() error {
	return h.err
}
