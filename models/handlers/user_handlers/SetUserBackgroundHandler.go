package user_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type SetUserBackgroundHandler struct {
	ctx          echo.Context
	code         int
	err          error
	url          string
	UserService  services.IUserService
	AuthService  services.IAuthService
	MinioService services.IMinioService
}

func NewSetUserBackgroundHandler(
	ctx echo.Context,
	UserService services.IUserService,
	AuthService services.IAuthService,
	MinioService services.IMinioService,
) *SetUserBackgroundHandler {
	return &SetUserBackgroundHandler{
		ctx:          ctx,
		code:         200,
		UserService:  UserService,
		AuthService:  AuthService,
		MinioService: MinioService,
	}
}

func (h *SetUserBackgroundHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	err := h.AuthService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	fmt.Printf("[INFO] SetUserBackgroundHandler.Handle{ userID: %v }\n", userID)
	backgroundFileName := h.ctx.QueryParam("background")
	fmt.Printf("[INFO] SetUserBackgroundHandler.Handle{ backgroundFileName: %v }\n", backgroundFileName)
	if backgroundFileName == "" {
		h.url, h.err = h.MinioService.GetDefault()
		if h.err != nil {
			return handlers.Lock(h, 404, h.err)
		}
	} else {
		h.url, h.err = h.MinioService.GetPresigned(userID, "background", backgroundFileName)
		if h.err != nil {
			h.url, h.err = h.MinioService.GetDefaultChoice(backgroundFileName)
			if h.err != nil {
				return handlers.Lock(h, 404, h.err)
			}
		}
	}
	_, err = h.
		UserService.
		UpdateUserBackgroundImage(&repository.UpdateUserBacgroundParams{
			ID:         userID,
			Background: &backgroundFileName,
		})
	if err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *SetUserBackgroundHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewStringResponse().Successful(h.ctx, h.url)
	}
}

func (h *SetUserBackgroundHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *SetUserBackgroundHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *SetUserBackgroundHandler) Code() int {
	return h.code
}

func (h *SetUserBackgroundHandler) Data() any {
	return h.url
}

func (h *SetUserBackgroundHandler) Error() error {
	return h.err
}
