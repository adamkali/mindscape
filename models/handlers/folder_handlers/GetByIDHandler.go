package folder_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetFolderByIDHandler struct {
	Data            *responses.FolderData
	err             error
	code            int
	ctx             echo.Context
	AuthService     services.IAuthService
	FolderService   services.IFolderService
	BookmarkService services.IBookmarkService
	NoteService     services.INoteService
}

func NewGetById(
	ctx echo.Context,
	FolderService services.IFolderService,
	BookmarkService services.IBookmarkService,
	NoteService services.INoteService,
	AuthService services.IAuthService,
) *GetFolderByIDHandler {
	return &GetFolderByIDHandler{
		ctx:             ctx,
		code:            200,
		FolderService:   FolderService,
		BookmarkService: BookmarkService,
		NoteService:     NoteService,
		AuthService:     AuthService,
	}
}

func (h *GetFolderByIDHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		handlers.Lock(h, 401, err)
	}
	var folderID uuid.UUID
	if folderID, err = uuid.Parse(h.ctx.Param("folder_id")); err != nil {
		handlers.Lock(h, 400, err)
	}
	folder := new(repository.Folder)
	if folder, err = h.FolderService.Get(folderID); err != nil {
		handlers.Lock(h, 404, err)
	}
	if folder.UserID != userID {
		handlers.Lock(h, 403, fmt.Errorf("unauthorized access to folder"))
	}
	folderData := responses.NewFolderData(*folder)
	if folderData.Bookmarks, err = h.BookmarkService.GetByFolder(folderID); err != nil {
		handlers.Lock(h, 500, err)
	}
	if folderData.Notes, err = h.NoteService.GetByFolder(folderID); err != nil {
		handlers.Lock(h, 500, err)
	}
	if folderData.Children, err = h.FolderService.GetByParent(folderID); err != nil {
		handlers.Lock(h, 500, err)
	}
	return h
}

func (h *GetFolderByIDHandler) JSON() error {
	if h.err != nil {
		return responses.NewFolderResponse().Fail(h.ctx, h.code, h.err)
	}
	return h.ctx.JSON(200, responses.NewFolderResponseWithData(*h.Data, true, "OK"))
}

func (h *GetFolderByIDHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *GetFolderByIDHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
