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

type DeleteFolderHandler struct {
	UserID          uuid.UUID
	FolderID        uuid.UUID
	FolderData      responses.FolderData
	Context         echo.Context
	Error           error
	Code            int
	Locked          bool
}

func NewDeleteHandler(context echo.Context) *DeleteFolderHandler {
	return &DeleteFolderHandler{
		Context: context,
		Locked:  false,
		Code:    200,
	}
}

func (h *DeleteFolderHandler) Lock(code int) *DeleteFolderHandler {
	h.Code = code
	h.Locked = true
	return h
}

func (h *DeleteFolderHandler) Response() error {
	if h.Locked && h.Error != nil {
		return h.JSON()
	} else {
		return h.JSON()
	}
}

func (h *DeleteFolderHandler) JSON() error {
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

	return h.Context.JSON(code, responses.FolderResponse{
		Data:    h.FolderData,
		Message: message,
		Success: !h.Locked,
	})
}

func (h *DeleteFolderHandler) Handle(fun any) *DeleteFolderHandler {
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

		case func(folderID uuid.UUID) (*repository.Folder, error):
			code = 404
			folder := new(repository.Folder)
			folder, h.Error = handle(h.FolderID)
			if h.Error != nil {
				return h.Lock(code)	
			}
			h.FolderData = responses.NewFolderData(*folder)
			break
		case func(folderID uuid.UUID) error:
			code = 500
			h.Error = handle(h.FolderID)
			if h.Error != nil {
				return h.Lock(code)
			}
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
