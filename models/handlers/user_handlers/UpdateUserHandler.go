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

type UpdateUserHandler struct {
	User   *repository.User
	Token  *string
	ctx    echo.Context
	err    error
	code   int
	vs     services.ValidatorService
	us     services.IUserService
	as     services.IAuthService
}

func NewUpdateUserHandler(ctx echo.Context, validator services.ValidatorService, userService services.IUserService, authService services.IAuthService) *UpdateUserHandler {
	return &UpdateUserHandler{
		ctx:    ctx,
		err:    nil,
		code:   200,
		vs:     validator,
		us:     userService,
		as:     authService,
	}
}


func (h *UpdateUserHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userId := claims.UserId
	err := h.as.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	request, err := h.vs.ValidateUpdateUserCredentialRequest(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	if request.ID != userId {
		return handlers.Lock(h, 403, fmt.Errorf("user ID mismatch: cannot update another user's credentials"))
	}
	h.User, err = h.us.UpdateUserCredentials(request)
	if err != nil {
		return handlers.Lock(h, 500, err)
	}
	h.Token, err = h.as.Update(*h.User)
	if err != nil {
		return handlers.Lock(h, 500, err)
	}

	return h
}

func (h *UpdateUserHandler) JSON() error {
	if h.Token == nil {
		h.Token = new(string)
	}
	if h.err != nil {
		return responses.NewLoginResponse().Fail(h.ctx, h.code, h.err)
	} else {
		return responses.NewLoginResponse().Successful(h.ctx, h.User, *h.Token)
	}

}

func (h *UpdateUserHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *UpdateUserHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *UpdateUserHandler) Code() int {
	return h.code
}

func (h *UpdateUserHandler) Data() any {
	return struct {
		User  *repository.User
		Token *string
	}{
		User:  h.User,
		Token: h.Token,
	}
}

func (h *UpdateUserHandler) Error() error {
	return h.err
}
