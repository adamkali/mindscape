package user_handlers

import (
	"net/http"
	"time"

	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/labstack/echo/v4"
)

const (
	// RefreshCookieName is the cookie holding the opaque refresh token.
	RefreshCookieName = "refresh_token"
	// RefreshCookiePath scopes the cookie so the browser only attaches it to
	// the refresh/logout endpoints, never to regular API calls.
	RefreshCookiePath = "/api/users/refresh"
)

// SetRefreshCookie attaches the refresh token to the response as an httpOnly
// cookie. JavaScript can never read it; the browser only sends it to
// RefreshCookiePath.
func SetRefreshCookie(ctx echo.Context, config *configuration.Configuration, refreshRaw string) {
	ctx.SetCookie(&http.Cookie{
		Name:     RefreshCookieName,
		Value:    refreshRaw,
		Path:     RefreshCookiePath,
		HttpOnly: true,
		Secure:   config.Server.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(config.RefreshTokenDuration()),
	})
}

// ClearRefreshCookie expires the refresh cookie (logout / failed refresh).
func ClearRefreshCookie(ctx echo.Context, config *configuration.Configuration) {
	ctx.SetCookie(&http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     RefreshCookiePath,
		HttpOnly: true,
		Secure:   config.Server.Cookie.Secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
