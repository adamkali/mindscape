// package user_handlers_test
//
// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"reflect"
// 	"testing"
//
// 	"github.com/adamkali/mindscape/models/handlers/user_handlers"
// 	"github.com/adamkali/mindscape/models/requests"
// 	"github.com/adamkali/mindscape/models/responses"
// 	"github.com/adamkali/mindscape/services"
// 	"github.com/labstack/echo/v4"
// )
//
//
// func TestNewRegisterHandler_Success_IsAdmin(t *testing.T) {
// 	userService := services.MockUserService{}
// 	authService := services.MockAuthService{}
// 	validatorService := services.ValidatorService{}
//
// 	// new http.Request
// 	request := new(http.Request)
// 	requestBody := requests.NewUserRequest{
// 		Username: "TestUser",
// 		Password: "AdminOfMe123!",
// 		Email:    "WlOwq@example.com",
// 		IsAdmin:  true,
// 	}
// 	requestBodyBytes, _ := json.Marshal(requestBody) 
// 	request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes)) 
//
// 	// new echo.Context
// 	context := echo.New().NewContext(request, nil)
// 	handler := user_handlers.NewRegisterHandler(
// 		context,
// 		validatorService,
// 		&userService,
// 		&authService,
// 	)
// 	if err:= handler.Handle(&requestBody).JSON(); err != nil {
// 		t.Error(err)
// 	}
//
// 	response := context.Response()
// }
package user_handlers_test
