package services

import (
	"context"
	"testing"

	"github.com/adamkali/mindscape/models/requests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// <method var=services.MockUserService.Create>
// <fixtures/>

// NewUserRequestParams provides base parameters for user creation tests
func NewUserRequestParams() *requests.NewUserRequest {
	return &requests.NewUserRequest{
		Username: "newtestuser",
		Email:    "newtest@example.com",
		Password: "Password123!",
		IsAdmin:  false,
	}
}

// WithBadUsername modifies request with invalid username
func WithBadUsername(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Username = ""
	return req
}

// WithBadEmail modifies request with invalid email
func WithBadEmail(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Email = "invalid-email"
	return req
}

// WithBadPassword modifies request with weak password
func WithBadPassword(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Password = "weak"
	return req
}

// WithDuplicateEmail modifies request with existing email
func WithDuplicateEmail(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Email = "test@example.com" // This email exists in seeded data
	req.Username = "uniqueusername" // Make username unique to test email duplication specifically
	return req
}

// WithDuplicateUsername modifies request with existing username
func WithDuplicateUsername(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Username = "testuser" // This username exists in seeded data
	req.Email = "uniqueemail@example.com" // Make email unique to test username duplication specifically
	return req
}

// WithAdminFlag modifies request to create admin user
func WithAdminFlag(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.IsAdmin = true
	req.Username = "newadminuser"
	req.Email = "newadmin@example.com"
	return req
}

// WithEmptyFields modifies request with all empty fields
func WithEmptyFields(req *requests.NewUserRequest) *requests.NewUserRequest {
	req.Username = ""
	req.Email = ""
	req.Password = ""
	return req
}

// WithForceCreateFailure modifies the service to force create failure
func WithForceCreateFailure(service *MockUserService) *MockUserService {
	service.ShouldFailCreate = true
	service.CreateErrorMessage = "Forced create failure"
	return service
}

// <runners/>

// Run_Create_ValidRequest executes create with valid request
func Run_Create_ValidRequest(t *testing.T, service *MockUserService) {
	req := NewUserRequestParams()
	user, err := service.Create(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.IsAdmin, user.Admin)
	assert.NotEmpty(t, user.BCryptHash)
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_Create_BadUsername executes create with invalid username
func Run_Create_BadUsername(t *testing.T, service *MockUserService) {
	req := WithBadUsername(NewUserRequestParams())
	user, err := service.Create(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_Create_DuplicateEmail executes create with duplicate email
func Run_Create_DuplicateEmail(t *testing.T, service *MockUserService) {
	req := WithDuplicateEmail(NewUserRequestParams())
	user, err := service.Create(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_Create_DuplicateUsername executes create with duplicate username
func Run_Create_DuplicateUsername(t *testing.T, service *MockUserService) {
	req := WithDuplicateUsername(NewUserRequestParams())
	user, err := service.Create(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username already exists")
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_Create_AdminUser executes create for admin user
func Run_Create_AdminUser(t *testing.T, service *MockUserService) {
	req := WithAdminFlag(NewUserRequestParams())
	user, err := service.Create(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.True(t, user.Admin)
	assert.Equal(t, 1, service.CreateCallCount)
}

// Run_Create_ForceFailure executes create with forced failure
func Run_Create_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceCreateFailure(service)
	req := NewUserRequestParams()
	user, err := service.Create(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Forced create failure")
	assert.Equal(t, 1, service.CreateCallCount)
}

// <tests>
// <evaluators>

// EvaluateCreateSuccess validates successful user creation
func EvaluateCreateSuccess(t *testing.T, service *MockUserService, expectedUsername, expectedEmail string, expectedAdmin bool) {
	assert.Equal(t, 1, service.CreateCallCount)
	assert.NotNil(t, service.LastCreateParams)
	assert.Equal(t, expectedUsername, service.LastCreateParams.Username)
	assert.Equal(t, expectedEmail, service.LastCreateParams.Email)
	assert.Equal(t, expectedAdmin, service.LastCreateParams.IsAdmin)
}

// EvaluateCreateFailure validates failed user creation
func EvaluateCreateFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.CreateCallCount)
	user, err := service.Create(service.LastCreateParams)
	assert.Error(t, err)
	assert.Nil(t, user)
	if expectedError != "" {
		assert.Contains(t, err.Error(), expectedError)
	}
}

// <map/>

// CreateTestMap defines all Create method test cases
var CreateTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidRequest":       Run_Create_ValidRequest,
	"BadUsername":        Run_Create_BadUsername,
	"DuplicateEmail":     Run_Create_DuplicateEmail,
	"DuplicateUsername":  Run_Create_DuplicateUsername,
	"AdminUser":          Run_Create_AdminUser,
	"ForceFailure":       Run_Create_ForceFailure,
}

// <hook/>

// Test_MockUserService_Create tests all Create method scenarios
func Test_MockUserService_Create(t *testing.T) {
	for name, testFunc := range CreateTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.Login>
// <fixtures/>

// LoginRequestParams provides base parameters for login tests
func LoginRequestParams() *requests.LoginRequest {
	return &requests.LoginRequest{
		Email:    "test@example.com",
		Password: "password123", // This matches the seeded test data
	}
}

// LoginRequestByUsername provides base parameters for username login
func LoginRequestByUsername() *requests.LoginRequest {
	return &requests.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
}

// WithWrongPassword modifies request with incorrect password
func WithWrongPassword(req *requests.LoginRequest) *requests.LoginRequest {
	req.Password = "wrongpassword"
	return req
}

// WithNonExistentEmail modifies request with non-existent email
func WithNonExistentEmail(req *requests.LoginRequest) *requests.LoginRequest {
	req.Email = "nonexistent@example.com"
	req.Username = ""
	return req
}

// WithNonExistentUsername modifies request with non-existent username
func WithNonExistentUsername(req *requests.LoginRequest) *requests.LoginRequest {
	req.Username = "nonexistentuser"
	req.Email = ""
	return req
}

// WithEmptyCredentials modifies request with empty email and username
func WithEmptyCredentials(req *requests.LoginRequest) *requests.LoginRequest {
	req.Email = ""
	req.Username = ""
	return req
}

// WithForceLoginFailure modifies the service to force login failure
func WithForceLoginFailure(service *MockUserService) *MockUserService {
	service.ShouldFailLogin = true
	service.LoginErrorMessage = "Forced login failure"
	return service
}

// <runners/>

// Run_Login_ValidEmail executes login with valid email
func Run_Login_ValidEmail(t *testing.T, service *MockUserService) {
	req := LoginRequestParams()
	user, err := service.Login(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_ValidUsername executes login with valid username
func Run_Login_ValidUsername(t *testing.T, service *MockUserService) {
	req := LoginRequestByUsername()
	user, err := service.Login(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_WrongPassword executes login with wrong password
func Run_Login_WrongPassword(t *testing.T, service *MockUserService) {
	req := WithWrongPassword(LoginRequestParams())
	user, err := service.Login(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid password")
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_NonExistentEmail executes login with non-existent email
func Run_Login_NonExistentEmail(t *testing.T, service *MockUserService) {
	req := WithNonExistentEmail(LoginRequestParams())
	user, err := service.Login(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_NonExistentUsername executes login with non-existent username
func Run_Login_NonExistentUsername(t *testing.T, service *MockUserService) {
	req := WithNonExistentUsername(LoginRequestParams())
	user, err := service.Login(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_EmptyCredentials executes login with empty credentials
func Run_Login_EmptyCredentials(t *testing.T, service *MockUserService) {
	req := WithEmptyCredentials(LoginRequestParams())
	user, err := service.Login(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email or username required")
	assert.Equal(t, 1, service.LoginCallCount)
}

// Run_Login_ForceFailure executes login with forced failure
func Run_Login_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceLoginFailure(service)
	req := LoginRequestParams()
	user, err := service.Login(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Forced login failure")
	assert.Equal(t, 1, service.LoginCallCount)
}

// <tests>
// <evaluators>

// EvaluateLoginSuccess validates successful login
func EvaluateLoginSuccess(t *testing.T, service *MockUserService, expectedIdentifier string) {
	assert.Equal(t, 1, service.LoginCallCount)
	assert.NotNil(t, service.LastLoginParams)
	if service.LastLoginParams.Email != "" {
		assert.Equal(t, expectedIdentifier, service.LastLoginParams.Email)
	} else {
		assert.Equal(t, expectedIdentifier, service.LastLoginParams.Username)
	}
}

// EvaluateLoginFailure validates failed login
func EvaluateLoginFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.LoginCallCount)
	assert.NotNil(t, service.LastLoginParams)
}

// <map/>

// LoginTestMap defines all Login method test cases
var LoginTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidEmail":         Run_Login_ValidEmail,
	"ValidUsername":      Run_Login_ValidUsername,
	"WrongPassword":      Run_Login_WrongPassword,
	"NonExistentEmail":   Run_Login_NonExistentEmail,
	"NonExistentUsername": Run_Login_NonExistentUsername,
	"EmptyCredentials":   Run_Login_EmptyCredentials,
	"ForceFailure":       Run_Login_ForceFailure,
}

// <hook/>

// Test_MockUserService_Login tests all Login method scenarios
func Test_MockUserService_Login(t *testing.T) {
	for name, testFunc := range LoginTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.Get>
// <fixtures/>

// ValidUserID provides a valid user ID for testing
func ValidUserID() uuid.UUID {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111") // From seeded data
}

// InvalidUserID provides an invalid user ID for testing
func InvalidUserID() uuid.UUID {
	return uuid.MustParse("99999999-9999-9999-9999-999999999999")
}

// WithForceGetFailure modifies the service to force get failure
func WithForceGetFailure(service *MockUserService) *MockUserService {
	service.ShouldFailGet = true
	service.GetErrorMessage = "Forced get failure"
	return service
}

// <runners/>

// Run_Get_ValidID executes get with valid user ID
func Run_Get_ValidID(t *testing.T, service *MockUserService) {
	userID := ValidUserID()
	user, err := service.Get(userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, userID, service.LastGetID)
}

// Run_Get_InvalidID executes get with invalid user ID
func Run_Get_InvalidID(t *testing.T, service *MockUserService) {
	userID := InvalidUserID()
	user, err := service.Get(userID)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, userID, service.LastGetID)
}

// Run_Get_ForceFailure executes get with forced failure
func Run_Get_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceGetFailure(service)
	userID := ValidUserID()
	user, err := service.Get(userID)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Forced get failure")
	assert.Equal(t, 1, service.GetCallCount)
}

// <tests>
// <evaluators>

// EvaluateGetSuccess validates successful user retrieval
func EvaluateGetSuccess(t *testing.T, service *MockUserService, expectedID uuid.UUID) {
	assert.Equal(t, 1, service.GetCallCount)
	assert.Equal(t, expectedID, service.LastGetID)
}

// EvaluateGetFailure validates failed user retrieval
func EvaluateGetFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.GetCallCount)
}

// <map/>

// GetTestMap defines all Get method test cases
var GetTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidID":      Run_Get_ValidID,
	"InvalidID":    Run_Get_InvalidID,
	"ForceFailure": Run_Get_ForceFailure,
}

// <hook/>

// Test_MockUserService_Get tests all Get method scenarios
func Test_MockUserService_Get(t *testing.T) {
	for name, testFunc := range GetTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.Remove>
// <fixtures/>

// WithForceRemoveFailure modifies the service to force remove failure
func WithForceRemoveFailure(service *MockUserService) *MockUserService {
	service.ShouldFailRemove = true
	service.RemoveErrorMessage = "Forced remove failure"
	return service
}

// <runners/>

// Run_Remove_ValidID executes remove with valid user ID
func Run_Remove_ValidID(t *testing.T, service *MockUserService) {
	userID := ValidUserID()
	initialCount := service.GetUserCount()
	
	err := service.Remove(userID)
	
	assert.NoError(t, err)
	assert.Equal(t, initialCount-1, service.GetUserCount())
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, userID, service.LastRemoveID)
	
	// Verify user is actually removed
	_, err = service.Get(userID)
	assert.Error(t, err)
}

// Run_Remove_InvalidID executes remove with invalid user ID
func Run_Remove_InvalidID(t *testing.T, service *MockUserService) {
	userID := InvalidUserID()
	initialCount := service.GetUserCount()
	
	err := service.Remove(userID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, initialCount, service.GetUserCount())
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, userID, service.LastRemoveID)
}

// Run_Remove_ForceFailure executes remove with forced failure
func Run_Remove_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceRemoveFailure(service)
	userID := ValidUserID()
	
	err := service.Remove(userID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Forced remove failure")
	assert.Equal(t, 1, service.RemoveCallCount)
}

// <tests>
// <evaluators>

// EvaluateRemoveSuccess validates successful user removal
func EvaluateRemoveSuccess(t *testing.T, service *MockUserService, expectedID uuid.UUID) {
	assert.Equal(t, 1, service.RemoveCallCount)
	assert.Equal(t, expectedID, service.LastRemoveID)
}

// EvaluateRemoveFailure validates failed user removal
func EvaluateRemoveFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.RemoveCallCount)
}

// <map/>

// RemoveTestMap defines all Remove method test cases
var RemoveTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidID":      Run_Remove_ValidID,
	"InvalidID":    Run_Remove_InvalidID,
	"ForceFailure": Run_Remove_ForceFailure,
}

// <hook/>

// Test_MockUserService_Remove tests all Remove method scenarios
func Test_MockUserService_Remove(t *testing.T) {
	for name, testFunc := range RemoveTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.GetAll>
// <fixtures/>

// WithForceGetAllFailure modifies the service to force getall failure
func WithForceGetAllFailure(service *MockUserService) *MockUserService {
	service.ShouldFailGetAll = true
	service.GetAllErrorMessage = "Forced getall failure"
	return service
}

// <runners/>

// Run_GetAll_Success executes getall successfully
func Run_GetAll_Success(t *testing.T, service *MockUserService) {
	users, err := service.GetAll()
	
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.GreaterOrEqual(t, len(users), 2) // At least 2 seeded users
	assert.Equal(t, 1, service.GetAllCallCount)
}

// Run_GetAll_ForceFailure executes getall with forced failure
func Run_GetAll_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceGetAllFailure(service)
	users, err := service.GetAll()
	
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "Forced getall failure")
	assert.Equal(t, 1, service.GetAllCallCount)
}

// <tests>
// <evaluators>

// EvaluateGetAllSuccess validates successful user retrieval
func EvaluateGetAllSuccess(t *testing.T, service *MockUserService, expectedCount int) {
	assert.Equal(t, 1, service.GetAllCallCount)
	users, err := service.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, expectedCount)
}

// EvaluateGetAllFailure validates failed user retrieval
func EvaluateGetAllFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.GetAllCallCount)
}

// <map/>

// GetAllTestMap defines all GetAll method test cases
var GetAllTestMap = map[string]func(*testing.T, *MockUserService){
	"Success":      Run_GetAll_Success,
	"ForceFailure": Run_GetAll_ForceFailure,
}

// <hook/>

// Test_MockUserService_GetAll tests all GetAll method scenarios
func Test_MockUserService_GetAll(t *testing.T) {
	for name, testFunc := range GetAllTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.Update>
// <fixtures/>

// UpdateProfileParams provides base parameters for update tests
func UpdateProfileParams() (uuid.UUID, string) {
	return ValidUserID(), "new-profile-pic.jpg"
}

// WithForceUpdateFailure modifies the service to force update failure
func WithForceUpdateFailure(service *MockUserService) *MockUserService {
	service.ShouldFailUpdate = true
	service.UpdateErrorMessage = "Forced update failure"
	return service
}

// <runners/>

// Run_Update_ValidParams executes update with valid parameters
func Run_Update_ValidParams(t *testing.T, service *MockUserService) {
	userID, profileName := UpdateProfileParams()
	
	user, err := service.Update(userID, profileName)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, profileName, *user.ProfilePicUrl)
	assert.Equal(t, 1, service.UpdateCallCount)
	assert.Equal(t, userID, service.LastUpdateID)
	assert.Equal(t, profileName, service.LastUpdateProfile)
}

// Run_Update_InvalidID executes update with invalid user ID
func Run_Update_InvalidID(t *testing.T, service *MockUserService) {
	userID := InvalidUserID()
	profileName := "new-profile.jpg"
	
	user, err := service.Update(userID, profileName)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 1, service.UpdateCallCount)
	assert.Equal(t, userID, service.LastUpdateID)
}

// Run_Update_ForceFailure executes update with forced failure
func Run_Update_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceUpdateFailure(service)
	userID, profileName := UpdateProfileParams()
	
	user, err := service.Update(userID, profileName)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Forced update failure")
	assert.Equal(t, 1, service.UpdateCallCount)
}

// <tests>
// <evaluators>

// EvaluateUpdateSuccess validates successful user update
func EvaluateUpdateSuccess(t *testing.T, service *MockUserService, expectedID uuid.UUID, expectedProfile string) {
	assert.Equal(t, 1, service.UpdateCallCount)
	assert.Equal(t, expectedID, service.LastUpdateID)
	assert.Equal(t, expectedProfile, service.LastUpdateProfile)
}

// EvaluateUpdateFailure validates failed user update
func EvaluateUpdateFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.UpdateCallCount)
}

// <map/>

// UpdateTestMap defines all Update method test cases
var UpdateTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidParams":  Run_Update_ValidParams,
	"InvalidID":    Run_Update_InvalidID,
	"ForceFailure": Run_Update_ForceFailure,
}

// <hook/>

// Test_MockUserService_Update tests all Update method scenarios
func Test_MockUserService_Update(t *testing.T) {
	for name, testFunc := range UpdateTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockUserService.UpdateUserCredentials>
// <fixtures/>

// UpdateCredentialsParams provides base parameters for credential update tests
func UpdateCredentialsParams() *requests.UpdateCredentialsRequest {
	return &requests.UpdateCredentialsRequest{
		ID:          ValidUserID(),
		Username:    "newusername",
		Email:       "newemail@example.com",
		OldPassword: "password123", // This matches the seeded data
		Password:    "NewPassword123!",
	}
}

// WithNoPasswordChange modifies request to not change password
func WithNoPasswordChange(req *requests.UpdateCredentialsRequest) *requests.UpdateCredentialsRequest {
	req.OldPassword = ""
	req.Password = ""
	return req
}

// WithWrongOldPassword modifies request with wrong old password
func WithWrongOldPassword(req *requests.UpdateCredentialsRequest) *requests.UpdateCredentialsRequest {
	req.OldPassword = "wrongpassword"
	return req
}

// WithDuplicateEmailCreds modifies request with existing email
func WithDuplicateEmailCreds(req *requests.UpdateCredentialsRequest) *requests.UpdateCredentialsRequest {
	req.Email = "admin@example.com" // This email exists in seeded data
	return req
}

// WithDuplicateUsernameCreds modifies request with existing username
func WithDuplicateUsernameCreds(req *requests.UpdateCredentialsRequest) *requests.UpdateCredentialsRequest {
	req.Username = "adminuser" // This username exists in seeded data
	return req
}

// WithInvalidUserIDCreds modifies request with invalid user ID
func WithInvalidUserIDCreds(req *requests.UpdateCredentialsRequest) *requests.UpdateCredentialsRequest {
	req.ID = InvalidUserID()
	return req
}

// WithForceUpdateCredsFailure modifies the service to force update credentials failure
func WithForceUpdateCredsFailure(service *MockUserService) *MockUserService {
	service.ShouldFailUpdateCreds = true
	service.UpdateCredsErrorMessage = "Forced update credentials failure"
	return service
}

// <runners/>

// Run_UpdateCreds_ValidParams executes update credentials with valid parameters
func Run_UpdateCreds_ValidParams(t *testing.T, service *MockUserService) {
	req := UpdateCredentialsParams()
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, 1, service.UpdateCredsCallCount)
	assert.Equal(t, req, service.LastUpdateCredsParams)
}

// Run_UpdateCreds_NoPasswordChange executes update without password change
func Run_UpdateCreds_NoPasswordChange(t *testing.T, service *MockUserService) {
	req := WithNoPasswordChange(UpdateCredentialsParams())
	originalUser, _ := service.Get(req.ID)
	originalHash := originalUser.BCryptHash
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, originalHash, user.BCryptHash) // Password should remain unchanged
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// Run_UpdateCreds_WrongOldPassword executes update with wrong old password
func Run_UpdateCreds_WrongOldPassword(t *testing.T, service *MockUserService) {
	req := WithWrongOldPassword(UpdateCredentialsParams())
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "old password is incorrect")
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// Run_UpdateCreds_DuplicateEmail executes update with duplicate email
func Run_UpdateCreds_DuplicateEmail(t *testing.T, service *MockUserService) {
	req := WithDuplicateEmailCreds(UpdateCredentialsParams())
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already in use")
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// Run_UpdateCreds_DuplicateUsername executes update with duplicate username
func Run_UpdateCreds_DuplicateUsername(t *testing.T, service *MockUserService) {
	req := WithDuplicateUsernameCreds(UpdateCredentialsParams())
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username already in use")
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// Run_UpdateCreds_InvalidUserID executes update with invalid user ID
func Run_UpdateCreds_InvalidUserID(t *testing.T, service *MockUserService) {
	req := WithInvalidUserIDCreds(UpdateCredentialsParams())
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// Run_UpdateCreds_ForceFailure executes update with forced failure
func Run_UpdateCreds_ForceFailure(t *testing.T, service *MockUserService) {
	WithForceUpdateCredsFailure(service)
	req := UpdateCredentialsParams()
	
	user, err := service.UpdateUserCredentials(req)
	
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Forced update credentials failure")
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// <tests>
// <evaluators>

// EvaluateUpdateCredsSuccess validates successful credential update
func EvaluateUpdateCredsSuccess(t *testing.T, service *MockUserService, expectedParams *requests.UpdateCredentialsRequest) {
	assert.Equal(t, 1, service.UpdateCredsCallCount)
	assert.Equal(t, expectedParams, service.LastUpdateCredsParams)
}

// EvaluateUpdateCredsFailure validates failed credential update
func EvaluateUpdateCredsFailure(t *testing.T, service *MockUserService, expectedError string) {
	assert.Equal(t, 1, service.UpdateCredsCallCount)
}

// <map/>

// UpdateCredsTestMap defines all UpdateUserCredentials method test cases
var UpdateCredsTestMap = map[string]func(*testing.T, *MockUserService){
	"ValidParams":        Run_UpdateCreds_ValidParams,
	"NoPasswordChange":   Run_UpdateCreds_NoPasswordChange,
	"WrongOldPassword":   Run_UpdateCreds_WrongOldPassword,
	"DuplicateEmail":     Run_UpdateCreds_DuplicateEmail,
	"DuplicateUsername":  Run_UpdateCreds_DuplicateUsername,
	"InvalidUserID":      Run_UpdateCreds_InvalidUserID,
	"ForceFailure":       Run_UpdateCreds_ForceFailure,
}

// <hook/>

// Test_MockUserService_UpdateUserCredentials tests all UpdateUserCredentials method scenarios
func Test_MockUserService_UpdateUserCredentials(t *testing.T) {
	for name, testFunc := range UpdateCredsTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockUserService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>