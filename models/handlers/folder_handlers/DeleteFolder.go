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

type DeleteFolderHandler struct {
	data          *repository.Folder
	err           error
	code          int
	ctx           echo.Context
	AuthService   services.IAuthService
	FolderService services.IFolderService
}

func NewDeleteHandler(
	ctx echo.Context,
	FolderService services.IFolderService,
	AuthService services.IAuthService,
) *DeleteFolderHandler {
	return &DeleteFolderHandler{
		ctx:           ctx,
		code:          200,
		FolderService: FolderService,
		AuthService:   AuthService,
	}
}

func (h *DeleteFolderHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	var folderID uuid.UUID
	if folderID, err = uuid.Parse(h.ctx.Param("folder_id")); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if h.data, err = h.FolderService.Get(folderID); err != nil {
		return handlers.Lock(h, 404, err)
	}
	if h.data.UserID != userID {
		return handlers.Lock(h, 403, fmt.Errorf("unauthorized access to folder"))
	}
	if err = h.FolderService.Remove(folderID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *DeleteFolderHandler) JSON() error {
	var message string
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	} else {
		message = "Folder deleted successfully"
		return responses.NewStringResponse().Successful(h.ctx, message)
	}
}

func (h *DeleteFolderHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *DeleteFolderHandler) Code() int {
	return h.code
}

func (h *DeleteFolderHandler) Data() any {
	return h.data
}

func (h *DeleteFolderHandler) Error() error {
	return h.err
}

func (h *DeleteFolderHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
