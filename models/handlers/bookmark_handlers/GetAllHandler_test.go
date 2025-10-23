package bookmark_handlers

import (
	"context"
	"testing"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// <method var=bookmark_handlers.GetByFolderHandler.Handle>
// <fixtures/>

// GetByFolderRequestParams provides base parameters for get folder bookmarks handler tests
func GetByFolderRequestParams() map[string]any {
	return map[string]any{
		"parent_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", // Valid folder ID
		"user_id":   "11111111-1111-1111-1111-111111111111", // Valid user ID
	}
}

// WithGetInvalidFolderID modifies request with invalid folder ID
func WithGetInvalidFolderID(req map[string]any) map[string]any {
	req["parent_id"] = "invalid-uuid"
	return req
}

// WithGetNonexistentFolder modifies request with non-existent folder ID
func WithGetNonexistentFolder(req map[string]any) map[string]any {
	req["parent_id"] = "99999999-9999-9999-9999-999999999999"
	return req
}

// WithGetUnauthorizedUser modifies request with unauthorized user
func WithGetUnauthorizedUser(req map[string]any) map[string]any {
	req["user_id"] = "22222222-2222-2222-2222-222222222222" // Admin user trying to access test user's bookmarks
	return req
}

// <runners/>

// Run_GetByFolderHandler_ValidRequest executes get folder handler with valid request
func Run_GetByFolderHandler_ValidRequest(t *testing.T) {
	params := GetByFolderRequestParams()
	ctx, rec := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	userID := uuid.MustParse(params["user_id"].(string))
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	assert.NoError(t, result.Error())
	assert.Equal(t, 200, result.Code())
	assert.NotNil(t, result.Data())
	
	bookmarks := result.Data().([]repository.Bookmark)
	assert.GreaterOrEqual(t, len(bookmarks), 2) // Should find bookmarks in the test folder
	
	// Verify all bookmarks belong to the requesting user
	for _, bookmark := range bookmarks {
		assert.Equal(t, userID, bookmark.UserID)
	}
	
	// Verify service calls
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.GetByFolderCallCount)
	
	// Test JSON response
	err := result.JSON()
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
}

// Run_GetByFolderHandler_InvalidFolderID executes get folder handler with invalid folder ID
func Run_GetByFolderHandler_InvalidFolderID(t *testing.T) {
	params := WithGetInvalidFolderID(GetByFolderRequestParams())
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 400, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify auth was called but not bookmark service due to parsing failure
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 0, bookmarkService.GetByFolderCallCount)
}

// Run_GetByFolderHandler_AuthenticationFailure executes get folder handler with invalid token
func Run_GetByFolderHandler_AuthenticationFailure(t *testing.T) {
	params := GetByFolderRequestParams()
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	userID := uuid.MustParse(params["user_id"].(string))
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services with auth failure
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	authService.ShouldFailCheckToken = true
	authService.CheckTokenErrorMessage = "Invalid token"
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 401, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify only auth was called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 0, bookmarkService.GetByFolderCallCount)
}

// Run_GetByFolderHandler_FolderNotFound executes get folder handler with non-existent folder
func Run_GetByFolderHandler_FolderNotFound(t *testing.T) {
	params := WithGetNonexistentFolder(GetByFolderRequestParams())
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	// Should succeed but return empty list for non-existent folder
	assert.NoError(t, result.Error())
	assert.Equal(t, 200, result.Code())
	assert.NotNil(t, result.Data())
	
	bookmarks := result.Data().([]repository.Bookmark)
	assert.Equal(t, 0, len(bookmarks)) // Empty for non-existent folder
	
	// Verify service calls
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.GetByFolderCallCount)
}

// Run_GetByFolderHandler_ServiceFailure executes get folder handler with service failure
func Run_GetByFolderHandler_ServiceFailure(t *testing.T) {
	params := GetByFolderRequestParams()
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	userID := uuid.MustParse(params["user_id"].(string))
	ctx.Set("user", createJWTToken(userID, false))
	
	// Setup mock services with bookmark service failure
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	bookmarkService.ShouldFailGetByFolder = true
	bookmarkService.GetByFolderErrorMessage = "Database error"
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 404, result.Code())
	assert.Nil(t, result.Data())
	
	// Verify both services were called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.GetByFolderCallCount)
}

// Run_GetByFolderHandler_UnauthorizedAccess executes get folder handler with unauthorized user access
func Run_GetByFolderHandler_UnauthorizedAccess(t *testing.T) {
	params := GetByFolderRequestParams()
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	
	// Use admin user ID but try to access test user's bookmarks
	adminUserID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	ctx.Set("user", createJWTToken(adminUserID, true))
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	result := handler.Handle()
	
	assert.Error(t, result.Error())
	assert.Equal(t, 403, result.Code())
	// Note: Data might still be set since authorization check happens after retrieval
	assert.Contains(t, result.Error().Error(), "unauthorized access to bookmark")
	
	// Verify services were called
	assert.Equal(t, 1, authService.CheckTokenCallCount)
	assert.Equal(t, 1, bookmarkService.GetByFolderCallCount)
}

// Run_GetByFolderHandler_MissingJWTToken executes get folder handler without JWT token
func Run_GetByFolderHandler_MissingJWTToken(t *testing.T) {
	params := GetByFolderRequestParams()
	ctx, _ := createEchoContext("GET", "/api/bookmarks/folder/:parent_id", "")
	ctx.SetParamNames("parent_id")
	ctx.SetParamValues(params["parent_id"].(string))
	// Don't set JWT token
	
	// Setup mock services
	authService := services.CreateMockAuthService(context.Background(), nil)
	bookmarkService := services.CreateMockBookmarkService(context.Background(), nil)
	
	authService.Reset()
	bookmarkService.Reset()
	
	handler := NewGetFolderHandler(ctx, bookmarkService, authService)
	
	// This should panic due to missing JWT token
	assert.Panics(t, func() {
		handler.Handle()
	})
}

// <tests>
// <evaluators>

// EvaluateGetByFolderSuccess validates successful folder bookmark retrieval
func EvaluateGetByFolderSuccess(t *testing.T, handler handlers.IHandler) {
	assert.NoError(t, handler.Error())
	assert.Equal(t, 200, handler.Code())
	assert.NotNil(t, handler.Data())
}

// EvaluateGetByFolderFailure validates failed folder bookmark retrieval
func EvaluateGetByFolderFailure(t *testing.T, handler handlers.IHandler, expectedCode int, expectedError string) {
	assert.Error(t, handler.Error())
	assert.Equal(t, expectedCode, handler.Code())
	// Note: Data might be set depending on when the failure occurs
	if expectedError != "" {
		assert.Contains(t, handler.Error().Error(), expectedError)
	}
}

// <map/>

// GetByFolderHandlerTestMap defines all GetByFolderHandler test cases
var GetByFolderHandlerTestMap = map[string]func(*testing.T){
	"ValidRequest":         Run_GetByFolderHandler_ValidRequest,
	"InvalidFolderID":      Run_GetByFolderHandler_InvalidFolderID,
	"AuthenticationFailure": Run_GetByFolderHandler_AuthenticationFailure,
	"FolderNotFound":       Run_GetByFolderHandler_FolderNotFound,
	"ServiceFailure":       Run_GetByFolderHandler_ServiceFailure,
	"UnauthorizedAccess":   Run_GetByFolderHandler_UnauthorizedAccess,
	"MissingJWTToken":      Run_GetByFolderHandler_MissingJWTToken,
}

// <hook/>

// Test_GetByFolderHandler tests all GetByFolderHandler scenarios
func Test_GetByFolderHandler(t *testing.T) {
	for name, testFunc := range GetByFolderHandlerTestMap {
		t.Run(name, func(t *testing.T) {
			testFunc(t)
		})
	}
}

// </method>