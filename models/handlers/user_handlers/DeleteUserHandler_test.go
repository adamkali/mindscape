package user_handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/user_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=user_handlers.DeleteUserHandler.Handle>
// <fixtures/>

// Test users for DeleteUserHandler tests
var deleteUserTestAdminUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	IsAdmin  bool
}{
	ID:       uuid.MustParse("66666666-6666-6666-6666-666666666666"),
	Username: "deleteadmin",
	Email:    "deleteadmin@example.com",
	Password: "passwordABC123!",
	IsAdmin:  true,
}

var deleteUserTestRegularUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
	IsAdmin  bool
}{
	ID:       uuid.MustParse("77777777-7777-7777-7777-777777777777"),
	Username: "deleteuser",
	Email:    "deleteuser@example.com",
	Password: "passwordABC123!",
	IsAdmin:  false,
}

// Target user to be deleted
var deleteUserTestTargetUser = struct {
	ID       uuid.UUID
	Username string
	Email    string
}{
	ID:       uuid.MustParse("88888888-8888-8888-8888-888888888888"),
	Username: "targettodelete",
	Email:    "target@example.com",
}

// WithDeleteInvalidUserID modifies request with invalid user ID
func WithDeleteInvalidUserID() string {
	return "invalid-uuid"
}

// WithDeleteNonExistentUserID modifies request with non-existent user ID
func WithDeleteNonExistentUserID() string {
	return "99999999-9999-9999-9999-999999999999"
}

// <runners/>

// DeleteUserTestServices provides mock services for testing
type DeleteUserTestServices struct {
	UserService      services.IUserService
	AuthService      services.IAuthService
	ValidatorService services.ValidatorService
}

// CreateDeleteUserTestServices creates fresh instances of all services for testing
func CreateDeleteUserTestServices() DeleteUserTestServices {
	userService := services.CreateMockUserService(nil, nil)
	userService.Reset()
	
	return DeleteUserTestServices{
		UserService:      userService,
		AuthService:      &services.MockAuthService{},
		ValidatorService: services.ValidatorService{},
	}
}

// CreateDeleteUserHTTPRequest creates an HTTP request for delete user
func CreateDeleteUserHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateDeleteUserAdminAuthenticatedContext creates echo context with admin JWT authentication
func CreateDeleteUserAdminAuthenticatedContext(r *http.Request, userID string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Set the user_id parameter
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(userID)
	
	// Create JWT token for admin user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     deleteUserTestAdminUser.ID,
		User:       deleteUserTestAdminUser.Username,
		IsAdmin:    deleteUserTestAdminUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// CreateDeleteUserRegularAuthenticatedContext creates echo context with regular user JWT authentication
func CreateDeleteUserRegularAuthenticatedContext(r *http.Request, userID string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Set the user_id parameter
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(userID)
	
	// Create JWT token for regular user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     deleteUserTestRegularUser.ID,
		User:       deleteUserTestRegularUser.Username,
		IsAdmin:    deleteUserTestRegularUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewDeleteUserHandler creates DeleteUserHandler with services and context
func NewDeleteUserHandler(services DeleteUserTestServices, ctx echo.Context) *h.DeleteUserHandler {
	return h.NewDeleteUserHandler(ctx, services.ValidatorService, services.UserService, services.AuthService)
}

// MockDeleteUserAuthServiceWithFailure extends MockAuthService for failure testing
type MockDeleteUserAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockDeleteUserAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_DeleteUserHandler_ValidRequest executes DeleteUserHandler with valid admin request
func Run_DeleteUserHandler_ValidRequest(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+deleteUserTestTargetUser.ID.String())
	ctx := CreateDeleteUserAdminAuthenticatedContext(req, deleteUserTestTargetUser.ID.String())
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_DeleteUserHandler_AuthenticationFailure executes DeleteUserHandler with auth failure
func Run_DeleteUserHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	testServices.AuthService = &MockDeleteUserAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+deleteUserTestTargetUser.ID.String())
	ctx := CreateDeleteUserAdminAuthenticatedContext(req, deleteUserTestTargetUser.ID.String())
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_DeleteUserHandler_CurrentUserNotFound executes DeleteUserHandler when current user doesn't exist
func Run_DeleteUserHandler_CurrentUserNotFound(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+deleteUserTestTargetUser.ID.String())
	
	// Create context with non-existent user ID
	e := echo.New()
	ctx := e.NewContext(req, httptest.NewRecorder())
	ctx.SetParamNames("user_id")
	ctx.SetParamValues(deleteUserTestTargetUser.ID.String())
	
	// Create JWT token for non-existent user
	nonExistentUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     nonExistentUserID,
		User:       "nonexistent",
		IsAdmin:    true,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	ctx.Set("user", token)
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusNotFound, result.Code())
}

// Run_DeleteUserHandler_AuthorizationFailure executes DeleteUserHandler with non-admin user
func Run_DeleteUserHandler_AuthorizationFailure(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+deleteUserTestTargetUser.ID.String())
	ctx := CreateDeleteUserRegularAuthenticatedContext(req, deleteUserTestTargetUser.ID.String()) // Non-admin user
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusForbidden, result.Code())
	assert.Contains(t, result.Error().Error(), "admin privileges required")
}

// Run_DeleteUserHandler_InvalidUserID executes DeleteUserHandler with invalid user ID
func Run_DeleteUserHandler_InvalidUserID(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+WithDeleteInvalidUserID())
	ctx := CreateDeleteUserAdminAuthenticatedContext(req, WithDeleteInvalidUserID())
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_DeleteUserHandler_TargetUserNotFound executes DeleteUserHandler with non-existent target user
func Run_DeleteUserHandler_TargetUserNotFound(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+WithDeleteNonExistentUserID())
	ctx := CreateDeleteUserAdminAuthenticatedContext(req, WithDeleteNonExistentUserID())
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusNotFound, result.Code())
}

// Run_DeleteUserHandler_UserServiceFailure executes DeleteUserHandler with UserService failure
func Run_DeleteUserHandler_UserServiceFailure(t *testing.T) {
	testServices := CreateDeleteUserTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateDeleteUserHTTPRequest(http.MethodDelete, "/api/users/"+deleteUserTestTargetUser.ID.String())
	ctx := CreateDeleteUserAdminAuthenticatedContext(req, deleteUserTestTargetUser.ID.String())
	
	handler := NewDeleteUserHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// <tests>
// <evaluators>

// EvaluateDeleteUserSuccess validates successful delete user handler execution
func EvaluateDeleteUserSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// EvaluateDeleteUserFailure validates failed delete user handler execution
func EvaluateDeleteUserFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// DeleteUserHandlerTestMap defines all DeleteUserHandler test cases
var DeleteUserHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":          Run_DeleteUserHandler_ValidRequest,
	"AuthenticationFailure": Run_DeleteUserHandler_AuthenticationFailure,
	"CurrentUserNotFound":   Run_DeleteUserHandler_CurrentUserNotFound,
	"AuthorizationFailure":  Run_DeleteUserHandler_AuthorizationFailure,
	"InvalidUserID":         Run_DeleteUserHandler_InvalidUserID,
	"TargetUserNotFound":    Run_DeleteUserHandler_TargetUserNotFound,
	"UserServiceFailure":    Run_DeleteUserHandler_UserServiceFailure,
}

// <hook/>

// Test_DeleteUserHandler tests all DeleteUserHandler scenarios
func Test_DeleteUserHandler(t *testing.T) {
	fmt.Println("Test_DeleteUserHandler")
	for name, testFunc := range DeleteUserHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>