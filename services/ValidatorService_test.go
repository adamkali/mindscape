package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	return httptest.NewRequest(
		method,
		path,
		strings.NewReader(string(requestJson)),
	)
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
	assert.Equal(t, r.IsAdmin, true)
}

func ValidateNewUserRequest_BadUsername(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Email and Username cannot be null",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadUsername_SQLInjection(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername_SQLInjection(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualErrorf(
		t,
		e,
		"Validation failed (%s) is not a valid username",
		"username'; DROP TABLE users;",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadUsername_SQLInjection_WithAdmin(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualErrorf(
		t,
		e,
		"Validation failed (%s) is not a valid username",
		"username'; DROP TABLE users;",
	)
	assert.Nil(t, r)
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
		"Validation faild. Seven Or More (false), Number (true), Upper (true), Special (true)",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoUpper(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoUpper(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation faild. Seven Or More (true), Number (true), Upper (false), Special (true)",
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
		"Validation faild. Seven Or More (true), Number (false), Upper (true), Special (true)",
	)
	assert.Nil(t, r)
}

func ValidateNewUserRequest_BadPassword_NoLower(t *testing.T) {
	r, e := Run_ValidateNewUserRequest_BadPassword_NoLower(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation faild. Seven Or More (true), Number (true), Upper (true), Special (true)",
	)
	assert.Nil(t, r)
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
		WithOnlyPassword(UpdateUserCredentialRequestJson()),
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
	assert.Equal(t, r.Password, "passwordABC123!")
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
	assert.EqualError(
		t,
		e,
		"ID cannot be null",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadEmail_SQLInjection(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Email cannot be null",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadUsername(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadUsername(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Email and Username cannot be null",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadUsername_SQLInjection(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed (adamkali) is not a valid username",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NotSevenOrMore(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (false), Upper (false), Special (false)",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoLower(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (false), Upper (false), Special (false)",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoUpper(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (false), Upper (false), Special (false)",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoNumber(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (false), Upper (false), Special (false)",
	)
	assert.Nil(t, r)
}
func ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial(t *testing.T) {
	r, e := Run_ValidateUpdateUserCredentialRequest_WithBadPassword_NoSpecial(services.ValidatorService{})
	assert.Error(t, e)
	assert.EqualError(
		t,
		e,
		"Validation failed. Seven Or More (false), Number (false), Upper (false), Special (false)",
	)
	assert.Nil(t, r)
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
