package folder_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateFolderHandler struct {
	UserID          uuid.UUID
	Params          *repository.CreateFolderParams
	FolderResponse responses.FolderData
	Context         echo.Context
	Error           error
	Code            int
	Locked          bool
}

func NewCreateHandler(context echo.Context) *CreateFolderHandler {
	return &CreateFolderHandler{
		Context: context,
		Locked:  false,
		Code:    200,
	}
}

func (h *CreateFolderHandler) Lock(code int) *CreateFolderHandler {
	h.Locked = true
	h.Code = code
	return h
}
	
func (h *CreateFolderHandler) Handle(fun any) *CreateFolderHandler {
	var code int
	if !h.Locked {
		switch handle := fun.(type) {
		case func(token string) error: // this is to the jwt token
			code = 401
			h.Error = handle(h.Context.Get("user").(*jwt.Token).Raw)
			if h.Error != nil {
				return h.Lock(code)
			}
			jwt_token := h.Context.Get("user").(*jwt.Token)
			claims := jwt_token.Claims.(*services.CustomJwt)
			h.UserID = claims.UserId
			break
		case func(echo.Context) (*repository.CreateFolderParams, error):
			code = 400
			h.Params, h.Error = handle(h.Context)
			if h.Error != nil {
				return h.Lock(code)
			}
			break
		case func(*repository.CreateFolderParams) (*repository.Folder, error):
			code = 500
			folder := new(repository.Folder)
			folder, h.Error = handle(h.Params)
			if h.Error != nil {
				return h.Lock(code)
			}
			h.FolderResponse = responses.NewFolderData(*folder)
			break
		default:
			code = 600
			h.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("Type assertion failed for type: %T\n", fun),
			)
		}
		if h.Error != nil {
			return h.Lock(code)
		}
	}
	return h
}

func (h *CreateFolderHandler) Response() error {
	if h.Locked && h.Error != nil {
		return h.JSON()
	} else {
		return h.JSON()
	}
}

func (h *CreateFolderHandler) JSON() error {
	var code int
	var message string
	if h.Locked && h.Error != nil {
		code = h.Code
		if code == 600 {
			message = "Misaligend handler on the server: " + h.Error.Error()
		} else {
			message = h.Error.Error()
		}
	} else if code == 200 {
		message = "OK"
	}
	response := responses.FolderResponse{
		Data:    h.FolderResponse,
		Message: message,
		Success: !h.Locked,
	}

	return h.Context.JSON(code, response)
}


