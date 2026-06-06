package user_handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamkali/mindscape/cmd/configuration"
	h "github.com/adamkali/mindscape/models/handlers/user_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=user_handlers.RefreshHandler.Handle>
// <fixtures/>

// NewRefreshTestRequest creates a refresh request carrying the refresh cookie
func NewRefreshTestRequest(cookieValue string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/users/refresh", nil)
	if cookieValue != "" {
		req.AddCookie(&http.Cookie{
			Name:  h.RefreshCookieName,
			Value: cookieValue,
		})
	}
	return req
}

// NewRefreshEchoContext creates an echo context keeping the recorder so the
// response cookies can be asserted
func NewRefreshEchoContext(r *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	return echo.New().NewContext(r, rec), rec
}

// NewRefreshHandler creates a RefreshHandler with services
func NewRefreshHandler(ctx echo.Context, authService services.IAuthService) *h.RefreshHandler {
	return h.NewRefreshHandler(ctx, &configuration.Configuration{}, authService)
}

// WithForceRefreshFailure modifies the mock to force RefreshSession failure
func WithForceRefreshFailure(authService *services.MockAuthService) *services.MockAuthService {
	authService.ShouldFailRefreshSession = true
	authService.RefreshSessionErrorMessage = "invalid or expired refresh token"
	return authService
}

// findCookie returns the named cookie from the recorded response, or nil
func findCookie(rec *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// <runners/>

// Run_RefreshHandler_ValidCookie executes RefreshHandler with a valid refresh cookie
func Run_RefreshHandler_ValidCookie(t *testing.T) {
	authService := &services.MockAuthService{}
	ctx, rec := NewRefreshEchoContext(NewRefreshTestRequest("valid-refresh-token"))

	result := NewRefreshHandler(ctx, authService).Handle()

	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.Equal(t, 1, authService.RefreshSessionCallCount)
	assert.Equal(t, "valid-refresh-token", authService.LastRefreshSessionToken)

	// rotated refresh cookie is set with the secure attributes
	cookie := findCookie(rec, h.RefreshCookieName)
	assert.NotNil(t, cookie)
	assert.NotEmpty(t, cookie.Value)
	assert.NotEqual(t, "valid-refresh-token", cookie.Value)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, h.RefreshCookiePath, cookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
}

// Run_RefreshHandler_MissingCookie executes RefreshHandler without a cookie
func Run_RefreshHandler_MissingCookie(t *testing.T) {
	authService := &services.MockAuthService{}
	ctx, _ := NewRefreshEchoContext(NewRefreshTestRequest(""))

	result := NewRefreshHandler(ctx, authService).Handle()

	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
	assert.Equal(t, 0, authService.RefreshSessionCallCount)
}

// Run_RefreshHandler_RotatedOrRevoked executes RefreshHandler with a stale token
func Run_RefreshHandler_RotatedOrRevoked(t *testing.T) {
	authService := WithForceRefreshFailure(&services.MockAuthService{})
	ctx, rec := NewRefreshEchoContext(NewRefreshTestRequest("stale-refresh-token"))

	result := NewRefreshHandler(ctx, authService).Handle()

	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())

	// the stale cookie is cleared so the client falls back to login
	cookie := findCookie(rec, h.RefreshCookieName)
	assert.NotNil(t, cookie)
	assert.Empty(t, cookie.Value)
	assert.Equal(t, -1, cookie.MaxAge)
}

// <tests>
// <map/>

// RefreshHandlerTestMap defines all RefreshHandler test cases
var RefreshHandlerTestMap = map[string]func(*testing.T){
	"ValidCookie":      Run_RefreshHandler_ValidCookie,
	"MissingCookie":    Run_RefreshHandler_MissingCookie,
	"RotatedOrRevoked": Run_RefreshHandler_RotatedOrRevoked,
}

// <hook/>

// Test_RefreshHandler tests all RefreshHandler scenarios
func Test_RefreshHandler(t *testing.T) {
	for name, testFunc := range RefreshHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>
