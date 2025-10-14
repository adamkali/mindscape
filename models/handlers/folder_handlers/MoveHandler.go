package folder_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// MoveHandler is a structure holding the services needed to
// preform moving a folder (MoveHandler.data) from one position
// to another.
type MoveHandler struct {
	data             *responses.FolderData
	err              error
	code             int
	ctx              echo.Context
	AuthService      services.IAuthService
	FolderService    services.IFolderService
	ValidatorService services.ValidatorService
}

func NewMoveHandler(
	ctx echo.Context,
	ValidatorService services.ValidatorService,
	FolderService services.IFolderService,
	AuthService services.IAuthService,
) *MoveHandler {
	return &MoveHandler{
		ctx:              ctx,
		code:             200,
		ValidatorService: ValidatorService,
		FolderService:    FolderService,
		AuthService:      AuthService,
	}
}

func (h MoveHandler) JSON() error {
	if h.err != nil {
		return responses.NewFolderResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewFolderResponse().SuccessfulWithData(h.ctx, *h.data)
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
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	var err error
	if err = h.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	request, err := h.ValidatorService.ValidateMoveFolderRequest(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	userID := claims.UserId
	folder := new(repository.Folder)
	fmt.Printf("[INFO] MoveHandler.Handle{ request: %v }\n", request)
	if folder, err = h.FolderService.Get(request.FolderID); err != nil {
		return handlers.Lock(h, 404, err)
	}
	if userID != request.UserID && userID != folder.UserID {
		return handlers.Lock(h, 403, fmt.Errorf("Unauthorized folder access."))
	}
	if err = h.FolderService.Move(folder.ID, request.NewParentID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	if request.NewParentID == nil {
		folder.ParentID = pgtype.UUID{
			Valid: true,
		}
	} else {
		folder.ParentID = pgtype.UUID{
			Bytes: *request.NewParentID,
			Valid: true,
		}
	}
	data := responses.NewFolderData(*folder)
	h.data = &data
	return h
}
