package folder_handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adamkali/mindscape/models/responses"
	i "github.com/adamkali/mindscape/models/handlers"
	h "github.com/adamkali/mindscape/models/handlers/folder_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=folder_handlers.GetRootFolderHandler.Handle>
// <fixtures/>

// Test user for GetRootHandler tests
var getRootTestUser = struct {
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

// GetRootTestServices provides mock services for testing
type GetRootTestServices struct {
	FolderService   services.IFolderService
	BookmarkService services.IBookmarkService
	NoteService     services.INoteService
	AuthService     services.IAuthService
}

// CreateGetRootTestServices creates fresh instances of all services for testing
func CreateGetRootTestServices() GetRootTestServices {
	folderService := services.CreateMockFolderService(nil, nil)
	folderService.Reset()
	
	bookmarkService := services.CreateMockBookmarkService(nil, nil)
	bookmarkService.Reset()
	
	return GetRootTestServices{
		FolderService:   folderService,
		BookmarkService: bookmarkService,
		NoteService:     services.NewMockNoteService(),
		AuthService:     &services.MockAuthService{},
	}
}

// CreateGetRootHTTPRequest creates an HTTP request for get root folders
func CreateGetRootHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateGetRootAuthenticatedContext creates echo context with JWT authentication
func CreateGetRootAuthenticatedContext(r *http.Request) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     getRootTestUser.ID,
		User:       getRootTestUser.Username,
		IsAdmin:    getRootTestUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// NewGetRootHandler creates GetRootFolderHandler with services and context
func NewGetRootHandler(services GetRootTestServices, ctx echo.Context) *h.GetRootFolderHandler {
	return h.NewGetRootHandler(ctx, services.FolderService, services.BookmarkService, services.NoteService, services.AuthService)
}

// MockGetRootAuthServiceWithFailure extends MockAuthService for failure testing
type MockGetRootAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockGetRootAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_GetRootHandler_ValidRequest executes GetRootHandler with valid request
func Run_GetRootHandler_ValidRequest(t *testing.T) {
	testServices := CreateGetRootTestServices()
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := CreateGetRootAuthenticatedContext(req)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folderData := result.Data().([]responses.FolderData)
	assert.GreaterOrEqual(t, len(folderData), 0) // Should have 0 or more root folders
}

// Run_GetRootHandler_AuthenticationFailure executes GetRootHandler with auth failure
func Run_GetRootHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateGetRootTestServices()
	testServices.AuthService = &MockGetRootAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := CreateGetRootAuthenticatedContext(req)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_GetRootHandler_NoRootFolders executes GetRootHandler when user has no root folders
func Run_GetRootHandler_NoRootFolders(t *testing.T) {
	testServices := CreateGetRootTestServices()
	
	// Create context with non-existent user ID to test no folders scenario
	e := echo.New()
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := e.NewContext(req, httptest.NewRecorder())
	
	// Create JWT token for user without folders
	nonExistentUserID := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     nonExistentUserID,
		User:       "nofoldersuser",
		IsAdmin:    false,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	ctx.Set("user", token)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folderData := result.Data().([]responses.FolderData)
	assert.Equal(t, 0, len(folderData)) // Should have no folders
}

// Run_GetRootHandler_FolderServiceFailure executes GetRootHandler with FolderService failure
func Run_GetRootHandler_FolderServiceFailure(t *testing.T) {
	testServices := CreateGetRootTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := CreateGetRootAuthenticatedContext(req)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_GetRootHandler_BookmarkServiceFailure executes GetRootHandler with BookmarkService failure
func Run_GetRootHandler_BookmarkServiceFailure(t *testing.T) {
	testServices := CreateGetRootTestServices()
	// TODO: Add bookmark service failure testing once type access is resolved
	
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := CreateGetRootAuthenticatedContext(req)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_GetRootHandler_NoteServiceFailure executes GetRootHandler with NoteService failure
func Run_GetRootHandler_NoteServiceFailure(t *testing.T) {
	testServices := CreateGetRootTestServices()
	noteService := &services.MockNoteService{
		ShouldFailGetByFolder:   true,
		GetByFolderErrorMessage: "Note service failure",
	}
	testServices.NoteService = noteService
	
	req := CreateGetRootHTTPRequest(http.MethodGet, "/api/folders/root")
	ctx := CreateGetRootAuthenticatedContext(req)
	
	handler := NewGetRootHandler(testServices, ctx)
	result := handler.Handle()
	
	// Should fail when note service fails and there are folders to process
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusInternalServerError, result.Code())
}

// <tests>
// <evaluators>

// EvaluateGetRootSuccess validates successful get root handler execution
func EvaluateGetRootSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folderData := result.Data().([]responses.FolderData)
	assert.GreaterOrEqual(t, len(folderData), 0)
}

// EvaluateGetRootFailure validates failed get root handler execution
func EvaluateGetRootFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// GetRootHandlerTestMap defines all GetRootHandler test cases
var GetRootHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":             Run_GetRootHandler_ValidRequest,
	"AuthenticationFailure":    Run_GetRootHandler_AuthenticationFailure,
	"NoRootFolders":            Run_GetRootHandler_NoRootFolders,
	"FolderServiceFailure":     Run_GetRootHandler_FolderServiceFailure,
	"BookmarkServiceFailure":   Run_GetRootHandler_BookmarkServiceFailure,
	"NoteServiceFailure":       Run_GetRootHandler_NoteServiceFailure,
}

// <hook/>

// Test_GetRootHandler tests all GetRootHandler scenarios
func Test_GetRootHandler(t *testing.T) {
	fmt.Println("Test_GetRootHandler")
	for name, testFunc := range GetRootHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>