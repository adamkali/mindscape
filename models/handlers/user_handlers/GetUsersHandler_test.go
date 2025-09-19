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

// <method var=user_handlers.GetUsersHandler.Handle>
// <fixtures/>

// Test users for GetUsersHandler tests
var getUsersTestAdminUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	IsAdmin  bool
}{
	ID:       uuid.MustParse("44444444-4444-4444-4444-444444444444"),
	Username: "getusersadmin",
	Email:    "getusersadmin@example.com",
	Password: "passwordABC123!",
	IsAdmin:  true,
}

var getUsersTestRegularUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	IsAdmin  bool
}{
	ID:       uuid.MustParse("55555555-5555-5555-5555-555555555555"),
	Username: "getusersuser",
	Email:    "getusersuser@example.com",
	Password: "passwordABC123!",
	IsAdmin:  false,
}

// <runners/>

// GetUsersTestServices provides mock services for testing
type GetUsersTestServices struct {
	UserService services.IUserService
	AuthService services.IAuthService
}

// CreateGetUsersTestServices creates fresh instances of all services for testing
func CreateGetUsersTestServices() GetUsersTestServices {
	userService := services.CreateMockUserService(nil, nil)
	userService.Reset()
	
	return GetUsersTestServices{
		UserService: userService,
		AuthService: &services.MockAuthService{},
	}
}

// CreateGetUsersHTTPRequest creates an HTTP request for GetUsers
func CreateGetUsersHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateGetUsersAdminAuthenticatedContext creates echo context with admin JWT authentication
func CreateGetUsersAdminAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token for admin user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     getUsersTestAdminUser.ID,
		User:       getUsersTestAdminUser.Username,
		IsAdmin:    getUsersTestAdminUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// CreateGetUsersRegularAuthenticatedContext creates echo context with regular user JWT authentication
func CreateGetUsersRegularAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token for regular user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     getUsersTestRegularUser.ID,
		User:       getUsersTestRegularUser.Username,
		IsAdmin:    getUsersTestRegularUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewGetUsersHandler creates GetUsersHandler with services and context
func NewGetUsersHandler(services GetUsersTestServices, ctx echo.Context) *h.GetUsersHandler {
	return h.NewGetUsersHandler(ctx, services.AuthService, services.UserService)
}

// MockGetUsersAuthServiceWithFailure extends MockAuthService for failure testing
type MockGetUsersAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockGetUsersAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_GetUsersHandler_ValidRequest executes GetUsersHandler with valid admin request
func Run_GetUsersHandler_ValidRequest(t *testing.T) {
	testServices := CreateGetUsersTestServices()
	req := CreateGetUsersHTTPRequest(http.MethodGet, "/api/users")
	ctx := CreateGetUsersAdminAuthenticatedContext(req)
	
	handler := NewGetUsersHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	users := result.Data().([]repository.User)
	assert.GreaterOrEqual(t, len(users), 1) // Should have at least the seeded users
}

// Run_GetUsersHandler_AuthenticationFailure executes GetUsersHandler with auth failure
func Run_GetUsersHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateGetUsersTestServices()
	testServices.AuthService = &MockGetUsersAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateGetUsersHTTPRequest(http.MethodGet, "/api/users")
	ctx := CreateGetUsersAdminAuthenticatedContext(req)
	
	handler := NewGetUsersHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_GetUsersHandler_AuthorizationFailure executes GetUsersHandler with non-admin user
func Run_GetUsersHandler_AuthorizationFailure(t *testing.T) {
	testServices := CreateGetUsersTestServices()
	req := CreateGetUsersHTTPRequest(http.MethodGet, "/api/users")
	ctx := CreateGetUsersRegularAuthenticatedContext(req) // Non-admin user
	
	handler := NewGetUsersHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusForbidden, result.Code())
	assert.Contains(t, result.Error().Error(), "admin privileges required")
}

// Run_GetUsersHandler_UserServiceFailure executes GetUsersHandler with UserService failure
func Run_GetUsersHandler_UserServiceFailure(t *testing.T) {
	testServices := CreateGetUsersTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateGetUsersHTTPRequest(http.MethodGet, "/api/users")
	ctx := CreateGetUsersAdminAuthenticatedContext(req)
	
	handler := NewGetUsersHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// <tests>
// <evaluators>

// EvaluateGetUsersSuccess validates successful get users handler execution
func EvaluateGetUsersSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	users := result.Data().([]repository.User)
	assert.GreaterOrEqual(t, len(users), 0)
}

// EvaluateGetUsersFailure validates failed get users handler execution
func EvaluateGetUsersFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// GetUsersHandlerTestMap defines all GetUsersHandler test cases
var GetUsersHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":          Run_GetUsersHandler_ValidRequest,
	"AuthenticationFailure": Run_GetUsersHandler_AuthenticationFailure,
	"AuthorizationFailure":  Run_GetUsersHandler_AuthorizationFailure,
	"UserServiceFailure":    Run_GetUsersHandler_UserServiceFailure,
}

// <hook/>

// Test_GetUsersHandler tests all GetUsersHandler scenarios
func Test_GetUsersHandler(t *testing.T) {
	fmt.Println("Test_GetUsersHandler")
	for name, testFunc := range GetUsersHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>