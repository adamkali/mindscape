package services

import (
	"context"
	"testing"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// <method var=services.MockBookmarkService.GetAll>
// <fixtures/>

// <runners/>

// Run_BookmarkGetAll_Success executes GetAll successfully
func Run_BookmarkGetAll_Success(t *testing.T, service *MockBookmarkService) {
	bookmarks, err := service.GetAll()
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.GreaterOrEqual(t, len(bookmarks), 4) // At least 4 seeded bookmarks
	assert.Equal(t, 1, service.GetAllCallCount)
}

// Run_BookmarkGetAll_ForceFailure executes GetAll with forced failure
func Run_BookmarkGetAll_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetAll = true
	service.GetAllErrorMessage = "Forced GetAll failure"
	
	bookmarks, err := service.GetAll()
	
	assert.Error(t, err)
	assert.Nil(t, bookmarks)
	assert.Contains(t, err.Error(), "Forced GetAll failure")
	assert.Equal(t, 1, service.GetAllCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetAllSuccess validates successful bookmark retrieval
func EvaluateBookmarkGetAllSuccess(t *testing.T, service *MockBookmarkService, expectedCount int) {
	assert.Equal(t, 1, service.GetAllCallCount)
}

// EvaluateBookmarkGetAllFailure validates failed bookmark retrieval
func EvaluateBookmarkGetAllFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetAllCallCount)
}

// <map/>

// BookmarkGetAllTestMap defines all GetAll method test cases
var BookmarkGetAllTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"Success":      Run_BookmarkGetAll_Success,
	"ForceFailure": Run_BookmarkGetAll_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetAll tests all GetAll method scenarios
func Test_MockBookmarkService_GetAll(t *testing.T) {
	for name, testFunc := range BookmarkGetAllTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.Get>
// <fixtures/>

// ValidBookmarkID provides a valid bookmark ID for testing
func ValidBookmarkID() uuid.UUID {
	return uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd") // From seeded data
}

// InvalidBookmarkID provides an invalid bookmark ID for testing
func InvalidBookmarkID() uuid.UUID {
	return uuid.MustParse("99999999-9999-9999-9999-999999999999")
}

// <runners/>

// Run_BookmarkGet_ValidID executes get with valid bookmark ID
func Run_BookmarkGet_ValidID(t *testing.T, service *MockBookmarkService) {
	bookmarkID := ValidBookmarkID()
	bookmark, err := service.Get(bookmarkID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, bookmarkID, bookmark.ID)
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, bookmarkID, service.LastGetID)
}

// Run_BookmarkGet_InvalidID executes get with invalid bookmark ID
func Run_BookmarkGet_InvalidID(t *testing.T, service *MockBookmarkService) {
	bookmarkID := InvalidBookmarkID()
	bookmark, err := service.Get(bookmarkID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark not found")
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, bookmarkID, service.LastGetID)
}

// Run_BookmarkGet_ForceFailure executes get with forced failure
func Run_BookmarkGet_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGet = true
	service.GetErrorMessage = "Forced Get failure"
	bookmarkID := ValidBookmarkID()
	
	bookmark, err := service.Get(bookmarkID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "Forced Get failure")
	assert.Equal(t, 1, service.GetCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetSuccess validates successful bookmark retrieval
func EvaluateBookmarkGetSuccess(t *testing.T, service *MockBookmarkService, expectedID uuid.UUID) {
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, expectedID, service.LastGetID)
}

// EvaluateBookmarkGetFailure validates failed bookmark retrieval
func EvaluateBookmarkGetFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetCallCount)
}

// <map/>

// BookmarkGetTestMap defines all Get method test cases
var BookmarkGetTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidID":      Run_BookmarkGet_ValidID,
	"InvalidID":    Run_BookmarkGet_InvalidID,
	"ForceFailure": Run_BookmarkGet_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_Get tests all Get method scenarios
func Test_MockBookmarkService_Get(t *testing.T) {
	for name, testFunc := range BookmarkGetTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.GetByFolder>
// <fixtures/>

// ValidBookmarkFolderID provides a valid folder ID for testing
func ValidBookmarkFolderID() uuid.UUID {
	return uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa") // From seeded data
}

// InvalidBookmarkFolderID provides an invalid folder ID for testing
func InvalidBookmarkFolderID() uuid.UUID {
	return uuid.MustParse("99999999-9999-9999-9999-999999999999")
}

// <runners/>

// Run_BookmarkGetByFolder_ValidFolder executes GetByFolder with valid folder ID
func Run_BookmarkGetByFolder_ValidFolder(t *testing.T, service *MockBookmarkService) {
	folderID := ValidBookmarkFolderID()
	bookmarks, err := service.GetByFolder(folderID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.GreaterOrEqual(t, len(bookmarks), 2) // At least 2 bookmarks in test folder
	assert.Equal(t, 1, service.GetByFolderCallCount)
	assert.Equal(t, folderID, service.LastGetByFolderID)
}

// Run_BookmarkGetByFolder_InvalidFolder executes GetByFolder with invalid folder ID
func Run_BookmarkGetByFolder_InvalidFolder(t *testing.T, service *MockBookmarkService) {
	folderID := InvalidBookmarkFolderID()
	bookmarks, err := service.GetByFolder(folderID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.Equal(t, 0, len(bookmarks)) // No bookmarks for invalid folder
	assert.Equal(t, 1, service.GetByFolderCallCount)
	assert.Equal(t, folderID, service.LastGetByFolderID)
}

// Run_BookmarkGetByFolder_ForceFailure executes GetByFolder with forced failure
func Run_BookmarkGetByFolder_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetByFolder = true
	service.GetByFolderErrorMessage = "Forced GetByFolder failure"
	folderID := ValidBookmarkFolderID()
	
	bookmarks, err := service.GetByFolder(folderID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmarks)
	assert.Contains(t, err.Error(), "Forced GetByFolder failure")
	assert.Equal(t, 1, service.GetByFolderCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetByFolderSuccess validates successful folder bookmark retrieval
func EvaluateBookmarkGetByFolderSuccess(t *testing.T, service *MockBookmarkService, expectedFolderID uuid.UUID) {
	assert.Equal(t, 1, service.GetByFolderCallCount)
	assert.Equal(t, expectedFolderID, service.LastGetByFolderID)
}

// EvaluateBookmarkGetByFolderFailure validates failed folder bookmark retrieval
func EvaluateBookmarkGetByFolderFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetByFolderCallCount)
}

// <map/>

// BookmarkGetByFolderTestMap defines all GetByFolder method test cases
var BookmarkGetByFolderTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidFolder":   Run_BookmarkGetByFolder_ValidFolder,
	"InvalidFolder": Run_BookmarkGetByFolder_InvalidFolder,
	"ForceFailure":  Run_BookmarkGetByFolder_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetByFolder tests all GetByFolder method scenarios
func Test_MockBookmarkService_GetByFolder(t *testing.T) {
	for name, testFunc := range BookmarkGetByFolderTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.GetByUser>
// <fixtures/>

// ValidBookmarkUserID provides a valid user ID for testing
func ValidBookmarkUserID() uuid.UUID {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111") // From seeded data
}

// InvalidBookmarkUserID provides an invalid user ID for testing
func InvalidBookmarkUserID() uuid.UUID {
	return uuid.MustParse("99999999-9999-9999-9999-999999999999")
}

// <runners/>

// Run_BookmarkGetByUser_ValidUser executes GetByUser with valid user ID
func Run_BookmarkGetByUser_ValidUser(t *testing.T, service *MockBookmarkService) {
	userID := ValidBookmarkUserID()
	bookmarks, err := service.GetByUser(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.GreaterOrEqual(t, len(bookmarks), 3) // At least 3 bookmarks for test user
	assert.Equal(t, 1, service.GetByUserCallCount)
	assert.Equal(t, userID, service.LastGetByUserID)
}

// Run_BookmarkGetByUser_InvalidUser executes GetByUser with invalid user ID
func Run_BookmarkGetByUser_InvalidUser(t *testing.T, service *MockBookmarkService) {
	userID := InvalidBookmarkUserID()
	bookmarks, err := service.GetByUser(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.Equal(t, 0, len(bookmarks)) // No bookmarks for invalid user
	assert.Equal(t, 1, service.GetByUserCallCount)
	assert.Equal(t, userID, service.LastGetByUserID)
}

// Run_BookmarkGetByUser_ForceFailure executes GetByUser with forced failure
func Run_BookmarkGetByUser_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetByUser = true
	service.GetByUserErrorMessage = "Forced GetByUser failure"
	userID := ValidBookmarkUserID()
	
	bookmarks, err := service.GetByUser(userID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmarks)
	assert.Contains(t, err.Error(), "Forced GetByUser failure")
	assert.Equal(t, 1, service.GetByUserCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetByUserSuccess validates successful user bookmark retrieval
func EvaluateBookmarkGetByUserSuccess(t *testing.T, service *MockBookmarkService, expectedUserID uuid.UUID) {
	assert.Equal(t, 1, service.GetByUserCallCount)
	assert.Equal(t, expectedUserID, service.LastGetByUserID)
}

// EvaluateBookmarkGetByUserFailure validates failed user bookmark retrieval
func EvaluateBookmarkGetByUserFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetByUserCallCount)
}

// <map/>

// BookmarkGetByUserTestMap defines all GetByUser method test cases
var BookmarkGetByUserTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidUser":    Run_BookmarkGetByUser_ValidUser,
	"InvalidUser":  Run_BookmarkGetByUser_InvalidUser,
	"ForceFailure": Run_BookmarkGetByUser_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetByUser tests all GetByUser method scenarios
func Test_MockBookmarkService_GetByUser(t *testing.T) {
	for name, testFunc := range BookmarkGetByUserTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.GetMostRecent>
// <fixtures/>

// <runners/>

// Run_BookmarkGetMostRecent_ValidUser executes GetMostRecent with valid user ID
func Run_BookmarkGetMostRecent_ValidUser(t *testing.T, service *MockBookmarkService) {
	userID := ValidBookmarkUserID()
	bookmark, err := service.GetMostRecent(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, userID, bookmark.UserID)
	assert.Equal(t, "Recent Bookmark", bookmark.Name) // Should be the most recent
	assert.Equal(t, 1, service.GetMostRecentCallCount)
	assert.Equal(t, userID, service.LastGetMostRecentUserID)
}

// Run_BookmarkGetMostRecent_InvalidUser executes GetMostRecent with invalid user ID
func Run_BookmarkGetMostRecent_InvalidUser(t *testing.T, service *MockBookmarkService) {
	userID := InvalidBookmarkUserID()
	bookmark, err := service.GetMostRecent(userID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "no bookmarks found for user")
	assert.Equal(t, 1, service.GetMostRecentCallCount)
	assert.Equal(t, userID, service.LastGetMostRecentUserID)
}

// Run_BookmarkGetMostRecent_ForceFailure executes GetMostRecent with forced failure
func Run_BookmarkGetMostRecent_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetMostRecent = true
	service.GetMostRecentErrorMessage = "Forced GetMostRecent failure"
	userID := ValidBookmarkUserID()
	
	bookmark, err := service.GetMostRecent(userID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "Forced GetMostRecent failure")
	assert.Equal(t, 1, service.GetMostRecentCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetMostRecentSuccess validates successful most recent bookmark retrieval
func EvaluateBookmarkGetMostRecentSuccess(t *testing.T, service *MockBookmarkService, expectedUserID uuid.UUID) {
	assert.Equal(t, 1, service.GetMostRecentCallCount)
	assert.Equal(t, expectedUserID, service.LastGetMostRecentUserID)
}

// EvaluateBookmarkGetMostRecentFailure validates failed most recent bookmark retrieval
func EvaluateBookmarkGetMostRecentFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetMostRecentCallCount)
}

// <map/>

// BookmarkGetMostRecentTestMap defines all GetMostRecent method test cases
var BookmarkGetMostRecentTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidUser":    Run_BookmarkGetMostRecent_ValidUser,
	"InvalidUser":  Run_BookmarkGetMostRecent_InvalidUser,
	"ForceFailure": Run_BookmarkGetMostRecent_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetMostRecent tests all GetMostRecent method scenarios
func Test_MockBookmarkService_GetMostRecent(t *testing.T) {
	for name, testFunc := range BookmarkGetMostRecentTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.GetMostRecents>
// <fixtures/>

// <runners/>

// Run_BookmarkGetMostRecents_ValidUser executes GetMostRecents with valid user ID
func Run_BookmarkGetMostRecents_ValidUser(t *testing.T, service *MockBookmarkService) {
	userID := ValidBookmarkUserID()
	bookmarks, err := service.GetMostRecents(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.GreaterOrEqual(t, len(bookmarks), 3) // At least 3 bookmarks for test user
	assert.Equal(t, 1, service.GetMostRecentsCallCount)
	assert.Equal(t, userID, service.LastGetMostRecentsUserID)
	
	// Verify ordering (most recent first)
	if len(bookmarks) >= 2 {
		assert.Equal(t, "Recent Bookmark", bookmarks[0].Name)
	}
}

// Run_BookmarkGetMostRecents_InvalidUser executes GetMostRecents with invalid user ID
func Run_BookmarkGetMostRecents_InvalidUser(t *testing.T, service *MockBookmarkService) {
	userID := InvalidBookmarkUserID()
	bookmarks, err := service.GetMostRecents(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.Equal(t, 0, len(bookmarks)) // No bookmarks for invalid user
	assert.Equal(t, 1, service.GetMostRecentsCallCount)
	assert.Equal(t, userID, service.LastGetMostRecentsUserID)
}

// Run_BookmarkGetMostRecents_ForceFailure executes GetMostRecents with forced failure
func Run_BookmarkGetMostRecents_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetMostRecents = true
	service.GetMostRecentsErrorMessage = "Forced GetMostRecents failure"
	userID := ValidBookmarkUserID()
	
	bookmarks, err := service.GetMostRecents(userID)
	
	assert.Error(t, err)
	assert.Nil(t, bookmarks)
	assert.Contains(t, err.Error(), "Forced GetMostRecents failure")
	assert.Equal(t, 1, service.GetMostRecentsCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetMostRecentsSuccess validates successful most recent bookmarks retrieval
func EvaluateBookmarkGetMostRecentsSuccess(t *testing.T, service *MockBookmarkService, expectedUserID uuid.UUID) {
	assert.Equal(t, 1, service.GetMostRecentsCallCount)
	assert.Equal(t, expectedUserID, service.LastGetMostRecentsUserID)
}

// EvaluateBookmarkGetMostRecentsFailure validates failed most recent bookmarks retrieval
func EvaluateBookmarkGetMostRecentsFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetMostRecentsCallCount)
}

// <map/>

// BookmarkGetMostRecentsTestMap defines all GetMostRecents method test cases
var BookmarkGetMostRecentsTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidUser":    Run_BookmarkGetMostRecents_ValidUser,
	"InvalidUser":  Run_BookmarkGetMostRecents_InvalidUser,
	"ForceFailure": Run_BookmarkGetMostRecents_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetMostRecents tests all GetMostRecents method scenarios
func Test_MockBookmarkService_GetMostRecents(t *testing.T) {
	for name, testFunc := range BookmarkGetMostRecentsTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.GetByDateRange>
// <fixtures/>

// DateRangeParams provides base parameters for date range tests
func DateRangeParams() *repository.FindBookmarksByUserIDDateTimeRangeParams {
	now := time.Now()
	startTime := now.Add(-2 * time.Hour)
	endTime := now.Add(1 * time.Hour)
	
	return &repository.FindBookmarksByUserIDDateTimeRangeParams{
		UserID:            ValidBookmarkUserID(),
		UpdatedDatetime:   &startTime,
		UpdatedDatetime_2: &endTime,
	}
}

// WithNarrowDateRange modifies date range to exclude some bookmarks
func WithNarrowDateRange(params *repository.FindBookmarksByUserIDDateTimeRangeParams) *repository.FindBookmarksByUserIDDateTimeRangeParams {
	now := time.Now()
	startTime := now.Add(-30 * time.Minute)
	endTime := now.Add(30 * time.Minute)
	
	params.UpdatedDatetime = &startTime
	params.UpdatedDatetime_2 = &endTime
	return params
}

// WithInvalidDateRangeUser modifies date range with invalid user
func WithInvalidDateRangeUser(params *repository.FindBookmarksByUserIDDateTimeRangeParams) *repository.FindBookmarksByUserIDDateTimeRangeParams {
	params.UserID = InvalidBookmarkUserID()
	return params
}

// <runners/>

// Run_BookmarkGetByDateRange_ValidParams executes GetByDateRange with valid parameters
func Run_BookmarkGetByDateRange_ValidParams(t *testing.T, service *MockBookmarkService) {
	params := DateRangeParams()
	bookmarks, err := service.GetByDateRange(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.GreaterOrEqual(t, len(bookmarks), 1) // Should find at least recent bookmark
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
	assert.Equal(t, params, service.LastGetByDateRangeParams)
}

// Run_BookmarkGetByDateRange_NarrowRange executes GetByDateRange with narrow date range
func Run_BookmarkGetByDateRange_NarrowRange(t *testing.T, service *MockBookmarkService) {
	params := WithNarrowDateRange(DateRangeParams())
	bookmarks, err := service.GetByDateRange(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	// Should find recent bookmark but not older ones
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
	assert.Equal(t, params, service.LastGetByDateRangeParams)
}

// Run_BookmarkGetByDateRange_InvalidUser executes GetByDateRange with invalid user
func Run_BookmarkGetByDateRange_InvalidUser(t *testing.T, service *MockBookmarkService) {
	params := WithInvalidDateRangeUser(DateRangeParams())
	bookmarks, err := service.GetByDateRange(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.Equal(t, 0, len(bookmarks)) // No bookmarks for invalid user
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
	assert.Equal(t, params, service.LastGetByDateRangeParams)
}

// Run_BookmarkGetByDateRange_NilParams executes GetByDateRange with nil parameters
func Run_BookmarkGetByDateRange_NilParams(t *testing.T, service *MockBookmarkService) {
	bookmarks, err := service.GetByDateRange(nil)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmarks)
	assert.Equal(t, 0, len(bookmarks)) // Empty for nil params
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
}

// Run_BookmarkGetByDateRange_ForceFailure executes GetByDateRange with forced failure
func Run_BookmarkGetByDateRange_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailGetByDateRange = true
	service.GetByDateRangeErrorMessage = "Forced GetByDateRange failure"
	params := DateRangeParams()
	
	bookmarks, err := service.GetByDateRange(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmarks)
	assert.Contains(t, err.Error(), "Forced GetByDateRange failure")
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkGetByDateRangeSuccess validates successful date range bookmark retrieval
func EvaluateBookmarkGetByDateRangeSuccess(t *testing.T, service *MockBookmarkService, expectedParams *repository.FindBookmarksByUserIDDateTimeRangeParams) {
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
	assert.Equal(t, expectedParams, service.LastGetByDateRangeParams)
}

// EvaluateBookmarkGetByDateRangeFailure validates failed date range bookmark retrieval
func EvaluateBookmarkGetByDateRangeFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.GetByDateRangeCallCount)
}

// <map/>

// BookmarkGetByDateRangeTestMap defines all GetByDateRange method test cases
var BookmarkGetByDateRangeTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidParams":  Run_BookmarkGetByDateRange_ValidParams,
	"NarrowRange":  Run_BookmarkGetByDateRange_NarrowRange,
	"InvalidUser":  Run_BookmarkGetByDateRange_InvalidUser,
	"NilParams":    Run_BookmarkGetByDateRange_NilParams,
	"ForceFailure": Run_BookmarkGetByDateRange_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_GetByDateRange tests all GetByDateRange method scenarios
func Test_MockBookmarkService_GetByDateRange(t *testing.T) {
	for name, testFunc := range BookmarkGetByDateRangeTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.Create>
// <fixtures/>

// CreateBookmarkParams provides base parameters for bookmark creation tests
func CreateBookmarkParams() *repository.CreateBookmarkParams {
	return &repository.CreateBookmarkParams{
		UserID:   ValidBookmarkUserID(),
		FolderID: ValidBookmarkFolderID(),
		Name:     "New Test Bookmark",
		Link:     "https://new-test.example.com",
	}
}

// WithCreateEmptyName modifies request with empty name
func WithCreateEmptyName(params *repository.CreateBookmarkParams) *repository.CreateBookmarkParams {
	params.Name = ""
	return params
}

// WithCreateEmptyLink modifies request with empty link
func WithCreateEmptyLink(params *repository.CreateBookmarkParams) *repository.CreateBookmarkParams {
	params.Link = ""
	return params
}

// WithCreateDuplicateName modifies request with existing name
func WithCreateDuplicateName(params *repository.CreateBookmarkParams) *repository.CreateBookmarkParams {
	params.Name = "Recent Bookmark" // This name exists in seeded data
	return params
}

// WithCreateInvalidUser modifies request with invalid user ID
func WithCreateInvalidUser(params *repository.CreateBookmarkParams) *repository.CreateBookmarkParams {
	params.UserID = InvalidBookmarkUserID()
	return params
}

// <runners/>

// Run_BookmarkCreate_ValidParams executes create with valid parameters
func Run_BookmarkCreate_ValidParams(t *testing.T, service *MockBookmarkService) {
	params := CreateBookmarkParams()
	initialCount := service.GetBookmarkCount()
	
	bookmark, err := service.Create(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, params.Name, bookmark.Name)
	assert.Equal(t, params.Link, bookmark.Link)
	assert.Equal(t, params.UserID, bookmark.UserID)
	assert.Equal(t, params.FolderID, bookmark.FolderID)
	assert.Equal(t, initialCount+1, service.GetBookmarkCount())
	assert.Equal(t, 1, service.CreateCallCount)
	assert.Equal(t, params, service.LastCreateParams)
}

// Run_BookmarkCreate_EmptyName executes create with empty name
func Run_BookmarkCreate_EmptyName(t *testing.T, service *MockBookmarkService) {
	params := WithCreateEmptyName(CreateBookmarkParams())
	initialCount := service.GetBookmarkCount()
	
	bookmark, err := service.Create(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark name cannot be empty")
	assert.Equal(t, initialCount, service.GetBookmarkCount())
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_BookmarkCreate_EmptyLink executes create with empty link
func Run_BookmarkCreate_EmptyLink(t *testing.T, service *MockBookmarkService) {
	params := WithCreateEmptyLink(CreateBookmarkParams())
	initialCount := service.GetBookmarkCount()
	
	bookmark, err := service.Create(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark link cannot be empty")
	assert.Equal(t, initialCount, service.GetBookmarkCount())
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_BookmarkCreate_DuplicateName executes create with duplicate name
func Run_BookmarkCreate_DuplicateName(t *testing.T, service *MockBookmarkService) {
	params := WithCreateDuplicateName(CreateBookmarkParams())
	initialCount := service.GetBookmarkCount()
	
	bookmark, err := service.Create(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark name already exists")
	assert.Equal(t, initialCount, service.GetBookmarkCount())
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_BookmarkCreate_NilParams executes create with nil parameters
func Run_BookmarkCreate_NilParams(t *testing.T, service *MockBookmarkService) {
	bookmark, err := service.Create(nil)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "create parameters cannot be nil")
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_BookmarkCreate_ForceFailure executes create with forced failure
func Run_BookmarkCreate_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailCreate = true
	service.CreateErrorMessage = "Forced Create failure"
	params := CreateBookmarkParams()
	
	bookmark, err := service.Create(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "Forced Create failure")
	assert.Equal(t, 1, service.CreateCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkCreateSuccess validates successful bookmark creation
func EvaluateBookmarkCreateSuccess(t *testing.T, service *MockBookmarkService, expectedParams *repository.CreateBookmarkParams) {
	assert.Equal(t, 1, service.CreateCallCount)
	assert.Equal(t, expectedParams, service.LastCreateParams)
}

// EvaluateBookmarkCreateFailure validates failed bookmark creation
func EvaluateBookmarkCreateFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.CreateCallCount)
}

// <map/>

// BookmarkCreateTestMap defines all Create method test cases
var BookmarkCreateTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidParams":   Run_BookmarkCreate_ValidParams,
	"EmptyName":     Run_BookmarkCreate_EmptyName,
	"EmptyLink":     Run_BookmarkCreate_EmptyLink,
	"DuplicateName": Run_BookmarkCreate_DuplicateName,
	"NilParams":     Run_BookmarkCreate_NilParams,
	"ForceFailure":  Run_BookmarkCreate_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_Create tests all Create method scenarios
func Test_MockBookmarkService_Create(t *testing.T) {
	for name, testFunc := range BookmarkCreateTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.Update>
// <fixtures/>

// UpdateBookmarkParams provides base parameters for bookmark update tests
func UpdateBookmarkParams() *repository.UpdateBookmarkParams {
	return &repository.UpdateBookmarkParams{
		ID:       ValidBookmarkID(),
		FolderID: ValidBookmarkFolderID(),
		Name:     "Updated Bookmark Name",
		Link:     "https://updated.example.com",
	}
}

// WithUpdateEmptyName modifies update request with empty name
func WithUpdateEmptyName(params *repository.UpdateBookmarkParams) *repository.UpdateBookmarkParams {
	params.Name = ""
	return params
}

// WithUpdateEmptyLink modifies update request with empty link
func WithUpdateEmptyLink(params *repository.UpdateBookmarkParams) *repository.UpdateBookmarkParams {
	params.Link = ""
	return params
}

// WithUpdateInvalidID modifies update request with invalid ID
func WithUpdateInvalidID(params *repository.UpdateBookmarkParams) *repository.UpdateBookmarkParams {
	params.ID = InvalidBookmarkID()
	return params
}

// WithUpdateDifferentFolder modifies update request to change folder
func WithUpdateDifferentFolder(params *repository.UpdateBookmarkParams) *repository.UpdateBookmarkParams {
	params.FolderID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	return params
}

// <runners/>

// Run_BookmarkUpdate_ValidParams executes update with valid parameters
func Run_BookmarkUpdate_ValidParams(t *testing.T, service *MockBookmarkService) {
	params := UpdateBookmarkParams()
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, params.Name, bookmark.Name)
	assert.Equal(t, params.Link, bookmark.Link)
	assert.Equal(t, params.FolderID, bookmark.FolderID)
	assert.Equal(t, params.ID, bookmark.ID)
	assert.Equal(t, 1, service.UpdateCallCount)
	assert.Equal(t, params.ID, service.LastUpdateID)
	assert.Equal(t, params, service.LastUpdateParams)
}

// Run_BookmarkUpdate_FolderChange executes update with folder change
func Run_BookmarkUpdate_FolderChange(t *testing.T, service *MockBookmarkService) {
	params := WithUpdateDifferentFolder(UpdateBookmarkParams())
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, params.FolderID, bookmark.FolderID)
	assert.Equal(t, 1, service.UpdateCallCount)
}

// Run_BookmarkUpdate_EmptyName executes update with empty name
func Run_BookmarkUpdate_EmptyName(t *testing.T, service *MockBookmarkService) {
	params := WithUpdateEmptyName(UpdateBookmarkParams())
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark name cannot be empty")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// Run_BookmarkUpdate_EmptyLink executes update with empty link
func Run_BookmarkUpdate_EmptyLink(t *testing.T, service *MockBookmarkService) {
	params := WithUpdateEmptyLink(UpdateBookmarkParams())
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark link cannot be empty")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// Run_BookmarkUpdate_InvalidID executes update with invalid bookmark ID
func Run_BookmarkUpdate_InvalidID(t *testing.T, service *MockBookmarkService) {
	params := WithUpdateInvalidID(UpdateBookmarkParams())
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark not found")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// Run_BookmarkUpdate_NilParams executes update with nil parameters
func Run_BookmarkUpdate_NilParams(t *testing.T, service *MockBookmarkService) {
	id := ValidBookmarkID()
	bookmark, err := service.Update(id, nil)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "update parameters cannot be nil")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// Run_BookmarkUpdate_ForceFailure executes update with forced failure
func Run_BookmarkUpdate_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailUpdate = true
	service.UpdateErrorMessage = "Forced Update failure"
	params := UpdateBookmarkParams()
	
	bookmark, err := service.Update(params.ID, params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "Forced Update failure")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkUpdateSuccess validates successful bookmark update
func EvaluateBookmarkUpdateSuccess(t *testing.T, service *MockBookmarkService, expectedID uuid.UUID, expectedParams *repository.UpdateBookmarkParams) {
	assert.Equal(t, 1, service.UpdateCallCount)
	assert.Equal(t, expectedID, service.LastUpdateID)
	assert.Equal(t, expectedParams, service.LastUpdateParams)
}

// EvaluateBookmarkUpdateFailure validates failed bookmark update
func EvaluateBookmarkUpdateFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.UpdateCallCount)
}

// <map/>

// BookmarkUpdateTestMap defines all Update method test cases
var BookmarkUpdateTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidParams":   Run_BookmarkUpdate_ValidParams,
	"FolderChange":  Run_BookmarkUpdate_FolderChange,
	"EmptyName":     Run_BookmarkUpdate_EmptyName,
	"EmptyLink":     Run_BookmarkUpdate_EmptyLink,
	"InvalidID":     Run_BookmarkUpdate_InvalidID,
	"NilParams":     Run_BookmarkUpdate_NilParams,
	"ForceFailure":  Run_BookmarkUpdate_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_Update tests all Update method scenarios
func Test_MockBookmarkService_Update(t *testing.T) {
	for name, testFunc := range BookmarkUpdateTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.Move>
// <fixtures/>

// MoveBookmarkParams provides base parameters for bookmark move tests
func MoveBookmarkParams() *repository.MoveBookmarkParams {
	return &repository.MoveBookmarkParams{
		ID:       ValidBookmarkID(),
		FolderID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), // Different folder
	}
}

// WithMoveInvalidID modifies move request with invalid ID
func WithMoveInvalidID(params *repository.MoveBookmarkParams) *repository.MoveBookmarkParams {
	params.ID = InvalidBookmarkID()
	return params
}

// WithMoveSameFolder modifies move request to same folder (no change)
func WithMoveSameFolder(params *repository.MoveBookmarkParams) *repository.MoveBookmarkParams {
	params.FolderID = ValidBookmarkFolderID() // Same as current folder
	return params
}

// <runners/>

// Run_BookmarkMove_ValidParams executes move with valid parameters
func Run_BookmarkMove_ValidParams(t *testing.T, service *MockBookmarkService) {
	params := MoveBookmarkParams()
	originalBookmark, _ := service.Get(params.ID)
	originalFolderID := originalBookmark.FolderID // Store the original folder ID
	service.GetCallCount = 0 // Reset the call count after the initial get
	
	bookmark, err := service.Move(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, params.FolderID, bookmark.FolderID)
	assert.Equal(t, params.ID, bookmark.ID)
	assert.NotEqual(t, originalFolderID, bookmark.FolderID) // Folder changed
	assert.Equal(t, 1, service.MoveCallCount)
	assert.Equal(t, params, service.LastMoveParams)
}

// Run_BookmarkMove_SameFolder executes move to same folder
func Run_BookmarkMove_SameFolder(t *testing.T, service *MockBookmarkService) {
	params := WithMoveSameFolder(MoveBookmarkParams())
	
	bookmark, err := service.Move(params)
	
	assert.NoError(t, err)
	assert.NotNil(t, bookmark)
	assert.Equal(t, params.FolderID, bookmark.FolderID)
	assert.Equal(t, 1, service.MoveCallCount)
}

// Run_BookmarkMove_InvalidID executes move with invalid bookmark ID
func Run_BookmarkMove_InvalidID(t *testing.T, service *MockBookmarkService) {
	params := WithMoveInvalidID(MoveBookmarkParams())
	
	bookmark, err := service.Move(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "bookmark not found")
	assert.Equal(t, 1, service.MoveCallCount)
}

// Run_BookmarkMove_NilParams executes move with nil parameters
func Run_BookmarkMove_NilParams(t *testing.T, service *MockBookmarkService) {
	bookmark, err := service.Move(nil)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "move parameters cannot be nil")
	assert.Equal(t, 1, service.MoveCallCount)
}

// Run_BookmarkMove_ForceFailure executes move with forced failure
func Run_BookmarkMove_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailMove = true
	service.MoveErrorMessage = "Forced Move failure"
	params := MoveBookmarkParams()
	
	bookmark, err := service.Move(params)
	
	assert.Error(t, err)
	assert.Nil(t, bookmark)
	assert.Contains(t, err.Error(), "Forced Move failure")
	assert.Equal(t, 1, service.MoveCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkMoveSuccess validates successful bookmark move
func EvaluateBookmarkMoveSuccess(t *testing.T, service *MockBookmarkService, expectedParams *repository.MoveBookmarkParams) {
	assert.Equal(t, 1, service.MoveCallCount)
	assert.Equal(t, expectedParams, service.LastMoveParams)
}

// EvaluateBookmarkMoveFailure validates failed bookmark move
func EvaluateBookmarkMoveFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.MoveCallCount)
}

// <map/>

// BookmarkMoveTestMap defines all Move method test cases
var BookmarkMoveTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidParams":  Run_BookmarkMove_ValidParams,
	"SameFolder":   Run_BookmarkMove_SameFolder,
	"InvalidID":    Run_BookmarkMove_InvalidID,
	"NilParams":    Run_BookmarkMove_NilParams,
	"ForceFailure": Run_BookmarkMove_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_Move tests all Move method scenarios
func Test_MockBookmarkService_Move(t *testing.T) {
	for name, testFunc := range BookmarkMoveTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockBookmarkService.Remove>
// <fixtures/>

// <runners/>

// Run_BookmarkRemove_ValidID executes remove with valid bookmark ID
func Run_BookmarkRemove_ValidID(t *testing.T, service *MockBookmarkService) {
	bookmarkID := ValidBookmarkID()
	initialCount := service.GetBookmarkCount()
	
	err := service.Remove(bookmarkID)
	
	assert.NoError(t, err)
	assert.Equal(t, initialCount-1, service.GetBookmarkCount())
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, bookmarkID, service.LastRemoveID)
	
	// Verify bookmark is actually removed
	_, err = service.Get(bookmarkID)
	assert.Error(t, err)
}

// Run_BookmarkRemove_InvalidID executes remove with invalid bookmark ID
func Run_BookmarkRemove_InvalidID(t *testing.T, service *MockBookmarkService) {
	bookmarkID := InvalidBookmarkID()
	initialCount := service.GetBookmarkCount()
	
	err := service.Remove(bookmarkID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bookmark not found")
	assert.Equal(t, initialCount, service.GetBookmarkCount())
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, bookmarkID, service.LastRemoveID)
}

// Run_BookmarkRemove_ForceFailure executes remove with forced failure
func Run_BookmarkRemove_ForceFailure(t *testing.T, service *MockBookmarkService) {
	service.ShouldFailRemove = true
	service.RemoveErrorMessage = "Forced Remove failure"
	bookmarkID := ValidBookmarkID()
	
	err := service.Remove(bookmarkID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Forced Remove failure")
	assert.Equal(t, 1, service.RemoveCallCount)
}

// <tests>
// <evaluators>

// EvaluateBookmarkRemoveSuccess validates successful bookmark removal
func EvaluateBookmarkRemoveSuccess(t *testing.T, service *MockBookmarkService, expectedID uuid.UUID) {
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, expectedID, service.LastRemoveID)
}

// EvaluateBookmarkRemoveFailure validates failed bookmark removal
func EvaluateBookmarkRemoveFailure(t *testing.T, service *MockBookmarkService, expectedError string) {
	assert.Equal(t, 1, service.RemoveCallCount)
}

// <map/>

// BookmarkRemoveTestMap defines all Remove method test cases
var BookmarkRemoveTestMap = map[string]func(*testing.T, *MockBookmarkService){
	"ValidID":      Run_BookmarkRemove_ValidID,
	"InvalidID":    Run_BookmarkRemove_InvalidID,
	"ForceFailure": Run_BookmarkRemove_ForceFailure,
}

// <hook/>

// Test_MockBookmarkService_Remove tests all Remove method scenarios
func Test_MockBookmarkService_Remove(t *testing.T) {
	for name, testFunc := range BookmarkRemoveTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockBookmarkService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>