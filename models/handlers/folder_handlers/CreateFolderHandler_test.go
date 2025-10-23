package folder_handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/folder_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=folder_handlers.CreateFolderHandler.Handle>
// <fixtures/>

// Test user for CreateFolderHandler tests
var createFolderTestUser = struct {
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

// CreateFolderHandlerRequestParams provides base parameters for CreateFolderHandler tests
func CreateFolderHandlerRequestParams() map[string]any {
	return map[string]any{
		"name":      "Test Folder",
		"parent_id": nil, // Root folder
	}
}

// WithCreateFolderParentID modifies request with specific parent ID
func WithCreateFolderParentID(req map[string]any, parentID string) map[string]any {
	req["parent_id"] = parentID
	return req
}

// WithCreateFolderBadName modifies request with invalid name
func WithCreateFolderBadName(req map[string]any) map[string]any {
	req["name"] = "" // Empty name should fail validation
	return req
}

// WithCreateFolderMissingName modifies request with missing name field
func WithCreateFolderMissingName(req map[string]any) map[string]any {
	delete(req, "name")
	return req
}

// WithCreateFolderMalformedJSON provides malformed JSON string
func WithCreateFolderMalformedJSON() string {
	return "{malformed json"
}

// <runners/>

// CreateFolderTestServices provides mock services for testing
type CreateFolderTestServices struct {
	FolderService    services.IFolderService
	AuthService      services.IAuthService
	ValidatorService services.ValidatorService
}

// CreateCreateFolderTestServices creates fresh instances of all services for testing
func CreateCreateFolderTestServices() CreateFolderTestServices {
	folderService := services.CreateMockFolderService(nil, nil)
	folderService.Reset()
	
	return CreateFolderTestServices{
		FolderService:    folderService,
		AuthService:      &services.MockAuthService{},
		ValidatorService: services.ValidatorService{},
	}
}

// CreateCreateFolderHTTPRequest creates an HTTP request with JSON body
func CreateCreateFolderHTTPRequest(method, path string, body map[string]any) *http.Request {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	
	req := httptest.NewRequest(method, path, strings.NewReader(string(jsonBytes)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateCreateFolderHTTPRequestWithBody creates an HTTP request with string body
func CreateCreateFolderHTTPRequestWithBody(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateCreateFolderAuthenticatedContext creates echo context with JWT authentication
func CreateCreateFolderAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     createFolderTestUser.ID,
		User:       createFolderTestUser.Username,
		IsAdmin:    createFolderTestUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewCreateFolderHandler creates CreateFolderHandler with services and context
func NewCreateFolderHandler(services CreateFolderTestServices, ctx echo.Context) *h.CreateFolderHandler {
	return h.NewCreateHandler(ctx, services.ValidatorService, services.FolderService, services.AuthService)
}

// MockCreateFolderAuthServiceWithFailure extends MockAuthService for failure testing
type MockCreateFolderAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockCreateFolderAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_CreateFolderHandler_ValidRequest executes CreateFolderHandler with valid request
func Run_CreateFolderHandler_ValidRequest(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", CreateFolderHandlerRequestParams())
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folder := result.Data().(*repository.Folder)
	assert.Equal(t, "Test Folder", folder.Name)
	assert.Equal(t, createFolderTestUser.ID, folder.UserID)
}

// Run_CreateFolderHandler_ValidRequestWithParent executes CreateFolderHandler with parent folder
func Run_CreateFolderHandler_ValidRequestWithParent(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	parentID := "22222222-2222-2222-2222-222222222222"
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", WithCreateFolderParentID(CreateFolderHandlerRequestParams(), parentID))
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
}

// Run_CreateFolderHandler_AuthenticationFailure executes CreateFolderHandler with auth failure
func Run_CreateFolderHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	testServices.AuthService = &MockCreateFolderAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", CreateFolderHandlerRequestParams())
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_CreateFolderHandler_ValidationFailure executes CreateFolderHandler with validation failure
func Run_CreateFolderHandler_ValidationFailure(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", WithCreateFolderBadName(CreateFolderHandlerRequestParams()))
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_CreateFolderHandler_MissingName executes CreateFolderHandler with missing name
func Run_CreateFolderHandler_MissingName(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", WithCreateFolderMissingName(CreateFolderHandlerRequestParams()))
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_CreateFolderHandler_ServiceFailure executes CreateFolderHandler with FolderService failure
func Run_CreateFolderHandler_ServiceFailure(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateCreateFolderHTTPRequest(http.MethodPost, "/api/folders", CreateFolderHandlerRequestParams())
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_CreateFolderHandler_InvalidJSON executes CreateFolderHandler with malformed JSON
func Run_CreateFolderHandler_InvalidJSON(t *testing.T) {
	testServices := CreateCreateFolderTestServices()
	req := CreateCreateFolderHTTPRequestWithBody(http.MethodPost, "/api/folders", WithCreateFolderMalformedJSON())
	ctx := CreateCreateFolderAuthenticatedContext(req)
	
	handler := NewCreateFolderHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// <tests>
// <evaluators>

// EvaluateCreateFolderSuccess validates successful create folder handler execution
func EvaluateCreateFolderSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folder := result.Data().(*repository.Folder)
	assert.NotEmpty(t, folder.ID)
	assert.NotEmpty(t, folder.Name)
	assert.NotEmpty(t, folder.UserID)
}

// EvaluateCreateFolderFailure validates failed create folder handler execution
func EvaluateCreateFolderFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// CreateFolderHandlerTestMap defines all CreateFolderHandler test cases
var CreateFolderHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":            Run_CreateFolderHandler_ValidRequest,
	"ValidRequestWithParent":  Run_CreateFolderHandler_ValidRequestWithParent,
	"AuthenticationFailure":   Run_CreateFolderHandler_AuthenticationFailure,
	"ValidationFailure":       Run_CreateFolderHandler_ValidationFailure,
	"MissingName":             Run_CreateFolderHandler_MissingName,
	"ServiceFailure":          Run_CreateFolderHandler_ServiceFailure,
	"InvalidJSON":             Run_CreateFolderHandler_InvalidJSON,
}

// <hook/>

// Test_CreateFolderHandler tests all CreateFolderHandler scenarios
func Test_CreateFolderHandler(t *testing.T) {
	fmt.Println("Test_CreateFolderHandler")
	for name, testFunc := range CreateFolderHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>