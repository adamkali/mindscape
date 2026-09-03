package folder_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type ShareFolderWithHouseHoldMemersHandler struct {
	result   bool
	err      error
	code     int
	ctx      echo.Context
	services *services.Registrar
}

func NewShareFolderWithHouseHoldMemersHandler(
	ctx echo.Context,
	services *services.Registrar,
) *ShareFolderWithHouseHoldMemersHandler {
	return &ShareFolderWithHouseHoldMemersHandler{
		ctx:      ctx,
		services: services,
	}
}

func ShareFolderWithHouseHoldMemersHandlerJsonHandler(
	ctx echo.Context,
	services *services.Registrar,
) error {
	return (NewShareFolderWithHouseHoldMemersHandler(ctx, services)).Handle().JSON()
}

func (h *ShareFolderWithHouseHoldMemersHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}
func (h *ShareFolderWithHouseHoldMemersHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}
func (h *ShareFolderWithHouseHoldMemersHandler) Code() int    { return h.code }
func (h *ShareFolderWithHouseHoldMemersHandler) Error() error { return h.err }
func (h *ShareFolderWithHouseHoldMemersHandler) Data() any    { return h.result }

func (h *ShareFolderWithHouseHoldMemersHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewBoolResponse().Successful(h.ctx)
}

// shareFolders clones the folder tree rooted at `folder` for sharedUser.
// It must create top-down: a folder's clone has to exist before its
// children can be created pointing at it as their new parent, so this is a
// straightforward parent-first recursion rather than a flat/order-independent
// walk. newParentID is nil for the root clone (no parent) and the previously
// created clone's ID for every folder below it.
//
// TODO: this is not transactional — FolderService.Create/BookmarkService.Create
// each commit independently, so a failure partway through the recursion leaves
// orphaned clones from folders/bookmarks already created. Fix by adding a
// shared *pgxpool.Pool to services.Registrar, having the handler Begin() one
// pgx.Tx, and rewriting this to operate on a *repository.Queries (via
// repository.New(tx)) instead of going through FolderService/BookmarkService,
// committing once at the top level.
func shareFolders(
	services *services.Registrar,
	sharedUser uuid.UUID,
	folder repository.Folder,
	newParentID *uuid.UUID,
) (repository.Folder, error) {
	parentID := pgtype.UUID{Valid: false}
	if newParentID != nil {
		parentID = pgtype.UUID{Bytes: *newParentID, Valid: true}
	}
	clonedFolder, err := services.FolderService.Create(&repository.CreateFolderParams{
		UserID:      sharedUser,
		ParentID:    parentID,
		Name:        folder.Name,
		Description: folder.Description,
	})
	if err != nil {
		return repository.Folder{}, err
	}

	bookmarks, err := services.BookmarkService.GetByFolder(folder.ID)
	if err != nil {
		return repository.Folder{}, err
	}
	for _, bookmark := range bookmarks {
		if _, err := services.BookmarkService.Create(&repository.CreateBookmarkParams{
			UserID:   sharedUser,
			FolderID: clonedFolder.ID,
			Name:     bookmark.Name,
			Link:     bookmark.Link,
		}); err != nil {
			return repository.Folder{}, err
		}
	}

	children, err := services.FolderService.GetByParent(folder.ID)
	if err != nil {
		return repository.Folder{}, err
	}
	for _, child := range children {
		if _, err := shareFolders(services, sharedUser, child, &clonedFolder.ID); err != nil {
			return repository.Folder{}, err
		}
	}

	return *clonedFolder, nil
}

func (h *ShareFolderWithHouseHoldMemersHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	request, err := h.services.ValidatorService.ShareFolderWithHouseholdMemeberRequestValidator(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	// check if the request memberId is a member of the household
	// if not this shold return a 403 error with description
	// this user is not in your household, and for security reasons dont share until they are
	_, err = h.services.HouseholdService.GetMember(request.HouseholdID, request.MemberID)
	if err != nil {
		return handlers.Lock(h, 403, err)
	}
	fmt.Printf("[INFO] ShareFolderWithHouseHoldMemersHandler.Handle{ request: %v }\n", request)

	// now get the folder structure from userID and request.FolderID
	folder, err := h.services.FolderService.Get(request.FolderID)
	if err != nil {
		return handlers.Lock(h, 404, err)
	}
	// check thatt the folder is yours
	if folder.UserID != userID {
		return handlers.Lock(h, 403, fmt.Errorf("Unauthorized folder access."))
	}

	// share it with the memberID: clone the folder tree top-down so each
	// folder's new parent already exists before its children are created
	if _, err := shareFolders(h.services, request.MemberID, *folder, nil); err != nil {
		return handlers.Lock(h, 500, err)
	}

	h.result = true
	return h
}
