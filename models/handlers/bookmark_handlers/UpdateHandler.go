package bookmark_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UpdateHandler struct {
	data             *repository.Bookmark
	err              error
	code             int
	ctx              echo.Context
	AuthService      services.IAuthService
	BookmarkService  services.IBookmarkService
	ValidatorService services.ValidatorService
}

func NewUpdateHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	BookmarkService services.IBookmarkService,
	AuthService services.IAuthService,
) *UpdateHandler {
	return &UpdateHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		BookmarkService:  BookmarkService,
		AuthService:      AuthService,
	}
}

func (h UpdateHandler) JSON() error {
	if h.err != nil {
		return responses.NewBookmarkResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewBookmarkResponse().Successful(h.ctx, *h.data)
}
func (h UpdateHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}
func (h UpdateHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
func (h UpdateHandler) Code() int {
	return h.code
}
func (h UpdateHandler) Error() error {
	return h.err
}
func (h UpdateHandler) Data() any {
	return h.data
}

func (h UpdateHandler) Handle() handlers.IHandler {
	fmt.Printf("[DEBUG] UpdateHandler.Handle Start\n")
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}

	request, err := h.ValidatorService.ValidateUpdateBookmarkRequest(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	userID := claims.UserId
	if h.data, err = h.BookmarkService.Get(request.BookmarkID); err != nil {
		return handlers.Lock(h, 404, err)
	}

	if userID != request.UserID && userID != h.data.UserID {
		return handlers.Lock(h, 403, fmt.Errorf("Unauthorized bookmark access."))
	}

	// Keep existing folder ID
	folderID := h.data.FolderID
	if h.data, err = h.BookmarkService.Update(request.BookmarkID, request.Into(folderID)); err != nil {
		return handlers.Lock(h, 500, err)
	}

	return h
}
