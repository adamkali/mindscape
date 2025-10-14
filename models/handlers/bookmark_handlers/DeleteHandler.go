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

type DeleteHandler struct {
	ctx             echo.Context
	code            int
	data            *repository.Bookmark
	err             error
	AuthService     services.IAuthService
	BookmarkService services.IBookmarkService
}
func NewDeleteHandler(
	ctx echo.Context,
	BookmarkService services.IBookmarkService,
	AuthService services.IAuthService,
) *DeleteHandler {
	return &DeleteHandler{
		ctx:             ctx,
		code:            200,
		BookmarkService: BookmarkService,
		AuthService:     AuthService,
	}
}

func (h *DeleteHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *DeleteHandler) Code() int {
	return h.code
}

func (h *DeleteHandler) Data() any {
	return h.data
}

func (h *DeleteHandler) Error() error {
	return h.err
}

func (h *DeleteHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}

func (h *DeleteHandler) JSON() error {
	var message string
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		message = "Bookmark deleted successfully"
		return responses.NewStringResponse().Successful(h.ctx, message)
	}
}

func (h *DeleteHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	var bookmarkID uuid.UUID
	bookmarkStringID := h.ctx.Param("bookmark_id")
	fmt.Printf("[INFO] bookmarkStringID: %v\n", bookmarkStringID)
	if bookmarkID, err = uuid.Parse(bookmarkStringID); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if h.data, err = h.BookmarkService.Get(bookmarkID); err != nil {
		return handlers.Lock(h, 404, err)
	}
	if h.data.UserID != userID {
		return handlers.Lock(h, 403, fmt.Errorf("unauthorized access to bookmark"))
	}
	if err = h.BookmarkService.Remove(bookmarkID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}


