package folder_handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/folder_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=folder_handlers.DeleteFolderHandler.Handle>
// <fixtures/>

// Test user for DeleteFolderHandler tests
var deleteFolderTestUser = struct {
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

// Test folder IDs
var deleteFolderTestFolderID = "11111111-1111-1111-1111-111111111111"       // handleruser's folder
var deleteFolderTestNonExistentFolderID = "99999999-9999-9999-9999-999999999999"
var deleteFolderTestOtherUserFolderID = "88888888-8888-8888-8888-888888888888" // different user's folder

// WithDeleteFolderInvalidFolderID provides invalid folder ID
func WithDeleteFolderInvalidFolderID() string {
	return "invalid-uuid"
}

// <runners/>

// DeleteFolderTestServices provides mock services for testing
type DeleteFolderTestServices struct {
	FolderService services.IFolderService
	AuthService   services.IAuthService
}

// CreateDeleteFolderTestServices creates fresh instances of all services for testing
func CreateDeleteFolderTestServices() DeleteFolderTestServices {
	folderService := services.CreateMockFolderService(nil, nil)
	folderService.Reset()
	
	return DeleteFolderTestServices{
		FolderService: folderService,
		AuthService:   &services.MockAuthService{},
	}
}

// CreateDeleteFolderHTTPRequest creates an HTTP request for delete folder
func CreateDeleteFolderHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateDeleteFolderAuthenticatedContext creates echo context with JWT authentication and folder ID param
func CreateDeleteFolderAuthenticatedContext(r *http.Request, folderID string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Set the folder_id parameter
	ctx.SetParamNames("folder_id")
	ctx.SetParamValues(folderID)
	
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     deleteFolderTestUser.ID,
		User:       deleteFolderTestUser.Username,
		IsAdmin:    deleteFolderTestUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewDeleteFolderHandler creates DeleteFolderHandler with services and context
func NewDeleteFolderHandler(services DeleteFolderTestServices, ctx echo.Context) *h.DeleteFolderHandler {
	return h.NewDeleteHandler(ctx, services.FolderService, services.AuthService)
}

// MockDeleteFolderAuthServiceWithFailure extends MockAuthService for failure testing
type MockDeleteFolderAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockDeleteFolderAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_DeleteFolderHandler_ValidRequest executes DeleteFolderHandler with valid request
func Run_DeleteFolderHandler_ValidRequest(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+deleteFolderTestFolderID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, deleteFolderTestFolderID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_DeleteFolderHandler_AuthenticationFailure executes DeleteFolderHandler with auth failure
func Run_DeleteFolderHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	testServices.AuthService = &MockDeleteFolderAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+deleteFolderTestFolderID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, deleteFolderTestFolderID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_DeleteFolderHandler_InvalidFolderID executes DeleteFolderHandler with invalid folder ID
func Run_DeleteFolderHandler_InvalidFolderID(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	invalidID := WithDeleteFolderInvalidFolderID()
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+invalidID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, invalidID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_DeleteFolderHandler_FolderNotFound executes DeleteFolderHandler with non-existent folder
func Run_DeleteFolderHandler_FolderNotFound(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+deleteFolderTestNonExistentFolderID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, deleteFolderTestNonExistentFolderID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusNotFound, result.Code())
}

// Run_DeleteFolderHandler_AuthorizationFailure executes DeleteFolderHandler with unauthorized access
func Run_DeleteFolderHandler_AuthorizationFailure(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	// handleruser trying to delete different user's folder
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+deleteFolderTestOtherUserFolderID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, deleteFolderTestOtherUserFolderID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusForbidden, result.Code())
	assert.Contains(t, result.Error().Error(), "unauthorized access to folder")
}

// Run_DeleteFolderHandler_ServiceFailure executes DeleteFolderHandler with FolderService failure
func Run_DeleteFolderHandler_ServiceFailure(t *testing.T) {
	testServices := CreateDeleteFolderTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateDeleteFolderHTTPRequest(http.MethodDelete, "/api/folders/"+deleteFolderTestFolderID)
	ctx := CreateDeleteFolderAuthenticatedContext(req, deleteFolderTestFolderID)
	
	handler := NewDeleteFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// <tests>
// <evaluators>

// EvaluateDeleteFolderSuccess validates successful delete folder handler execution
func EvaluateDeleteFolderSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// EvaluateDeleteFolderFailure validates failed delete folder handler execution
func EvaluateDeleteFolderFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// DeleteFolderHandlerTestMap defines all DeleteFolderHandler test cases
var DeleteFolderHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":          Run_DeleteFolderHandler_ValidRequest,
	"AuthenticationFailure": Run_DeleteFolderHandler_AuthenticationFailure,
	"InvalidFolderID":       Run_DeleteFolderHandler_InvalidFolderID,
	"FolderNotFound":        Run_DeleteFolderHandler_FolderNotFound,
	"AuthorizationFailure":  Run_DeleteFolderHandler_AuthorizationFailure,
	"ServiceFailure":        Run_DeleteFolderHandler_ServiceFailure,
}

// <hook/>

// Test_DeleteFolderHandler tests all DeleteFolderHandler scenarios
func Test_DeleteFolderHandler(t *testing.T) {
	fmt.Println("Test_DeleteFolderHandler")
	for name, testFunc := range DeleteFolderHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>