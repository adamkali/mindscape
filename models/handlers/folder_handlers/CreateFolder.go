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

type CreateFolderHandler struct {
	Data             *repository.Folder
	err              error
	code             int
	ctx              echo.Context
	AuthService      services.IAuthService
	FolderService    services.IFolderService
	ValidatorService services.ValidatorService
}

func NewCreateHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	FolderService services.IFolderService,
	AuthService services.IAuthService,
) *CreateFolderHandler {
	return &CreateFolderHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		FolderService:    FolderService,
		AuthService:      AuthService,
	}
}

func (h *CreateFolderHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	request, err := h.ValidatorService.ValidateCreateFolderRequest(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	fmt.Printf("[INFO] CreateFolderHandler.Handle{ request: %v }\n", request)
	params := &repository.CreateFolderParams{
		UserID:   userID,
		Name:     request.Name,
		ParentID: request.ParentID,
	}
	if h.Data, err = h.FolderService.Create(params); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *CreateFolderHandler) JSON() error {
	if h.err != nil {
		return responses.NewFolderResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewFolderResponse().Successful(h.ctx, *h.Data)
}

func (h *CreateFolderHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *CreateFolderHandler) SetError(err error) handlers.IHandler {
	h.err = fmt.Errorf("%d Error: %s", h.code, err.Error())
	return h
}
