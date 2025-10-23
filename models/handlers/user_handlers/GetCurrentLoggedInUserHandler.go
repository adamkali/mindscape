package user_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type GetCurrentLoggedInUserHandler struct {
	ctx  echo.Context
	err  error
	code int
	Vs   services.ValidatorService
	Us   services.IUserService
	As   services.IAuthService
	data *repository.User
}

func NewGetCurrentLoggedInUserHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	UserService services.IUserService,
	AuthService services.IAuthService,
) *GetCurrentLoggedInUserHandler {
	return &GetCurrentLoggedInUserHandler{
		ctx:  ctx,
		err:  nil,
		code: 200,
		Vs:   ValidatorService,
		Us:   UserService,
		As:   AuthService,
	}
}

func (h *GetCurrentLoggedInUserHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	err := h.As.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	h.data, err = h.Us.Get(claims.UserId)
	if err != nil {
		return handlers.Lock(h, 404, err)
	}
	return h
}

func (h *GetCurrentLoggedInUserHandler) JSON() error {
	if h.err != nil {
		return responses.NewUserResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewUserResponse().Successful(h.ctx, h.data)
}

func (h *GetCurrentLoggedInUserHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *GetCurrentLoggedInUserHandler) Code() int {
	return h.code
}

func (h *GetCurrentLoggedInUserHandler) Data() any {
	return h.data
}

func (h *GetCurrentLoggedInUserHandler) Error() error {
	return h.err
}
func (h *GetCurrentLoggedInUserHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}
