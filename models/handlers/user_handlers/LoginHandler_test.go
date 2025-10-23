package user_handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adamkali/mindscape/db/repository"
	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/user_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=user_handlers.LoginHandler.Handle>
// <fixtures/>

// LoginHandlerRequestParams provides base parameters for LoginHandler tests
func LoginHandlerRequestParams() map[string]any {
	return map[string]any{
		"email":    "handler@example.com",
		"password": "passwordABC123!",
	}
}

// LoginHandlerRequestByUsername provides base parameters for username login
func LoginHandlerRequestByUsername() map[string]any {
	return map[string]any{
		"username": "handleruser",
		"password": "passwordABC123!",
	}
}

// WithWrongPassword modifies request with incorrect password (but still valid format)
func WithWrongPassword(req map[string]any) map[string]any {
	req["password"] = "WrongPassword123!"
	return req
}

// WithNonExistentEmail modifies request with non-existent email
func WithNonExistentEmail(req map[string]any) map[string]any {
	req["email"] = "nonexistent@example.com"
	delete(req, "username")
	return req
}

// WithNonExistentUsername modifies request with non-existent username
func WithNonExistentUsername(req map[string]any) map[string]any {
	req["username"] = "nonexistentuser"
	delete(req, "email")
	return req
}

// WithEmptyCredentials modifies request with empty email and username
func WithEmptyCredentials(req map[string]any) map[string]any {
	req["email"] = ""
	req["username"] = ""
	return req
}

// WithMissingPassword modifies request without password field
func WithMissingPassword(req map[string]any) map[string]any {
	delete(req, "password")
	return req
}

// WithEmptyPassword modifies request with empty password
func WithEmptyPassword(req map[string]any) map[string]any {
	req["password"] = ""
	return req
}

// WithInvalidJSON provides malformed JSON string
func WithInvalidJSON() string {
	return "{invalid json"
}

// WithForceUserServiceFailure modifies service to force UserService.Login failure
func WithForceUserServiceFailure(userService *services.MockUserService) *services.MockUserService {
	userService.ShouldFailLogin = true
	userService.LoginErrorMessage = "User service login failure"
	return userService
}

// WithForceAuthServiceFailure modifies service to force AuthService.Update failure
type MockAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFail bool
}

func (m *MockAuthServiceWithFailure) Update(user repository.User) (*string, error) {
	if m.ShouldFail {
		return nil, fmt.Errorf("Auth service update failure")
	}
	return m.MockAuthService.Update(user)
}

// <runners/>

// NewTestRequest creates HTTP request from JSON data
func NewTestRequest(method string, path string, request map[string]any) *http.Request {
	requestJson, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(string(requestJson)),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewTestRequestWithBody creates HTTP request from string body
func NewTestRequestWithBody(method string, path string, body string) *http.Request {
	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewEchoContext creates new echo context
func NewEchoContext(r *http.Request) echo.Context {
	return echo.New().NewContext(r, httptest.NewRecorder())
}

// NewLoginHandler creates LoginHandler with services
func NewLoginHandler(r *http.Request, userService services.IUserService, authService services.IAuthService) *h.LoginHandler {
	return h.NewLoginHandler(
		NewEchoContext(r),
		services.ValidatorService{},
		userService,
		authService,
	)
}

// Run_LoginHandler_ValidEmail executes LoginHandler with valid email
func Run_LoginHandler_ValidEmail(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", LoginHandlerRequestParams())
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	data := result.Data().(struct {
		Authenticated *repository.User
		Token         *string
	})
	assert.NotNil(t, data.Authenticated)
	assert.NotNil(t, data.Token)
	assert.Equal(t, "handler@example.com", data.Authenticated.Email)
}

// Run_LoginHandler_ValidUsername executes LoginHandler with valid username
func Run_LoginHandler_ValidUsername(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", LoginHandlerRequestByUsername())
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	data := result.Data().(struct {
		Authenticated *repository.User
		Token         *string
	})
	assert.NotNil(t, data.Authenticated)
	assert.NotNil(t, data.Token)
	assert.Equal(t, "handleruser", data.Authenticated.Username)
}

// Run_LoginHandler_WrongPassword executes LoginHandler with wrong password
func Run_LoginHandler_WrongPassword(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithWrongPassword(LoginHandlerRequestParams()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
	assert.Contains(t, result.Error().Error(), "invalid password")
}

// Run_LoginHandler_NonExistentEmail executes LoginHandler with non-existent email
func Run_LoginHandler_NonExistentEmail(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithNonExistentEmail(LoginHandlerRequestParams()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
	assert.Contains(t, result.Error().Error(), "user not found")
}

// Run_LoginHandler_NonExistentUsername executes LoginHandler with non-existent username
func Run_LoginHandler_NonExistentUsername(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithNonExistentUsername(LoginHandlerRequestByUsername()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
	assert.Contains(t, result.Error().Error(), "user not found")
}

// Run_LoginHandler_EmptyCredentials executes LoginHandler with empty credentials
func Run_LoginHandler_EmptyCredentials(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithEmptyCredentials(LoginHandlerRequestParams()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_LoginHandler_MissingPassword executes LoginHandler with missing password
func Run_LoginHandler_MissingPassword(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithMissingPassword(LoginHandlerRequestParams()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_LoginHandler_EmptyPassword executes LoginHandler with empty password
func Run_LoginHandler_EmptyPassword(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", WithEmptyPassword(LoginHandlerRequestParams()))
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_LoginHandler_InvalidJSON executes LoginHandler with malformed JSON
func Run_LoginHandler_InvalidJSON(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &services.MockAuthService{}
	
	req := NewTestRequestWithBody(http.MethodPost, "/api/users/login", WithInvalidJSON())
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_LoginHandler_UserServiceFailure executes LoginHandler with UserService failure
func Run_LoginHandler_UserServiceFailure(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	WithForceUserServiceFailure(userService)
	authService := &services.MockAuthService{}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", LoginHandlerRequestParams())
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
	assert.Contains(t, result.Error().Error(), "User service login failure")
}

// Run_LoginHandler_AuthServiceFailure executes LoginHandler with AuthService failure
func Run_LoginHandler_AuthServiceFailure(t *testing.T) {
	userService := &services.MockUserService{}
	userService.Reset()
	authService := &MockAuthServiceWithFailure{ShouldFail: true}
	
	req := NewTestRequest(http.MethodPost, "/api/users/login", LoginHandlerRequestParams())
	handler := NewLoginHandler(req, userService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusInternalServerError, result.Code())
	assert.Contains(t, result.Error().Error(), "Auth service update failure")
}

// <tests>
// <evaluators>

// EvaluateLoginSuccess validates successful login handler execution
func EvaluateLoginSuccess(t *testing.T, result i.IHandler, expectedIdentifier string) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	data := result.Data().(struct {
		Authenticated *repository.User
		Token         *string
	})
	assert.NotNil(t, data.Authenticated)
	assert.NotNil(t, data.Token)
}

// EvaluateLoginFailure validates failed login handler execution
func EvaluateLoginFailure(t *testing.T, result i.IHandler, expectedCode int, expectedErrorSubstring string) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
	if expectedErrorSubstring != "" {
		assert.Contains(t, result.Error().Error(), expectedErrorSubstring)
	}
}

// <map/>

// LoginHandlerTestMap defines all LoginHandler test cases
var LoginHandlerTestMap = map[string]func(*testing.T){
	"ValidEmail":           Run_LoginHandler_ValidEmail,
	"ValidUsername":        Run_LoginHandler_ValidUsername,
	"WrongPassword":        Run_LoginHandler_WrongPassword,
	"NonExistentEmail":     Run_LoginHandler_NonExistentEmail,
	"NonExistentUsername":  Run_LoginHandler_NonExistentUsername,
	"EmptyCredentials":     Run_LoginHandler_EmptyCredentials,
	"MissingPassword":      Run_LoginHandler_MissingPassword,
	"EmptyPassword":        Run_LoginHandler_EmptyPassword,
	"InvalidJSON":          Run_LoginHandler_InvalidJSON,
	"UserServiceFailure":   Run_LoginHandler_UserServiceFailure,
	"AuthServiceFailure":   Run_LoginHandler_AuthServiceFailure,
}

// <hook/>

// Test_LoginHandler tests all LoginHandler scenarios
func Test_LoginHandler(t *testing.T) {
	fmt.Println("Test_LoginHandler")
	for name, testFunc := range LoginHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>