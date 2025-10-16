package user_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


type BackgroundHandler struct {
	ctx              echo.Context
	code             int
	err              error
	url              string
	Query            string
	ValidatorService services.ValidatorService
	AuthService      services.IAuthService
	MinioService     services.IMinioService
}

func NewBackgroundHandler(
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

func (h *BackgroundHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	err := h.AuthService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}

	h.Query = h.ctx.QueryParam("url")
	if h.Query == "" {
		h.url, h.err = h.MinioService.GetDefault()
		if h.err != nil {
			return handlers.Lock(h,500, h.err)
		}
	}
	h.url, h.err = h.MinioService.GetPresigned(userID, h.Query)
	if h.err != nil {
	    h.url, h.err = h.MinioService.GetDefaultChoice(h.Query)
		if h.err != nil {
		    return handlers.Lock(h,500, h.err)
		}
	}
	return h
}

func (h *BackgroundHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewStringResponse().Successful(h.ctx, h.url)
	}
}

func (h *BackgroundHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *BackgroundHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *BackgroundHandler) Code() int {
	return h.code
}

func (h *BackgroundHandler) Data() any {
	return h.url
}

func (h *BackgroundHandler) Error() error {
	return h.err
}
