package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/requests"
	"github.com/adamkali/mindscape/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	CreateUser            = "/api/users/signup"
	LoginUser             = "/api/users/login"
	UpdateUserCredentials = "/api/users/creds"
)

// <service var=services.ValidatorService>
// <fixtures/>
func NewEchoContext(r *http.Request) echo.Context {
	return echo.New().NewContext(r, httptest.NewRecorder())
}
func Request(method string, path string, request map[string]any) *http.Request {
	requestJson, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(string(requestJson)),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// <method var=services.ValidatorService.ValidateNewUserRequest>
// <fixtures>
func WithIsAdmin(request map[string]any) map[string]any {
	request["is_admin"] = true
	return request
}

func WithBadUsername(request map[string]any) map[string]any {
	request["username"] = "!!username!!"
	return request
}

func WithBadUsername_SQLInjection(request map[string]any) map[string]any {
	request["username"] = "username'; DROP TABLE users;"
	return request
}

func WithBadEmail(request map[string]any) map[string]any {
	request["email"] = "chatarewecooked"
	return request
}

func WithBadPassword_NotSevenOrMore(request map[string]any) map[string]any {
	request["password"] = "Aa1!"
	return request
}

func WithBadPassword_NoNumber(request map[string]any) map[string]any {
	request["password"] = "passwordabc!"
	return request
}

func WithBadPassword_NoSpecialCharacter(request map[string]any) map[string]any {
	request["password"] = "passwordABC123"
	return request
}

func WithBadPassword_NoUpper(request map[string]any) map[string]any {
	request["password"] = "passwordabc123!"
	return request
}

func WithBadPassword_NoLower(request map[string]any) map[string]any {
	request["password"] = "PASSWORDABC123!"
	return request
}

func WithBadPassword_NoSpecial(request map[string]any) map[string]any {
	request["password"] = "passwordABC123"
	return request
}

func WithShortMixedPassword(request map[string]any) map[string]any {
	request["password"] = "Cheese360!"
	return request
}

func UserRequestJson() map[string]any {
	return map[string]any{
		"username": "adamkali",
		"email":    "iHr4l@example.com",
		"password": "passwordABC123!",
		"is_admin": false,
	}
}

// </fixtures>
// <tests>
// <runners>
func Run_ValidateNewUserRequest(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		UserRequestJson(),
	)))
}

func Run_ValidateNewUserRequest_WithAdmin(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithIsAdmin(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadUsername(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadUsername(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadUsername_SQLInjection(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadUsername_SQLInjection(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithIsAdmin(WithBadUsername_SQLInjection(UserRequestJson())))))
}

func Run_ValidateNewUserRequest_BadEmail(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadEmail(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadEmail_WithAdmin(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithIsAdmin(WithBadEmail(UserRequestJson())))))
}

func Run_ValidateNewUserRequest_BadPassword_NotSevenOrMore(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NotSevenOrMore(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadPassword_NoUpper(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NoUpper(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadPassword_NoSpecialCharacter(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NoSpecialCharacter(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadPassword_NoNumber(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NoNumber(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadPassword_NoLower(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NoLower(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_BadPassword_NoSpecial(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithBadPassword_NoSpecial(UserRequestJson()))))
}

func Run_ValidateNewUserRequest_ShortMixedPassword(v services.ValidatorService) (*requests.NewUserRequest, error) {
	return v.ValidateNewUserRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithShortMixedPassword(UserRequestJson()))))
}

// </runners>
// <evaluators>
func ValidateNewUserRequest_Default(t *testing.T) {
	r, e := Run_ValidateNewUserRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
	assert.Equal(t, r.IsAdmin, false)
}

func ValidateNewUserRequest_WithAdmin(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_WithAdmin(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
	assert.Equal(t, r.IsAdmin, false)
}

func ValidateNewUserRequest_BadUsername(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Username, "!!username!!")
}

func ValidateNewUserRequest_BadUsername_SQLInjection(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername_SQLInjection(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Username, "username'; DROP TABLE users;")
}

func ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Username, "username'; DROP TABLE users;")
	assert.Equal(t, r.IsAdmin, false)
}

func ValidateNewUserRequest_BadEmail(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadEmail(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed (chatarewecooked) is not a valid email address",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadEmail_WithAdmin(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadEmail_WithAdmin(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed (chatarewecooked) is not a valid email address",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NotSevenOrMore(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NotSevenOrMore(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (true), Upper (true), Special (true)",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoUpper(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoUpper(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (true), Number (true), Upper (false), Special (true)",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoSpecialCharacter(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoSpecialCharacter(services.ValidatorService{})
	assert.Error(t, e)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoNumber(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoNumber(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (true), Number (false), Upper (false), Special (true)",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoLower(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoLower(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "PASSWORDABC123!")
}

func ValidateNewUserRequest_ShortMixedPassword(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_ShortMixedPassword(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "Cheese360!")
}

// </evaluators>
// <map/>
var Map_ValidateNewUserRequest = map[string]func(t *testing.T){
	"Default":                            ValidateNewUserRequest_Default,
	"WithAdmin":                          ValidateNewUserRequest_WithAdmin,
	"BadUsername":                        ValidateNewUserRequest_BadUsername,
	"BadUsername_SQLInjection":           ValidateNewUserRequest_BadUsername_SQLInjection,
	"BadUsername_SQLInjection_WithAdmin": ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin,
	"BadEmail":                           ValidateNewUserRequest_BadEmail,
	"BadEmail_WithAdmin":                 ValidateNewUserRequest_BadEmail_WithAdmin,
	"BadPassword_NotSevenOrMore":         ValidateNewUserRequest_BadPassword_NotSevenOrMore,
	"BadPassword_NoUpper":                ValidateNewUserRequest_BadPassword_NoUpper,
	"BadPassword_NoSpecialCharacter":     ValidateNewUserRequest_BadPassword_NoSpecialCharacter,
	"BadPassword_NoNumber":               ValidateNewUserRequest_BadPassword_NoNumber,
	"BadPassword_NoLower":                ValidateNewUserRequest_BadPassword_NoLower,
	"ShortMixedPassword":                 ValidateNewUserRequest_ShortMixedPassword,
}

// </tests>
// <hook/>
func Test_ValidateNewUserRequest(t *testing.T) {
	fmt.Println("Test_ValidateNewUserRequest")
	for k, v := range Map_ValidateNewUserRequest {
		t.Run(k, v)
	}
}

// </method>
// <method var=services.ValidatorService.ValidateLoginRequest>
// <fixtures>
// <fixture.default/>
func LoginRequestJson() map[string]any {
	return map[string]any{
		"username": "adamkali",
		"email":    "iHr4l@example.com",
		"password": "passwordABC123!",
	}
}

func WithUsername(request map[string]any) map[string]any {
	request["email"] = ""
	return request
}

func WithEmail(request map[string]any) map[string]any {
	request["username"] = ""
	return request
}

func WithOnlyUsername(request map[string]any) map[string]any {
	request["email"] = ""
	request["password"] = ""
	return request
}

func WithOnlyEmail(request map[string]any) map[string]any {
	request["username"] = ""
	request["password"] = ""
	return request
}

func WithOnlyPassword(request map[string]any) map[string]any {
	request["username"] = ""
	request["email"] = ""
	return request
}

// </fixtures>
// <runners>
func Run_ValidateLoginRequest(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		LoginRequestJson(),
	)))
}

func Run_ValidateLoginRequest_WithUsername(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithUsername(LoginRequestJson()),
	)))
}

func Run_ValidateLoginRequest_WithEmail(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithEmail(LoginRequestJson()),
	)))
}

func Run_ValidateLoginRequest_WithOnlyUsername(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithOnlyUsername(LoginRequestJson()),
	)))
}

func Run_ValidateLoginRequest_WithOnlyEmail(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithOnlyEmail(LoginRequestJson()),
	)))
}

func Run_ValidateLoginRequest_WithOnlyPassword(s services.ValidatorService) (*requests.LoginRequest, error) {
	return s.ValidateLoginRequest(NewEchoContext(Request(
		http.MethodPost,
		CreateUser,
		WithOnlyPassword(LoginRequestJson()),
	)))
}

// </runners>
// <tests>
// <evaluators>
func ValidateLoginRequest_Default(t *testing.T) {
	r, e := Run_ValidateLoginRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginRequest_WithUsername(t *testing.T) {
	r, e := Run_ValidateLoginRequest_WithUsername(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginRequest_WithEmail(t *testing.T) {
	r, e := Run_ValidateLoginRequest_WithEmail(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginRequest_WithOnlyUsername(t *testing.T) {
	r, e := Run_ValidateLoginRequest_WithOnlyUsername(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"You must send a password",
	)
	assert.Nil(t, r)
}

func ValidateLoginRequest_WithOnlyEmail(t *testing.T) {
	r, e := Run_ValidateLoginRequest_WithOnlyEmail(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"You must send a password",
	)
	assert.Nil(t, r)
}

func ValidateLoginRequest_WithOnlyPassword(t *testing.T) {
	r, e := Run_ValidateLoginRequest_WithOnlyPassword(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Email and Username cannot be null",
	)
	assert.Nil(t, r)
}

// </evaluators>
// <map/>
var Map_ValidateLoginRequest = map[string]func(*testing.T){
	"Default":          ValidateLoginRequest_Default,
	"WithUsername":     ValidateLoginRequest_WithUsername,
	"WithEmail":        ValidateLoginRequest_WithEmail,
	"WithOnlyUsername": ValidateLoginRequest_WithOnlyUsername,
	"WithOnlyEmail":    ValidateLoginRequest_WithOnlyEmail,
	"WithOnlyPassword": ValidateLoginRequest_WithOnlyPassword,
}

// </tests>
// <hook/>
func Test_ValidateLoginRequest(t *testing.T) {
	fmt.Println("Test_ValidateLoginRequest")
	for k, v := range Map_ValidateLoginRequest {
		t.Run(k, v)
	}
}

// <method var=services.ValidatorService.ValidateUpdateUserCredentialRequest>
// <fixtures>
// <fixture.default/>
func UpdateUserCredentialRequestJson() map[string]any {
	return map[string]any{
		"id":           "e38e78a4-2ca3-4c59-a3ea-a2019866e593",
		"username":     "adamkali",
		"email":        "iHr4l@example.com",
		"password":     "WowPoggers123!",
		"old_password": "passwordABC123!",
	}
}

func WithNilID(request map[string]any) map[string]any {
	request["id"] = uuid.Nil
	return request
}

func WithBadID(request map[string]any) map[string]any {
	request["id"] = "e38e78a4-2ca3-4c59-a3ea"
	return request
}

func WithBadEmail_SQLInjection(request map[string]any) map[string]any {
	request["email"] = "email'; DROP TABLE users;"
	return request
}

func WithPassword_NoOldPassword(request map[string]any) map[string]any {
	request["old_password"] = ""
	return request
}

// </fixtures>
// <runners>
func Run_ValidateUpdateUserCredentialRequest(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		UpdateUserCredentialRequestJson(),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithNilID(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithNilID(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadID(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadID(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadEmail_SQLInjection(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadUsername(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadUsername(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadUsername_SQLInjection(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoOldPassword(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithPassword_NoOldPassword(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadPassword_NotSevenOrMore(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadPassword_NoLower(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadPassword_NoUpper(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadPassword_NoNumber(UpdateUserCredentialRequestJson()),
	)))
}
func Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial(s services.ValidatorService) (*requests.UpdateCredentialsRequest, error) {
	return s.ValidateUpdateUserCredentialRequest(NewEchoContext(Request(
		http.MethodPost,
		UpdateUserCredentials,
		WithBadPassword_NoSpecial(UpdateUserCredentialRequestJson()),
	)))
}

// </runners>
// <evaluators>
func ValidateUpdateUserCredentialRequest_Default(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "WowPoggers123!")
}
func ValidateUpdateUserCredentialRequest_WithNilID(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithNilID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"ID cannot be null",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadID(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length")
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "missing '@'")
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadUsername(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadUsername(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Username, "!!username!!")
}
func ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Username, "username'; DROP TABLE users;")
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "Aa1!")
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "PASSWORDABC123!")
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "passwordabc123!")
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "passwordabc!")
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Password, "passwordABC123")
}

// </evaluators>
// <map/>
var Map_ValidateUpdateUserCredentialRequest = map[string]func(*testing.T){
	"Default":                        ValidateUpdateUserCredentialRequest_Default,
	"WithNilID":                      ValidateUpdateUserCredentialRequest_WithNilID,
	"WithBadID":                      ValidateUpdateUserCredentialRequest_WithBadID,
	"WithBadEmail_SQLInjection":      ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection,
	"WithBadUsername":                ValidateUpdateUserCredentialRequest_WithBadUsername,
	"WithBadUsername_SQLInjection":   ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection,
	"WithBadPassword_NotSevenOrMore": ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore,
	"WithBadPassword_NoLower":        ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower,
	"WithBadPassword_NoUpper":        ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper,
	"WithBadPassword_NoNumber":       ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber,
	"WithBadPassword_NoSpecial":      ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial,
}

// <hook/>
func Test_ValidateUpdateUserCredentialRequest(t *testing.T) {
	fmt.Println("Test_ValidateUpdateUserCredentialRequest")
	for k, v := range Map_ValidateUpdateUserCredentialRequest {
		t.Run(k, v)
	}
}

// </method>
// <method var=services.ValidatorService.ValidateCreateFolderRequest>
// <fixtures>
// <fixture.default/>
func CreateFolderRequestJson() map[string]any {
	return map[string]any{
		"user_id":     "e38e78a4-2ca3-4c59-a3ea-a2019866e593",
		"parent_id":   "34789a5b-daa1-4f03-b912-e4e4e79dae53",
		"name":        "Test Folder",
		"description": "This is a test folder",
	}
}
func WithBadParentID(request map[string]any) map[string]any {
	request["parent_id"] = "e38e78a4-2ca3-4c59-a3ea"
	return request
}

func WithEmptyName(request map[string]any) map[string]any {
	request["name"] = ""
	return request
}

func WithNilName(request map[string]any) map[string]any {
	delete(request, "name")
	return request
}

func WithSpecialCharacterName_Slash(request map[string]any) map[string]any {
	request["name"] = "Test/Folder"
	return request
}

func WithSpecialCharacterName_Backslash(request map[string]any) map[string]any {
	request["name"] = "Test\\Folder"
	return request
}

func WithSpecialCharacterName_Colon(request map[string]any) map[string]any {
	request["name"] = "Test:Folder"
	return request
}

func WithSpecialCharacterName_Asterisk(request map[string]any) map[string]any {
	request["name"] = "Test*Folder"
	return request
}

func WithSpecialCharacterName_Question(request map[string]any) map[string]any {
	request["name"] = "Test?Folder"
	return request
}

func WithSpecialCharacterName_Quote(request map[string]any) map[string]any {
	request["name"] = "Test\"Folder"
	return request
}

func WithSpecialCharacterName_LessThan(request map[string]any) map[string]any {
	request["name"] = "Test<Folder"
	return request
}

func WithSpecialCharacterName_GreaterThan(request map[string]any) map[string]any {
	request["name"] = "Test>Folder"
	return request
}

func WithSpecialCharacterName_Pipe(request map[string]any) map[string]any {
	request["name"] = "Test|Folder"
	return request
}

func WithValidSpecialName(request map[string]any) map[string]any {
	request["name"] = "Test Folder - Project (2024)"
	return request
}

func WithEmojiName(request map[string]any) map[string]any {
	request["name"] = "📂 Project Folder 🚀"
	return request
}

// </fixtures>
// <runners>
func Run_ValidateCreateFolderRequest(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		CreateFolderRequestJson(),
	)))
}

func Run_ValidateCreateFolderRequest_WithEmptyName(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithEmptyName(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithNilName(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithNilName(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Slash(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Slash(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Backslash(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Backslash(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Colon(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Colon(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Asterisk(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Asterisk(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Question(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Question(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Quote(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Quote(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_LessThan(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_LessThan(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_GreaterThan(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_GreaterThan(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Pipe(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithSpecialCharacterName_Pipe(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithValidSpecialName(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithValidSpecialName(CreateFolderRequestJson()),
	)))
}

func Run_ValidateCreateFolderRequest_WithEmojiName(s services.ValidatorService) (*repository.CreateFolderParams, error) {
	return s.ValidateCreateFolderRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/folders/create",
		WithEmojiName(CreateFolderRequestJson()),
	)))
}

// </runners>
// <tests>
// <evaluators>
func ValidateCreateFolderRequest_Default(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "Test Folder")
	assert.Equal(t, *r.Description, "This is a test folder")
}

func ValidateCreateFolderRequest_WithEmptyName(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithEmptyName(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot be empty",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithNilName(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithNilName(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot be empty",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Slash(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Slash(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Backslash(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Backslash(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Colon(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Colon(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Asterisk(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Asterisk(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Question(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Question(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Quote(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Quote(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_LessThan(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_LessThan(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_GreaterThan(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_GreaterThan(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithSpecialCharacterName_Pipe(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithSpecialCharacterName_Pipe(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder name cannot contain any special characters",
	)
	assert.Nil(t, r)
}

func ValidateCreateFolderRequest_WithValidSpecialName(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithValidSpecialName(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "Test Folder - Project (2024)")
	assert.Equal(t, *r.Description, "This is a test folder")
}

func ValidateCreateFolderRequest_WithEmojiName(t *testing.T) {
	r, e := Run_ValidateCreateFolderRequest_WithEmojiName(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "📂 Project Folder 🚀")
	assert.Equal(t, *r.Description, "This is a test folder")
}

// </evaluators>
// <map/>
var Map_ValidateCreateFolderRequest = map[string]func(*testing.T){
	"Default":                           ValidateCreateFolderRequest_Default,
	"WithEmptyName":                     ValidateCreateFolderRequest_WithEmptyName,
	"WithNilName":                       ValidateCreateFolderRequest_WithNilName,
	"WithSpecialCharacterName_Slash":    ValidateCreateFolderRequest_WithSpecialCharacterName_Slash,
	"WithSpecialCharacterName_Backslash": ValidateCreateFolderRequest_WithSpecialCharacterName_Backslash,
	"WithSpecialCharacterName_Colon":    ValidateCreateFolderRequest_WithSpecialCharacterName_Colon,
	"WithSpecialCharacterName_Asterisk": ValidateCreateFolderRequest_WithSpecialCharacterName_Asterisk,
	"WithSpecialCharacterName_Question": ValidateCreateFolderRequest_WithSpecialCharacterName_Question,
	"WithSpecialCharacterName_Quote":    ValidateCreateFolderRequest_WithSpecialCharacterName_Quote,
	"WithSpecialCharacterName_LessThan": ValidateCreateFolderRequest_WithSpecialCharacterName_LessThan,
	"WithSpecialCharacterName_GreaterThan": ValidateCreateFolderRequest_WithSpecialCharacterName_GreaterThan,
	"WithSpecialCharacterName_Pipe":     ValidateCreateFolderRequest_WithSpecialCharacterName_Pipe,
	"WithValidSpecialName":              ValidateCreateFolderRequest_WithValidSpecialName,
	"WithEmojiName":                     ValidateCreateFolderRequest_WithEmojiName,
}

// <hook/>
func Test_ValidateCreateFolderRequest(t *testing.T) {
	fmt.Println("Test_ValidateCreateFolderRequest")
	for k, v := range Map_ValidateCreateFolderRequest {
		t.Run(k, v)
	}
}

// </method>
// <method var=services.ValidatorService.ValidateLoginFormRequest>
// <fixtures>
// <fixture.default/>
func LoginFormRequestJson() map[string]any {
	return map[string]any{
		"username": "adamkali",
		"email":    "iHr4l@example.com",
		"password": "passwordABC123!",
	}
}

func WithFormUsername(request map[string]any) map[string]any {
	request["email"] = ""
	return request
}

func WithFormEmail(request map[string]any) map[string]any {
	request["username"] = ""
	return request
}

func WithFormOnlyUsername(request map[string]any) map[string]any {
	request["email"] = ""
	request["password"] = ""
	return request
}

func WithFormOnlyEmail(request map[string]any) map[string]any {
	request["username"] = ""
	request["password"] = ""
	return request
}

func WithFormOnlyPassword(request map[string]any) map[string]any {
	request["username"] = ""
	request["email"] = ""
	return request
}

func WithFormNoPassword(request map[string]any) map[string]any {
	request["password"] = ""
	return request
}

// </fixtures>
// <runners>
func Run_ValidateLoginFormRequest(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("username", "adamkali")
	req.PostForm.Set("email", "iHr4l@example.com")
	req.PostForm.Set("password", "passwordABC123!")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormUsername(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("username", "adamkali")
	req.PostForm.Set("password", "passwordABC123!")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormEmail(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("email", "iHr4l@example.com")
	req.PostForm.Set("password", "passwordABC123!")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormOnlyUsername(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("username", "adamkali")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormOnlyEmail(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("email", "iHr4l@example.com")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormOnlyPassword(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("password", "passwordABC123!")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

func Run_ValidateLoginFormRequest_WithFormNoPassword(s services.ValidatorService) (*requests.LoginRequest, error) {
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = make(map[string][]string)
	req.PostForm.Set("username", "adamkali")
	req.PostForm.Set("email", "iHr4l@example.com")
	return s.ValidateLoginFormRequest(echo.New().NewContext(req, httptest.NewRecorder()))
}

// </runners>
// <tests>
// <evaluators>
func ValidateLoginFormRequest_Default(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginFormRequest_WithFormUsername(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormUsername(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "adamkali")
	assert.Equal(t, r.Email, "")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginFormRequest_WithFormEmail(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormEmail(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Username, "")
	assert.Equal(t, r.Email, "iHr4l@example.com")
	assert.Equal(t, r.Password, "passwordABC123!")
}

func ValidateLoginFormRequest_WithFormOnlyUsername(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormOnlyUsername(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"You must send a password",
	)
	assert.Nil(t, r)
}

func ValidateLoginFormRequest_WithFormOnlyEmail(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormOnlyEmail(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"You must send a password",
	)
	assert.Nil(t, r)
}

func ValidateLoginFormRequest_WithFormOnlyPassword(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormOnlyPassword(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Email and Username cannot be null",
	)
	assert.Nil(t, r)
}

func ValidateLoginFormRequest_WithFormNoPassword(t *testing.T) {
	r, e := Run_ValidateLoginFormRequest_WithFormNoPassword(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"You must send a password",
	)
	assert.Nil(t, r)
}

// </evaluators>
// <map/>
var Map_ValidateLoginFormRequest = map[string]func(*testing.T){
	"Default":           ValidateLoginFormRequest_Default,
	"WithFormUsername":  ValidateLoginFormRequest_WithFormUsername,
	"WithFormEmail":     ValidateLoginFormRequest_WithFormEmail,
	"WithFormOnlyUsername": ValidateLoginFormRequest_WithFormOnlyUsername,
	"WithFormOnlyEmail": ValidateLoginFormRequest_WithFormOnlyEmail,
	"WithFormOnlyPassword": ValidateLoginFormRequest_WithFormOnlyPassword,
	"WithFormNoPassword": ValidateLoginFormRequest_WithFormNoPassword,
}

// <hook/>
func Test_ValidateLoginFormRequest(t *testing.T) {
	fmt.Println("Test_ValidateLoginFormRequest")
	for k, v := range Map_ValidateLoginFormRequest {
		t.Run(k, v)
	}
}

// </method>
// <method var=services.ValidatorService.CreateBookmarkRequest>
// <fixtures>
// <fixture.default/>
func CreateBookmarkRequestJson() map[string]any {
	return map[string]any{
		"user_id":   "e38e78a4-2ca3-4c59-a3ea-a2019866e593",
		"folder_id": "34789a5b-daa1-4f03-b912-e4e4e79dae53",
		"name":      "My Bookmark",
		"link":      "https://example.com",
	}
}

func WithEmptyBookmarkName(request map[string]any) map[string]any {
	request["name"] = ""
	return request
}

func WithNilBookmarkName(request map[string]any) map[string]any {
	delete(request, "name")
	return request
}

func WithEmptyLink(request map[string]any) map[string]any {
	request["link"] = ""
	return request
}

func WithNilLink(request map[string]any) map[string]any {
	delete(request, "link")
	return request
}

func WithInvalidLink(request map[string]any) map[string]any {
	request["link"] = "not-a-valid-url"
	return request
}

func WithEmojiInBookmarkName(request map[string]any) map[string]any {
	request["name"] = "🔖 My Favorite Site 🚀"
	return request
}

func WithEmojiInLink(request map[string]any) map[string]any {
	request["link"] = "https://example.com/🚀"
	return request
}

func WithQueryParamsLink(request map[string]any) map[string]any {
	request["link"] = "https://us-east-1.console.aws.amazon.com/s3/object/bucket?region=us-east-1&prefix=CA/court/"
	return request
}

func WithValidHTTPLink(request map[string]any) map[string]any {
	request["link"] = "http://example.com"
	return request
}

func WithValidHTTPSLink(request map[string]any) map[string]any {
	request["link"] = "https://example.com/page"
	return request
}

func WithJavascriptLink(request map[string]any) map[string]any {
	request["link"] = "javascript:alert('xss')"
	return request
}

func WithDataLink(request map[string]any) map[string]any {
	request["link"] = "data:text/html,<script>alert('xss')</script>"
	return request
}

func WithFTPLink(request map[string]any) map[string]any {
	request["link"] = "ftp://example.com/file.txt"
	return request
}

func WithNilUserID(request map[string]any) map[string]any {
	delete(request, "user_id")
	return request
}

func WithInvalidUserID(request map[string]any) map[string]any {
	request["user_id"] = "invalid-uuid"
	return request
}

func WithNilFolderID(request map[string]any) map[string]any {
	delete(request, "folder_id")
	return request
}

func WithInvalidFolderID(request map[string]any) map[string]any {
	request["folder_id"] = "invalid-uuid"
	return request
}

// </fixtures>
// <runners>
func Run_CreateBookmarkRequest(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		CreateBookmarkRequestJson(),
	)))
}

func Run_CreateBookmarkRequest_WithEmptyBookmarkName(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithEmptyBookmarkName(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithNilBookmarkName(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithNilBookmarkName(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithEmptyLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithEmptyLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithNilLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithNilLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithInvalidLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithInvalidLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithEmojiInBookmarkName(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithEmojiInBookmarkName(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithEmojiInLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithEmojiInLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithQueryParamsLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithQueryParamsLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithValidHTTPLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithValidHTTPLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithValidHTTPSLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithValidHTTPSLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithJavascriptLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithJavascriptLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithDataLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithDataLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithFTPLink(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithFTPLink(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithNilUserID(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithNilUserID(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithInvalidUserID(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithInvalidUserID(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithNilFolderID(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithNilFolderID(CreateBookmarkRequestJson()),
	)))
}

func Run_CreateBookmarkRequest_WithInvalidFolderID(s services.ValidatorService) (*repository.CreateBookmarkParams, error) {
	return s.CreateBookmarkRequest(NewEchoContext(Request(
		http.MethodPost,
		"/api/bookmarks/create",
		WithInvalidFolderID(CreateBookmarkRequestJson()),
	)))
}

// </runners>
// <tests>
// <evaluators>
func CreateBookmarkRequest_Default(t *testing.T) {
	r, e := Run_CreateBookmarkRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "My Bookmark")
	assert.Equal(t, r.Link, "https://example.com")
}

func CreateBookmarkRequest_WithEmptyBookmarkName(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithEmptyBookmarkName(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Bookmark name cannot be empty",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithNilBookmarkName(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithNilBookmarkName(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Bookmark name cannot be empty",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithEmptyLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithEmptyLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Bookmark link cannot be empty",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithNilLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithNilLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Bookmark link cannot be empty",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithInvalidLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithInvalidLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Only HTTP and HTTPS URLs are allowed",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithEmojiInBookmarkName(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithEmojiInBookmarkName(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "🔖 My Favorite Site 🚀")
	assert.Equal(t, r.Link, "https://example.com")
}

func CreateBookmarkRequest_WithEmojiInLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithEmojiInLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Link cannot contain emojis",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithQueryParamsLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithQueryParamsLink(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Name, "My Bookmark")
	assert.Equal(t, r.Link, "https://us-east-1.console.aws.amazon.com/s3/object/bucket?region=us-east-1&prefix=CA/court/")
}

func CreateBookmarkRequest_WithValidHTTPLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithValidHTTPLink(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.Name, "My Bookmark")
	assert.Equal(t, r.Link, "http://example.com")
}

func CreateBookmarkRequest_WithValidHTTPSLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithValidHTTPSLink(services.ValidatorService{})
	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, r.Name, "My Bookmark")
	assert.Equal(t, r.Link, "https://example.com/page")
}

func CreateBookmarkRequest_WithJavascriptLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithJavascriptLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Only HTTP and HTTPS URLs are allowed",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithDataLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithDataLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Only HTTP and HTTPS URLs are allowed",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithFTPLink(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithFTPLink(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Only HTTP and HTTPS URLs are allowed",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithNilUserID(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithNilUserID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"User ID cannot be null",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithInvalidUserID(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithInvalidUserID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length: 12")
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithNilFolderID(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithNilFolderID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Folder ID cannot be null",
	)
	assert.Nil(t, r)
}

func CreateBookmarkRequest_WithInvalidFolderID(t *testing.T) {
	r, e := Run_CreateBookmarkRequest_WithInvalidFolderID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length: 12")
	assert.Nil(t, r)
}

// </evaluators>
// <map/>
var Map_CreateBookmarkRequest = map[string]func(*testing.T){
	"Default":                    CreateBookmarkRequest_Default,
	"WithEmptyBookmarkName":      CreateBookmarkRequest_WithEmptyBookmarkName,
	"WithNilBookmarkName":        CreateBookmarkRequest_WithNilBookmarkName,
	"WithEmptyLink":              CreateBookmarkRequest_WithEmptyLink,
	"WithNilLink":                CreateBookmarkRequest_WithNilLink,
	"WithInvalidLink":            CreateBookmarkRequest_WithInvalidLink,
	"WithEmojiInBookmarkName":    CreateBookmarkRequest_WithEmojiInBookmarkName,
	"WithEmojiInLink":            CreateBookmarkRequest_WithEmojiInLink,
	"WithQueryParamsLink":        CreateBookmarkRequest_WithQueryParamsLink,
	"WithValidHTTPLink":          CreateBookmarkRequest_WithValidHTTPLink,
	"WithValidHTTPSLink":         CreateBookmarkRequest_WithValidHTTPSLink,
	"WithJavascriptLink":         CreateBookmarkRequest_WithJavascriptLink,
	"WithDataLink":               CreateBookmarkRequest_WithDataLink,
	"WithFTPLink":                CreateBookmarkRequest_WithFTPLink,
	"WithNilUserID":              CreateBookmarkRequest_WithNilUserID,
	"WithInvalidUserID":          CreateBookmarkRequest_WithInvalidUserID,
	"WithNilFolderID":            CreateBookmarkRequest_WithNilFolderID,
	"WithInvalidFolderID":        CreateBookmarkRequest_WithInvalidFolderID,
}

// <hook/>
func Test_CreateBookmarkRequest(t *testing.T) {
	fmt.Println("Test_CreateBookmarkRequest")
	for k, v := range Map_CreateBookmarkRequest {
		t.Run(k, v)
	}
}

// </method>

// <method var=services.ValidatorService.ValidateMoveBookmarkRequest>
// <fixtures>
// <fixture.default/>
func MoveBookmarkRequestJson() map[string]any {
	return map[string]any{
		"userId":      "e38e78a4-2ca3-4c59-a3ea-a2019866e593",
		"bookmarkId": "f47f89c5-ebb2-4f04-c913-f5f5f80eaf64",
		"newParentId": "34789a5b-daa1-4f03-b912-e4e4e79dae53",
	}
}

func WithNilMoveUserID(request map[string]any) map[string]any {
	delete(request, "userId")
	return request
}

func WithNilMoveBookmarkID(request map[string]any) map[string]any {
	delete(request, "bookmarkId")
	return request
}

func WithNilMoveParentID(request map[string]any) map[string]any {
	delete(request, "newParentId")
	return request
}

func WithInvalidMoveUserID(request map[string]any) map[string]any {
	request["userId"] = "invalid-uuid"
	return request
}

func WithInvalidMoveBookmarkID(request map[string]any) map[string]any {
	request["bookmarkId"] = "invalid-uuid"
	return request
}

func WithInvalidMoveParentID(request map[string]any) map[string]any {
	request["newParentId"] = "invalid-uuid"
	return request
}

// </fixtures>
// <runners>
func Run_ValidateMoveBookmarkRequest(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		MoveBookmarkRequestJson(),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithNilMoveUserID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithNilMoveUserID(MoveBookmarkRequestJson()),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithNilMoveBookmarkID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithNilMoveBookmarkID(MoveBookmarkRequestJson()),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithNilMoveParentID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithNilMoveParentID(MoveBookmarkRequestJson()),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithInvalidMoveUserID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithInvalidMoveUserID(MoveBookmarkRequestJson()),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithInvalidMoveBookmarkID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithInvalidMoveBookmarkID(MoveBookmarkRequestJson()),
	)))
}

func Run_ValidateMoveBookmarkRequest_WithInvalidMoveParentID(s services.ValidatorService) (*requests.MoveBookmarkRequest, error) {
	return s.ValidateMoveBookmarkRequest(NewEchoContext(Request(
		http.MethodPatch,
		"/api/bookmarks/move",
		WithInvalidMoveParentID(MoveBookmarkRequestJson()),
	)))
}

// </runners>
// <tests>
// <evaluators>
func ValidateMoveBookmarkRequest_Default(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest(services.ValidatorService{})
	assert.NoError(t, e)
	assert.Equal(t, r.UserID.String(), "e38e78a4-2ca3-4c59-a3ea-a2019866e593")
	assert.Equal(t, r.BookmarkID.String(), "f47f89c5-ebb2-4f04-c913-f5f5f80eaf64")
	assert.Equal(t, r.NewParentID.String(), "34789a5b-daa1-4f03-b912-e4e4e79dae53")
}

func ValidateMoveBookmarkRequest_WithNilMoveUserID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithNilMoveUserID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"User ID cannot be null",
	)
	assert.Nil(t, r)
}

func ValidateMoveBookmarkRequest_WithNilMoveBookmarkID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithNilMoveBookmarkID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Bookmark ID cannot be null",
	)
	assert.Nil(t, r)
}

func ValidateMoveBookmarkRequest_WithNilMoveParentID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithNilMoveParentID(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"New parent folder ID cannot be null - bookmarks must belong to a folder",
	)
	assert.Nil(t, r)
}

func ValidateMoveBookmarkRequest_WithInvalidMoveUserID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithInvalidMoveUserID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length: 12")
	assert.Nil(t, r)
}

func ValidateMoveBookmarkRequest_WithInvalidMoveBookmarkID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithInvalidMoveBookmarkID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length: 12")
	assert.Nil(t, r)
}

func ValidateMoveBookmarkRequest_WithInvalidMoveParentID(t *testing.T) {
	r, e := Run_ValidateMoveBookmarkRequest_WithInvalidMoveParentID(services.ValidatorService{})
	assert.Error(t, e)
	assert.Contains(t, e.Error(), "invalid UUID length: 12")
	assert.Nil(t, r)
}

// </evaluators>
// <map/>
var Map_ValidateMoveBookmarkRequest = map[string]func(*testing.T){
	"Default":                      ValidateMoveBookmarkRequest_Default,
	"WithNilMoveUserID":           ValidateMoveBookmarkRequest_WithNilMoveUserID,
	"WithNilMoveBookmarkID":       ValidateMoveBookmarkRequest_WithNilMoveBookmarkID,
	"WithNilMoveParentID":         ValidateMoveBookmarkRequest_WithNilMoveParentID,
	"WithInvalidMoveUserID":       ValidateMoveBookmarkRequest_WithInvalidMoveUserID,
	"WithInvalidMoveBookmarkID":   ValidateMoveBookmarkRequest_WithInvalidMoveBookmarkID,
	"WithInvalidMoveParentID":     ValidateMoveBookmarkRequest_WithInvalidMoveParentID,
}

// <hook/>
func Test_ValidateMoveBookmarkRequest(t *testing.T) {
	fmt.Println("Test_ValidateMoveBookmarkRequest")
	for k, v := range Map_ValidateMoveBookmarkRequest {
		t.Run(k, v)
	}
}

// </method>
// </service>
