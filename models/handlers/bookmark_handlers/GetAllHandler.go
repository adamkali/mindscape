package bookmark_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetByFolderHandler struct {
	UserID    uuid.UUID
	Bookmarks []repository.Bookmark
	Context   echo.Context
	Error     error
	Code      int
	Locked    bool
}

func NewGetFolderHandler(e echo.Context) *GetByFolderHandler {
	return &GetByFolderHandler{
		Context: e,
		Locked:  false,
		Code:    200,
		Error:   nil,
	}
}

func (h *GetByFolderHandler) Lock(code int) *GetByFolderHandler {
	h.Locked = true
	h.Code = code
	return h
}

func (h *GetByFolderHandler) Handle(fun any) *GetByFolderHandler {
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
		case func(uuid.UUID) ([]repository.Bookmark, error):
			code = 400
			var parentId uuid.UUID
			h.UserID, h.Error = uuid.Parse(h.Context.Param("parent_id"))
			if h.Error != nil {
				return h.Lock(code)
			}
			h.Bookmarks, h.Error = handle(parentId)
			for _, bookmark := range h.Bookmarks {
				if bookmark.UserID != h.UserID {
					return h.Lock(404)
				}
			}
			code = 404
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

func (h *GetByFolderHandler) JSON() error {
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
	return h.Context.JSON(code, responses.NewBookmarksResponse(
		h.Bookmarks,
		!h.Locked,
		message,
	))
}
