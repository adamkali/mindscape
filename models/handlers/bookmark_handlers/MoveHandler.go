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

type MoveHandler struct {
	data            *repository.Bookmark 
	err             error
	code            int
	ctx             echo.Context
	AuthService     services.IAuthService
	BookmarkService services.IBookmarkService
	ValidatorService services.ValidatorService
}

func NewMoveHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	BookmarkService services.IBookmarkService,
	AuthService services.IAuthService,
) *MoveHandler {
	return &MoveHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		BookmarkService:  BookmarkService,
		AuthService:      AuthService,
	}
}


func (h MoveHandler) JSON() error {
	if h.err != nil {
		return responses.NewBookmarkResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewBookmarkResponse().Successful(h.ctx, *h.data)
}
func (h MoveHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}
func (h MoveHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
func (h MoveHandler) Code() int {
	return h.code
}
func (h MoveHandler) Error() error {
	return h.err
}
func (h MoveHandler) Data() any {
	return h.data
}

func (h MoveHandler) Handle() handlers.IHandler {
	fmt.Printf("[DEBUG] MoveHandler.Handle Start\n")
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
    fmt.Printf("[INFO] MoveHandler.Handle{ jwt_token: %v }\n", jwt_token)
	request, err := h.ValidatorService.ValidateMoveBookmarkRequest(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	fmt.Printf("[INFO] MoveHandler.Handle{ request: %v }\n", request)
	userID := claims.UserId
	if h.data, err = h.BookmarkService.Get(request.BookmarkID); err != nil {
		return handlers.Lock(h, 404, err)
	}
	fmt.Printf("[INFO] MoveHandler.Handle{ h.data: %v }\n", h.data)
	if userID != request.UserID && userID != h.data.UserID {
		return handlers.Lock(h, 403, fmt.Errorf("Unauthorized folder access."))
	}
	fmt.Printf("[INFO] MoveHandler.Handle{ userID: %v }\n", userID)
	if h.data, err = h.BookmarkService.Move(request.Into()); err != nil {
		return handlers.Lock(h, 500, err)
	}
	fmt.Printf("[INFO] MoveHandler.Handle{ h.data: %v }\n", h.data)
	return h
}
