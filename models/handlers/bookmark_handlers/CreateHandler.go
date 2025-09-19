package bookmark_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CreateHandler struct {
	data             *repository.Bookmark
	err              error
	code             int
	ctx              echo.Context
	AuthService      services.IAuthService
	BookmarkService  services.IBookmarkService
	ValidatorService services.ValidatorService
}

func NewCreateHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	BookmarkService services.IBookmarkService,
	AuthService services.IAuthService,
) *CreateHandler {
	return &CreateHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		BookmarkService:  BookmarkService,
		AuthService:      AuthService,
	}
}

func (h *CreateHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	request := new(repository.CreateBookmarkParams)
	if request, err = h.ValidatorService.CreateBookmarkRequest(h.ctx); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if h.data, err = h.BookmarkService.Create(request); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *CreateHandler) JSON() error {
	if h.err != nil {
		return responses.NewBookmarkResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewBookmarkResponse().Successful(h.ctx, *h.data)
}

func (h *CreateHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *CreateHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *CreateHandler) Data() any {
	return h.data
}

func (h *CreateHandler) Code() int {
	return h.code
}

func (h *CreateHandler) Error() error {
	return h.err
}
