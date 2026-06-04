package user_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

// RefreshHandler exchanges a valid refresh-token cookie for a new access JWT
// and a rotated refresh cookie. It is unauthenticated by design: the access
// JWT is expected to be expired when this is called, so the refresh cookie is
// the credential.
type RefreshHandler struct {
	user        *repository.User
	token       *string
	ctx         echo.Context
	code        int
	err         error
	Config      *configuration.Configuration
	AuthService services.IAuthService
}

func NewRefreshHandler(
	ctx echo.Context,
	config *configuration.Configuration,
	AuthService services.IAuthService,
) *RefreshHandler {
	return &RefreshHandler{
		ctx:         ctx,
		code:        200,
		user:        nil,
		token:       nil,
		err:         nil,
		Config:      config,
		AuthService: AuthService,
	}
}

func (h *RefreshHandler) Handle() handlers.IHandler {
	cookie, err := h.ctx.Cookie(RefreshCookieName)
	if err != nil || cookie.Value == "" {
		return handlers.Lock(h, 401, fmt.Errorf("missing refresh token"))
	}
	access, newRefreshRaw, user, err := h.AuthService.RefreshSession(cookie.Value)
	if err != nil {
		// The token was already rotated, revoked, or expired — clear the
		// stale cookie so the client falls back to a clean login.
		ClearRefreshCookie(h.ctx, h.Config)
		return handlers.Lock(h, 401, err)
	}
	h.user = user
	h.token = &access
	SetRefreshCookie(h.ctx, h.Config, newRefreshRaw)
	return h
}

func (h *RefreshHandler) JSON() error {
	var jwt string
	if h.token != nil {
		jwt = *h.token
	}
	if h.err == nil {
		return responses.NewLoginResponse().Successful(h.ctx, h.user, jwt)
	}
	return responses.NewLoginResponse().Fail(h.ctx, h.code, h.err)
}

func (h *RefreshHandler) Data() any {
	return struct {
		User  *repository.User
		Token *string
	}{
		User:  h.user,
		Token: h.token,
	}
}

func (h *RefreshHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *RefreshHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *RefreshHandler) Error() error {
	return h.err
}

func (h *RefreshHandler) Code() int {
	return h.code
}
