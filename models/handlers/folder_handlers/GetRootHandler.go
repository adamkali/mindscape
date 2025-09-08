package folder_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type GetRootFolderHandler struct {
	Data            []responses.FolderData
	err             error
	code            int
	ctx             echo.Context
	AuthService     services.IAuthService
	FolderService   services.IFolderService
	BookmarkService services.IBookmarkService
	NoteService     services.INoteService
}

func NewGetRootHandler(
	ctx echo.Context,
	FolderService services.IFolderService,
	BookmarkService services.IBookmarkService,
	NoteService services.INoteService,
	AuthService services.IAuthService,
) *GetRootFolderHandler {
	return &GetRootFolderHandler{
		ctx:             ctx,
		code:            200,
		FolderService:   FolderService,
		BookmarkService: BookmarkService,
		NoteService:     NoteService,
		AuthService:     AuthService,
	}
}

func (h *GetRootFolderHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		handlers.Lock(h, 401, err)
	}
	folders := make([]repository.Folder, 0)
	if folders, err = h.FolderService.GetRoot(userID); err != nil {
		handlers.Lock(h, 404, err)
	}
	h.Data = make([]responses.FolderData, 0)
	for _, folder := range folders {
		folderData := responses.NewFolderData(folder)
		if folderData.Bookmarks, err = h.BookmarkService.GetByFolder(*folderData.ID); err != nil {
			handlers.Lock(h, 500, err)
		}
		if folderData.Notes, err = h.NoteService.GetByFolder(*folderData.ID); err != nil {
			handlers.Lock(h, 500, err)
		}
		if folderData.Children, err = h.FolderService.GetByParent(*folderData.ID); err != nil {
			handlers.Lock(h, 500, err)
		}
		h.Data = append(h.Data, folderData)
	}
	return h
}

func (h *GetRootFolderHandler) JSON() error {
	if h.err != nil {
		return responses.NewFoldersResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewFoldersResponse().Successful(h.ctx, h.Data)
}

func (h *GetRootFolderHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}


func (h *GetRootFolderHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
