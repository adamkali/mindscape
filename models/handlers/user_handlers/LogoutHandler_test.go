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

// <method var=user_handlers.LogoutHandler.Handle>
// <fixtures/>

// NewLogoutTestRequest creates a logout request carrying the refresh cookie
func NewLogoutTestRequest(cookieValue string) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/api/users/refresh", nil)
	if cookieValue != "" {
		req.AddCookie(&http.Cookie{
			Name:  h.RefreshCookieName,
			Value: cookieValue,
		})
	}
	return req
}

// NewLogoutEchoContext creates an echo context keeping the recorder
func NewLogoutEchoContext(r *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	return echo.New().NewContext(r, rec), rec
}

// NewLogoutHandler creates a LogoutHandler with services
func NewLogoutHandler(ctx echo.Context, authService services.IAuthService) *h.LogoutHandler {
	return h.NewLogoutHandler(ctx, &configuration.Configuration{}, authService)
}

// WithForceRevokeFailure modifies the mock to force RevokeSession failure
func WithForceRevokeFailure(authService *services.MockAuthService) *services.MockAuthService {
	authService.ShouldFailRevokeSession = true
	return authService
}

// <runners/>

// Run_LogoutHandler_ValidCookie executes LogoutHandler with a refresh cookie
func Run_LogoutHandler_ValidCookie(t *testing.T) {
	authService := &services.MockAuthService{}
	ctx, rec := NewLogoutEchoContext(NewLogoutTestRequest("session-refresh-token"))

	result := NewLogoutHandler(ctx, authService).Handle()

	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.Equal(t, 1, authService.RevokeSessionCallCount)
	assert.Equal(t, "session-refresh-token", authService.LastRevokeSessionToken)

	// cookie is cleared
	cookie := findCookie(rec, h.RefreshCookieName)
	assert.NotNil(t, cookie)
	assert.Empty(t, cookie.Value)
	assert.Equal(t, -1, cookie.MaxAge)
}

// Run_LogoutHandler_MissingCookie executes LogoutHandler without a cookie (idempotent)
func Run_LogoutHandler_MissingCookie(t *testing.T) {
	authService := &services.MockAuthService{}
	ctx, rec := NewLogoutEchoContext(NewLogoutTestRequest(""))

	result := NewLogoutHandler(ctx, authService).Handle()

	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.Equal(t, 0, authService.RevokeSessionCallCount)

	// cookie is still cleared — logout is idempotent
	cookie := findCookie(rec, h.RefreshCookieName)
	assert.NotNil(t, cookie)
	assert.Empty(t, cookie.Value)
}

// Run_LogoutHandler_RevokeFailure executes LogoutHandler with RevokeSession failure
func Run_LogoutHandler_RevokeFailure(t *testing.T) {
	authService := WithForceRevokeFailure(&services.MockAuthService{})
	ctx, _ := NewLogoutEchoContext(NewLogoutTestRequest("session-refresh-token"))

	result := NewLogoutHandler(ctx, authService).Handle()

	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusInternalServerError, result.Code())
}

// <tests>
// <map/>

// LogoutHandlerTestMap defines all LogoutHandler test cases
var LogoutHandlerTestMap = map[string]func(*testing.T){
	"ValidCookie":   Run_LogoutHandler_ValidCookie,
	"MissingCookie": Run_LogoutHandler_MissingCookie,
	"RevokeFailure": Run_LogoutHandler_RevokeFailure,
}

// <hook/>

// Test_LogoutHandler tests all LogoutHandler scenarios
func Test_LogoutHandler(t *testing.T) {
	for name, testFunc := range LogoutHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>
