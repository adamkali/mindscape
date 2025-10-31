package user_handlers

import (
	"fmt"

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
		fmt.Printf("[ERROR] BackgroundHandler.Handle{ userID: %v\n --> %s }\n", userID)
		return handlers.Lock(h, 401, err)
	}
	q := h.ctx.QueryParam("background")
	fmt.Printf("[DEBUG] BackgroundHandler.Handle{ background: %s }\n", h.ctx.QueryParam("background"))
	fmt.Printf("[INFO] BackgroundHandler.Handle{ userID: %v, background: %s }\n", userID, q)
	if q == "" {
		h.url, h.err = h.MinioService.GetDefault()
		if h.err != nil {
			fmt.Printf("[ERROR] BackgroundHandler.MinioService.GetDefault{\nuserID: %v\n --> %s }\n", userID, h.err)
			return handlers.Lock(h,500, h.err)
		}
	}
	fmt.Printf("[INFO] BackgroundHandler.MinioService.GetPresigned{ userID: %v, background: %s }\n", userID, q)
	h.url, h.err = h.MinioService.GetPresigned(userID, "background", q)
	if h.err != nil {
		fmt.Printf("[WARNING] BackgroundHandler.MinioService.GetPresigned{\nuserID: %v,\nbackground: %s,\n} --> %s }\n", userID, q, h.err)
		fmt.Printf("[INFO] Defaulting to BackgroundHandler.MinioService.GetDefaultChoice \n", userID, q)
	    h.url, h.err = h.MinioService.GetDefaultChoice(q)
		if h.err != nil {
			fmt.Printf("[ERROR] BackgroundHandler.MinioService.GetDefaultChoice{\nuserID: %v,\nbackground: %s\n --> %s }\n", userID, q, h.err)
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
