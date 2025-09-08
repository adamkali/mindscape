package bookmark_handlers

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

type GetByFolderHandler struct {
	Data         []repository.Bookmark
	err          error
	code         int
	ctx          echo.Context
	AuthService  services.IAuthService
	BookmarkService services.IBookmarkService
}

func NewGetFolderHandler(
	ctx echo.Context,
	BookmarkService services.IBookmarkService,
	AuthService services.IAuthService,
) *GetByFolderHandler {
	return &GetByFolderHandler{
		ctx:             ctx,
		code:            200,
		BookmarkService: BookmarkService,
		AuthService:     AuthService,
	}
}

func (h *GetByFolderHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		handlers.Lock(h, 401, err)
	}
	var parentID uuid.UUID
	if parentID, err = uuid.Parse(h.ctx.Param("parent_id")); err != nil {
		handlers.Lock(h, 400, err)
	}
	if h.Data, err = h.BookmarkService.GetByFolder(parentID); err != nil {
		handlers.Lock(h, 404, err)
	}
	for _, bookmark := range h.Data {
		if bookmark.UserID != userID {
			handlers.Lock(h, 403, fmt.Errorf("unauthorized access to bookmark"))
		}
	}
	return h
}

func (h *GetByFolderHandler) JSON() error {
	if h.err != nil {
		return responses.NewBookmarksResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewBookmarksResponse().Successful(h.ctx, h.Data)
}

func (h *GetByFolderHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *GetByFolderHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}