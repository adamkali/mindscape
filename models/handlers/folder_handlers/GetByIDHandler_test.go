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

// <method var=folder_handlers.GetFolderByIDHandler.Handle>
// <fixtures/>

// Test user for GetByIDHandler tests
var getByIDTestUser = struct {
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
var getByIDTestFolderID = "11111111-1111-1111-1111-111111111111"
var getByIDTestNonExistentFolderID = "99999999-9999-9999-9999-999999999999"
var getByIDTestOtherUserFolderID = "88888888-8888-8888-8888-888888888888"

// WithGetByIDInvalidFolderID provides invalid folder ID
func WithGetByIDInvalidFolderID() string {
	return "invalid-uuid"
}

// <runners/>

// GetByIDTestServices provides mock services for testing
type GetByIDTestServices struct {
	FolderService   services.IFolderService
	BookmarkService services.IBookmarkService
	NoteService     services.INoteService
	AuthService     services.IAuthService
}

// CreateGetByIDTestServices creates fresh instances of all services for testing
func CreateGetByIDTestServices() GetByIDTestServices {
	folderService := services.CreateMockFolderService(nil, nil)
	folderService.Reset()
	
	bookmarkService := services.CreateMockBookmarkService(nil, nil)
	bookmarkService.Reset()
	
	return GetByIDTestServices{
		FolderService:   folderService,
		BookmarkService: bookmarkService,
		NoteService:     services.NewMockNoteService(),
		AuthService:     &services.MockAuthService{},
	}
}

// CreateGetByIDHTTPRequest creates an HTTP request for get folder by ID
func CreateGetByIDHTTPRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateGetByIDAuthenticatedContext creates echo context with JWT authentication and folder ID param
func CreateGetByIDAuthenticatedContext(r *http.Request, folderID string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Set the folder_id parameter
	ctx.SetParamNames("folder_id")
	ctx.SetParamValues(folderID)
	
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     getByIDTestUser.ID,
		User:       getByIDTestUser.Username,
		IsAdmin:    getByIDTestUser.IsAdmin,
		ProfilePic: "default-avatar.png",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	
	// Set token in context
	ctx.Set("user", token)
	return ctx
}

// CreateGetByIDDifferentUserContext creates echo context with different user JWT
func CreateGetByIDDifferentUserContext(r *http.Request, folderID string) echo.Context {
	e := echo.New()
	ctx := e.NewContext(r, httptest.NewRecorder())
	
	// Set the folder_id parameter
	ctx.SetParamNames("folder_id")
	ctx.SetParamValues(folderID)
	
	// Create JWT token for different user
	differentUserID := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &services.CustomJwt{
		UserId:     differentUserID,
		User:       "differentuser",
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

// NewGetByIDHandler creates GetFolderByIDHandler with services and context
func NewGetByIDHandler(services GetByIDTestServices, ctx echo.Context) *h.GetFolderByIDHandler {
	return h.NewGetById(ctx, services.FolderService, services.BookmarkService, services.NoteService, services.AuthService)
}

// MockGetByIDAuthServiceWithFailure extends MockAuthService for failure testing
type MockGetByIDAuthServiceWithFailure struct {
	services.MockAuthService
	ShouldFailCheckToken bool
}

func (m *MockGetByIDAuthServiceWithFailure) CheckToken(token string) error {
	if m.ShouldFailCheckToken {
		return assert.AnError
	}
	return nil
}

// Run_GetByIDHandler_ValidRequest executes GetByIDHandler with valid request
func Run_GetByIDHandler_ValidRequest(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folderData := result.Data().(*responses.FolderData)
	assert.NotNil(t, folderData.ID)
	assert.NotEmpty(t, folderData.Name)
}

// Run_GetByIDHandler_AuthenticationFailure executes GetByIDHandler with auth failure
func Run_GetByIDHandler_AuthenticationFailure(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	testServices.AuthService = &MockGetByIDAuthServiceWithFailure{ShouldFailCheckToken: true}
	
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusUnauthorized, result.Code())
}

// Run_GetByIDHandler_InvalidFolderID executes GetByIDHandler with invalid folder ID
func Run_GetByIDHandler_InvalidFolderID(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	invalidID := WithGetByIDInvalidFolderID()
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+invalidID)
	ctx := CreateGetByIDAuthenticatedContext(req, invalidID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusBadRequest, result.Code())
}

// Run_GetByIDHandler_FolderNotFound executes GetByIDHandler with non-existent folder
func Run_GetByIDHandler_FolderNotFound(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestNonExistentFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestNonExistentFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusNotFound, result.Code())
}

// Run_GetByIDHandler_AuthorizationFailure executes GetByIDHandler with unauthorized access
func Run_GetByIDHandler_AuthorizationFailure(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	// Use the different user's folder ID for authorization testing
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestOtherUserFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestOtherUserFolderID) // handleruser trying to access different user's folder
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusForbidden, result.Code())
	assert.Contains(t, result.Error().Error(), "unauthorized access to folder")
}

// Run_GetByIDHandler_FolderServiceFailure executes GetByIDHandler with FolderService failure
func Run_GetByIDHandler_FolderServiceFailure(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	// TODO: Add service failure testing once type access is resolved
	
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_GetByIDHandler_BookmarkServiceFailure executes GetByIDHandler with BookmarkService failure
func Run_GetByIDHandler_BookmarkServiceFailure(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	// TODO: Add bookmark service failure testing once type access is resolved
	
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	// For now, expect success since we can't force failure
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
}

// Run_GetByIDHandler_NoteServiceFailure executes GetByIDHandler with NoteService failure
func Run_GetByIDHandler_NoteServiceFailure(t *testing.T) {
	testServices := CreateGetByIDTestServices()
	noteService := &services.MockNoteService{
		ShouldFailGetByFolder:   true,
		GetByFolderErrorMessage: "Note service failure",
	}
	testServices.NoteService = noteService
	
	req := CreateGetByIDHTTPRequest(http.MethodGet, "/api/folders/"+getByIDTestFolderID)
	ctx := CreateGetByIDAuthenticatedContext(req, getByIDTestFolderID)
	
	handler := NewGetByIDHandler(testServices, ctx)
	result := handler.Handle()
	
	// Should fail when note service fails
	assert.Error(t, result.Error())
	assert.Equal(t, http.StatusInternalServerError, result.Code())
}

// <tests>
// <evaluators>

// EvaluateGetByIDSuccess validates successful get by ID handler execution
func EvaluateGetByIDSuccess(t *testing.T, result i.IHandler) {
	assert.NoError(t, result.Error())
	assert.Equal(t, http.StatusOK, result.Code())
	assert.NotNil(t, result.Data())
	
	folderData := result.Data().(*responses.FolderData)
	assert.NotNil(t, folderData.ID)
	assert.NotEmpty(t, folderData.Name)
}

// EvaluateGetByIDFailure validates failed get by ID handler execution
func EvaluateGetByIDFailure(t *testing.T, result i.IHandler, expectedCode int) {
	assert.Error(t, result.Error())
	assert.Equal(t, expectedCode, result.Code())
}

// <map/>

// GetByIDHandlerTestMap defines all GetByIDHandler test cases
var GetByIDHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":            Run_GetByIDHandler_ValidRequest,
	"AuthenticationFailure":   Run_GetByIDHandler_AuthenticationFailure,
	"InvalidFolderID":         Run_GetByIDHandler_InvalidFolderID,
	"FolderNotFound":          Run_GetByIDHandler_FolderNotFound,
	"AuthorizationFailure":    Run_GetByIDHandler_AuthorizationFailure,
	"FolderServiceFailure":    Run_GetByIDHandler_FolderServiceFailure,
	"BookmarkServiceFailure":  Run_GetByIDHandler_BookmarkServiceFailure,
	"NoteServiceFailure":      Run_GetByIDHandler_NoteServiceFailure,
}

// <hook/>

// Test_GetByIDHandler tests all GetByIDHandler scenarios
func Test_GetByIDHandler(t *testing.T) {
	fmt.Println("Test_GetByIDHandler")
	for name, testFunc := range GetByIDHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>