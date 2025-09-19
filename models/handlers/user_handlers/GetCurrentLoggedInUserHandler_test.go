package user_handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/user_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=user_handlers.GetCurrentLoggedInUserHandler.Handle>
// <fixtures/>

// Test user for GetCurrentLoggedInUserHandler tests
var getCurrentUserTestUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	IsAdmin  bool
}{
	ID:       uuid.MustParse("33333333-3333-3333-3333-333333333333"), // Using existing handleruser
	Username: "handleruser",
	Email:    "handler@example.com",
	Password: "passwordABC123!",
	IsAdmin:  false,
}

// <runners/>

// GetCurrentUserTestServices provides mock services for testing
type GetCurrentUserTestServices struct {
	UserService      services.IUserService
	AuthService      services.IAuthService
	ValidatorService services.ValidatorService
}

// CreateGetCurrentUserTestServices creates fresh instances of all services for testing
func CreateGetCurrentUserTestServices() GetCurrentUserTestServices {
	userService := services.CreateMockUserService(nil, nil)
	userService.Reset()
	
	return GetCurrentUserTestServices{
		UserService:      userService,
		AuthService:      &services.MockAuthService{},
		ValidatorService: services.ValidatorService{},
	}
}

// CreateGetCurrentUserHTTPRequest creates an HTTP request for get current user
func CreateGetCurrentUserHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateGetCurrentUserAuthenticatedContext creates echo context with JWT authentication
func CreateGetCurrentUserAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     getCurrentUserTestUser.ID,
		User:       getCurrentUserTestUser.Username,
		IsAdmin:    getCurrentUserTestUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// CreateGetCurrentUserNonExistentAuthenticatedContext creates echo context with non-existent user JWT
func CreateGetCurrentUserNonExistentAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token for non-existent user
	nonExistentUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     nonExistentUserID,
		User:       "nonexistent",
		IsAdmin:    false,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewGetCurrentUserHandler creates GetCurrentLoggedInUserHandler with services and context
func NewGetCurrentUserHandler(services GetCurrentUserTestServices, ctx echo.Context) *h.GetCurrentLoggedInUserHandler {
	return h.NewGetCurrentLoggedInUserHandler(ctx, services.ValidatorService, services.UserService, services.AuthService)
}

// MockGetCurrentUserAuthServiceWithFailure extends MockAuthService for failure testing
type MockGetCurrentUserAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockGetCurrentUserAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_GetCurrentUserHandler_ValidRequest executes GetCurrentLoggedInUserHandler with valid request
func Run_GetCurrentUserHandler_ValidRequest(t *testing.T) {
	testServices := CreateGetCurrentUserTestServices()
	req := CreateGetCurrentUserHTTPRequest(http.MethodGet, "/api/users/me")
	ctx := CreateGetCurrentUserAuthenticatedContext(req)
	
	handler := NewGetCurrentUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	user := result.Data().(*repository.User)
	assert.Equal(t, getCurrentUserTestUser.ID, user.ID)
	assert.Equal(t, getCurrentUserTestUser.Username, user.Username)
	assert.Equal(t, getCurrentUserTestUser.Email, user.Email)
}

// Run_GetCurrentUserHandler_AuthenticationFailure executes GetCurrentLoggedInUserHandler with auth failure
func Run_GetCurrentUserHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateGetCurrentUserTestServices()
	testServices.AuthService = &MockGetCurrentUserAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateGetCurrentUserHTTPRequest(http.MethodGet, "/api/users/me")
	ctx := CreateGetCurrentUserAuthenticatedContext(req)
	
	handler := NewGetCurrentUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_GetCurrentUserHandler_UserNotFound executes GetCurrentLoggedInUserHandler when user doesn't exist
func Run_GetCurrentUserHandler_UserNotFound(t *testing.T) {
	testServices := CreateGetCurrentUserTestServices()
	req := CreateGetCurrentUserHTTPRequest(http.MethodGet, "/api/users/me")
	ctx := CreateGetCurrentUserNonExistentAuthenticatedContext(req)
	
	handler := NewGetCurrentUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusNotFound, result.Code())
}

// <tests>
// <evaluators>

// EvaluateGetCurrentUserSuccess validates successful get current user handler execution
func EvaluateGetCurrentUserSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	user := result.Data().(*repository.User)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Email)
}

// EvaluateGetCurrentUserFailure validates failed get current user handler execution
func EvaluateGetCurrentUserFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// GetCurrentUserHandlerTestMap defines all GetCurrentLoggedInUserHandler test cases
var GetCurrentUserHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":          Run_GetCurrentUserHandler_ValidRequest,
	"AuthenticationFailure": Run_GetCurrentUserHandler_AuthenticationFailure,
	"UserNotFound":          Run_GetCurrentUserHandler_UserNotFound,
}

// <hook/>

// Test_GetCurrentLoggedInUserHandler tests all GetCurrentLoggedInUserHandler scenarios
func Test_GetCurrentLoggedInUserHandler(t *testing.T) {
	fmt.Println("Test_GetCurrentLoggedInUserHandler")
	for name, testFunc := range GetCurrentUserHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>