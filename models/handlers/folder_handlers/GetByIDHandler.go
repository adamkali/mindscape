package folder_handlers

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)


type GetFolderByIDHandler struct {
	UserID          uuid.UUID
	FolderID        uuid.UUID
	FolderResponse responses.FolderData
	Context         echo.Context
	Error           error
	Code            int
	Locked          bool
}

func NewGetById(e echo.Context) *GetFolderByIDHandler {
	return &GetFolderByIDHandler{
		Context: e,
		Locked:  false,
		Error:   nil,
		Code:    200,
	}
}

func (grfh *GetFolderByIDHandler) Lock(code int) *GetFolderByIDHandler {
	grfh.Locked = true
	grfh.Code = code
	return grfh
}

func (grfh *GetFolderByIDHandler) Handle(fun any) *GetFolderByIDHandler {
	var code int
	if !grfh.Locked {
		switch handle := fun.(type) {
		case func(token string) error: // this is to the jwt token
			code = 401
			grfh.Error = handle(grfh.Context.Get("user").(*jwt.Token).Raw)
			if grfh.Error != nil {
				return grfh.Lock(code)
			}
			jwt_token := grfh.Context.Get("user").(*jwt.Token)
			claims := jwt_token.Claims.(*services.CustomJwt)
			grfh.UserID = claims.UserId
			break 

		case func(id uuid.UUID) (*repository.Folder, error):
			code = 400
			folder := new(repository.Folder) 
			grfh.FolderID, grfh.Error = uuid.Parse(grfh.Context.Param("folder_id"))
			if grfh.Error != nil {
				grfh.Error = echo.NewHTTPError(
					code,
					fmt.Sprintf("Type assertion failed for %s: %T\n", runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name(), fun),
				)
				return grfh.Lock(code)
				
			}
			code = 404
			folder, grfh.Error = handle(grfh.FolderID)
			grfh.FolderResponse = responses.NewFolderData(*folder)
			break

		case func(id uuid.UUID) ([]repository.Bookmark, error):
			code = 404
			grfh.FolderResponse.Bookmarks, grfh.Error = handle(grfh.FolderID)
			break

		case func(id uuid.UUID) ([]repository.Note, error):
			code = 404
			grfh.FolderResponse.Notes, grfh.Error = handle(grfh.FolderID)
			break

		case func(id uuid.UUID) ([]repository.Folder, error):
			code = 404
			grfh.FolderResponse.Children, grfh.Error = handle(grfh.FolderID)
			break


		default:
			code = 600
			grfh.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("Type assertion failed for %s: %T\n", runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name(), fun),
			)
		}
		if grfh.Error != nil {
			grfh.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("GetFolderByIDHandler.%s caused: %s", runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name(), grfh.Error.Error()),
			)
			return grfh.Lock(code)
		}
	}
	return grfh
}

func (h *GetFolderByIDHandler) JSON() error {
	var code int
	var message string
	if h.Locked && h.Error != nil {
		code = h.Code
		if code == 600 {
			message = fmt.Sprintf("Misaligend handler on the server: %s", h.Error.Error())
		} else {
			message = h.Error.Error()
		}
	} else if code == 200 {
		message = "OK"
	}

	return h.Context.JSON(code,
	    responses.NewFolderResponseWithData(
		    h.FolderResponse,
		    !h.Locked,
		    message,
	    ),
	)
}
