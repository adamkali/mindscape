package handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/requests"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UpdateUserHandler struct {
	UserID            uuid.UUID
	UpdateUserRequest *requests.UpdateCredentialsRequest
	User              *repository.User
	Token             *string
	Context           echo.Context
	Error             error
	Code              int
	Locked            bool
}

func NewUpdateUserHandler(ctx echo.Context) *UpdateUserHandler {
	return &UpdateUserHandler{
		Context: ctx,
		Locked:  false,
		Error:   nil,
		Code:    200,
	}
}

func (h *UpdateUserHandler) Lock(code int) *UpdateUserHandler {
	h.Locked = true
	h.Code = code
	return h
}

// 1. get the user id from the jwt token
// 2. get the request from the context
// 3. check that the user id in the request matches the user id in the jwt token
// 4. update the user
// 5. refresh the jwt in the database
func (h *UpdateUserHandler) Handle(fun any) *UpdateUserHandler {
	var code int
	if !h.Locked {
		switch handle := fun.(type) {

		case func(token string) error:
			code = 401
			jwt_token := h.Context.Get("user").(*jwt.Token)
			claims := jwt_token.Claims.(*services.CustomJwt)
			h.UserID = claims.UserId
			h.Error = handle(jwt_token.Raw)
			if h.Error != nil {
				return h.Lock(code)
			}
			break

		case func(e echo.Context) (*requests.UpdateCredentialsRequest, error):
			code = 400
			h.UpdateUserRequest, h.Error = handle(h.Context)
			fmt.Printf("UpdateUserRequest: %v\n", h.UpdateUserRequest)
			if h.Error != nil {
				return h.Lock(code)
			}
			if h.UpdateUserRequest.ID != h.UserID {
				code = 403
				h.Error = echo.NewHTTPError(code, "Unauthorized")
				return h.Lock(code)
			}
			break

		case func(params *requests.UpdateCredentialsRequest) (*repository.User, error):
			code = 500
			h.User, h.Error = handle(h.UpdateUserRequest)
			if h.Error != nil {
				return h.Lock(code)
			}
			break

		case func(user repository.User) (*string, error):
			code = 500
			h.Token, h.Error = handle(*h.User)
			if h.Error != nil {
				return h.Lock(code)
			}
			break

		default:
			fmt.Printf("Type assertion failed for type: %T\n", fun)
			code = 600
			h.Error = echo.NewHTTPError(code, "Misaligned handler on the server")

		}
		if h.Error != nil {
			return h.Lock(code)
		}
	}
	return h
}

func (h *UpdateUserHandler) JSON() error {
	var code int
	var message string
	if h.Locked && h.Error != nil {
		code = h.Code
		if code == 600 {
			message = "Misaligend handler on the server"
		} else {
			message = h.Error.Error()
		}
	} else if code == 200 {
		message = "OK"
	}
	if h.Token == nil {
		h.Token = new(string)
	}
	return h.Context.JSON(code, responses.UpdateResponse{
		Message: message,
		Success: !h.Locked,
		Data:    responses.UserDataFromRepository(h.User),
		JWT:     *h.Token,
	})
}
