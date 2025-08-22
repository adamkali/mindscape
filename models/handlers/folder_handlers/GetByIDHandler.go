package folder_handlers

import (
	"fmt"
	"sync"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)


type GetFolderByIDHandler struct {
	UserID          uuid.UUID
	FolderResponses []responses.FolderData
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
		case func(uuid.UUID) ([]repository.Folder, error):
			code = 400
			folders := make([]repository.Folder, 0)
			grfh.FolderResponses = make([]responses.FolderData, len(folders))
			id, err := uuid.Parse(grfh.Context.Param("id"))
			if err != nil {
				grfh.Error = err
				return grfh.Lock(code)
			}
			code = 404
			folders, grfh.Error = handle(id)
			if grfh.Error != nil {
				return grfh.Lock(code)
			}
			for i, folder := range folders {
				grfh.FolderResponses[i] = responses.NewFolderData(folder)
			}
			break
		case func(uuid.UUID) ([]repository.Bookmark, error):
			code = 404
			errors := make([]error, len(grfh.FolderResponses))
			var wg sync.WaitGroup
			wg.Add(len(grfh.FolderResponses))
			// go func to create the folder responses
			bookmakData := func(folderData *responses.FolderData, er error) {
				folderData.Bookmarks, er = handle(*folderData.ID)
				wg.Done()
			}
			for i := range grfh.FolderResponses {
				go bookmakData(&grfh.FolderResponses[i], errors[i])
			}
			wg.Wait()
			for _, err := range errors {
				if err != nil {
					grfh.Error = err
					return grfh.Lock(code)
				}
			}
			break

		case func(uuid.UUID) ([]repository.Note, error):
			code = 404
			errors := make([]error, len(grfh.FolderResponses))
			var wg sync.WaitGroup
			wg.Add(len(grfh.FolderResponses))
			notesData := func(folderData *responses.FolderData, er error) {
				folderData.Notes, er = handle(*folderData.ID)
				wg.Done()
			}
			for i := range grfh.FolderResponses {
				go notesData(&grfh.FolderResponses[i], errors[i])
			}
			wg.Wait()
			for _, err := range errors {
				if err != nil {
					grfh.Error = err
					return grfh.Lock(code)
				}
			}
			break

		default:
			code = 600
			grfh.Error = echo.NewHTTPError(
				code,
				fmt.Sprintf("Type assertion failed for type: %T\n", fun),
			)
		}
		if grfh.Error != nil {
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
		responses.NewFoldersResponse(
			h.FolderResponses,
			!h.Locked,
			message,
		))
}
