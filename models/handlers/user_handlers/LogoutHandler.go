package user_handlers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

// LogoutHandler revokes the session belonging to the refresh-token cookie and
// clears the cookie. Only this device's session is deleted — other
// browsers/devices for the same user stay logged in.
type LogoutHandler struct {
	ctx         echo.Context
	code        int
	err         error
	Config      *configuration.Configuration
	AuthService services.IAuthService
}

func NewLogoutHandler(
	ctx echo.Context,
	config *configuration.Configuration,
	AuthService services.IAuthService,
) *LogoutHandler {
	return &LogoutHandler{
		ctx:         ctx,
		code:        200,
		err:         nil,
		Config:      config,
		AuthService: AuthService,
	}
}

func (h *LogoutHandler) Handle() handlers.IHandler {
	cookie, err := h.ctx.Cookie(RefreshCookieName)
	if err == nil && cookie.Value != "" {
		if err := h.AuthService.RevokeSession(cookie.Value); err != nil {
			return handlers.Lock(h, 500, err)
		}
	}
	// Always clear the cookie, even if there was no session to revoke —
	// logout should be idempotent.
	ClearRefreshCookie(h.ctx, h.Config)
	return h
}

func (h *LogoutHandler) JSON() error {
	if h.err == nil {
		return responses.NewStringResponse().Successful(h.ctx, "logged out")
	}
	return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
}

func (h *LogoutHandler) Data() any {
	return nil
}

func (h *LogoutHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *LogoutHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *LogoutHandler) Error() error {
	return h.err
}

func (h *LogoutHandler) Code() int {
	return h.code
}
