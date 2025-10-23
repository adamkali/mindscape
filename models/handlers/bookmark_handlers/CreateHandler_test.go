package bookmark_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// <method var=bookmark_handlers.CreateHandler.Handle>
// <fixtures/>

// Embedded utility functions to avoid package conflicts
func createEchoContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func createJWTToken(userID uuid.UUID, isAdmin bool) *jwt.Token {
	claims := &services.CustomJwt{
		UserId:  userID,
		IsAdmin: isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Raw = "test-token-string"
	return token
}

// CreateBookmarkRequestParams provides base parameters for create bookmark handler tests
func CreateBookmarkRequestParams() map[string]any {
	return map[string]any{
		"user_id":   "11111111-1111-1111-1111-111111111111",
		"folder_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"name":      "Test Bookmark",
		"link":      "https://test.example.com",
	}
}

// WithCreateInvalidJSON modifies request with malformed JSON
func WithCreateInvalidJSON(req map[string]any) string {
	return `{"name": "Test", "link": "invalid json`
}

// WithCreateEmptyName modifies request with empty name
func WithCreateEmptyName(req map[string]any) map[string]any {
	req["name"] = ""
	return req
}

// WithCreateEmptyLink modifies request with empty link
func WithCreateEmptyLink(req map[string]any) map[string]any {
	req["link"] = ""
	return req
}

// WithCreateMissingFields modifies request with missing required fields
func WithCreateMissingFields(req map[string]any) map[string]any {
	delete(req, "name")
	delete(req, "link")
	return req
}

// <runners/>

// Run_CreateBookmarkHandler_ValidRequest executes create handler with valid request
func Run_CreateBookmarkHandler_ValidRequest(t *testing.T) {
	params := CreateBookmarkRequestParams()
	jsonBody, _ := json.Marshal(params)
	ctx, rec := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, 200, result.Code())
	assert.NotNil(t, result.Data())
	
	bookmark := result.Data().(*repository.Bookmark)
	assert.Equal(t, "Test Bookmark", bookmark.Name)
	assert.Equal(t, "https://test.example.com", bookmark.Link)
	assert.Equal(t, userID, bookmark.UserID)
	
	// Verify service calls
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.CreateCallCount)
	
	// Test JSON response
	err := result.JSON()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// Run_CreateBookmarkHandler_InvalidJSON executes create handler with malformed JSON
func Run_CreateBookmarkHandler_InvalidJSON(t *testing.T) {
	invalidJSON := WithCreateInvalidJSON(CreateBookmarkRequestParams())
	ctx, _ := createEchoContext("POST", "/api/bookmarks", invalidJSON)
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 400, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify auth was called but not bookmark service due to validation failure
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 0, bookmarkService.CreateCallCount)
}

// Run_CreateBookmarkHandler_AuthenticationFailure executes create handler with invalid token
func Run_CreateBookmarkHandler_AuthenticationFailure(t *testing.T) {
	params := CreateBookmarkRequestParams()
	jsonBody, _ := json.Marshal(params)
	ctx, _ := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services with auth failure
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	authService.ShouldFailCheckToken = true
	authService.CheckTokenErrorMessage = "Invalid token"
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 401, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify only auth was called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 0, bookmarkService.CreateCallCount)
}

// Run_CreateBookmarkHandler_ValidationFailure executes create handler with invalid data
func Run_CreateBookmarkHandler_ValidationFailure(t *testing.T) {
	params := WithCreateEmptyName(CreateBookmarkRequestParams())
	jsonBody, _ := json.Marshal(params)
	ctx, _ := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 400, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify auth was called but not bookmark service due to validation failure
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 0, bookmarkService.CreateCallCount)
}

// Run_CreateBookmarkHandler_ServiceFailure executes create handler with service failure
func Run_CreateBookmarkHandler_ServiceFailure(t *testing.T) {
	params := CreateBookmarkRequestParams()
	jsonBody, _ := json.Marshal(params)
	ctx, _ := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services with bookmark service failure
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	bookmarkService.ShouldFailCreate = true
	bookmarkService.CreateErrorMessage = "Database error"
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 500, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify both services were called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.CreateCallCount)
}

// Run_CreateBookmarkHandler_MissingJWTToken executes create handler without JWT token
func Run_CreateBookmarkHandler_MissingJWTToken(t *testing.T) {
	params := CreateBookmarkRequestParams()
	jsonBody, _ := json.Marshal(params)
	ctx, _ := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	// Don't set JWT token
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	
	// This should panic due to missing JWT token
	assert.Panics(t, func() {
		handler.Handle()
	})
}

// Run_CreateBookmarkHandler_DuplicateName executes create handler with duplicate bookmark name
func Run_CreateBookmarkHandler_DuplicateName(t *testing.T) {
	params := CreateBookmarkRequestParams()
	params["name"] = "Recent Bookmark" // Existing name in seeded data
	jsonBody, _ := json.Marshal(params)
	ctx, _ := createEchoContext("POST", "/api/bookmarks", string(jsonBody))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	validatorService := services.ValidatorService{}
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewCreateHandler(ctx, validatorService, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 500, result.Code())
	assert.Nil(t, result.Data())
	assert.Contains(t, result.Error().Error(), "bookmark name already exists")
	
	// Verify services were called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.CreateCallCount)
}

// <tests>
// <evaluators>

// EvaluateCreateBookmarkSuccess validates successful bookmark creation
func EvaluateCreateBookmarkSuccess(t *testing.T, handler handlers.IHandler) {
	assert.NoError(t, handler.Error())
	assert.Equal(t, 200, handler.Code())
	assert.NotNil(t, handler.Data())
}

// EvaluateCreateBookmarkFailure validates failed bookmark creation
func EvaluateCreateBookmarkFailure(t *testing.T, handler handlers.IHandler, expectedCode int, expectedError string) {
	assert.Error(t, handler.Error())
	assert.Equal(t, expectedCode, handler.Code())
	assert.Nil(t, handler.Data())
	if expectedError != "" {
		assert.Contains(t, handler.Error().Error(), expectedError)
	}
}

// <map/>

// CreateBookmarkHandlerTestMap defines all CreateHandler test cases
var CreateBookmarkHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":         Run_CreateBookmarkHandler_ValidRequest,
	"InvalidJSON":          Run_CreateBookmarkHandler_InvalidJSON,
	"AuthenticationFailure": Run_CreateBookmarkHandler_AuthenticationFailure,
	"ValidationFailure":    Run_CreateBookmarkHandler_ValidationFailure,
	"ServiceFailure":       Run_CreateBookmarkHandler_ServiceFailure,
	"MissingJWTToken":      Run_CreateBookmarkHandler_MissingJWTToken,
	"DuplicateName":        Run_CreateBookmarkHandler_DuplicateName,
}

// <hook/>

// Test_CreateBookmarkHandler tests all CreateHandler scenarios
func Test_CreateBookmarkHandler(t *testing.T) {
	for name, testFunc := range CreateBookmarkHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>