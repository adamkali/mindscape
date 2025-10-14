package controllers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/models/handlers/folder_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type FolderController struct {
	Name             string
	Config           *configuration.Configuration
	AuthService      services.IAuthService
	UserService      services.IUserService
	FolderService    services.IFolderService
	BookmarkService  services.IBookmarkService
	NoteService      services.INoteService
	ValidatorService *services.ValidatorService
}

func (uc FolderController) ControllerName() string {
	return uc.Name
}

func BuildFolderController(p *Registrar) FolderController {
	return FolderController{
		Name:             "/folders",
		Config:           p.Config,
		AuthService:      p.AuthService,
		UserService:      p.UserService,
		FolderService:    p.FolderService,
		BookmarkService:  p.BookmarkService,
		NoteService:      p.NoteService,
		ValidatorService: p.ValidatorService,
	}
}

// @Summary Get the Root Folders associated with the user
// @Description Get the Root Folders associated with the user by Authorization Header
// @Description and will also try to get the children of the folder as well
//
// @ID          GetRootFolders
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     responses.FoldersResponse
// @Failure     401                 {object}     responses.FoldersResponse
// @Failure     404                 {object}     responses.FoldersResponse
// @Failure     500                 {object}     responses.FoldersResponse
// @Router      /folders [get]
func (folderController FolderController) GetRootFolders(e echo.Context) error {
	return folder_handlers.NewGetRootHandler(
		e,
		folderController.FolderService,
		folderController.BookmarkService,
		folderController.NoteService,
		folderController.AuthService,
	).Handle().JSON()
}

// @Summary Get the Folders associated with the user under A Parent Folder
// @Description Get the Folders associated with the user under A Parent Folder by Authorization Header
// @Description and will also try to get the children of the folder as well
//
// @ID          GetFolders
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       folder_id           path         string                         true "Folder ID"                default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders/{folder_id} [get]
func (folderController FolderController) GetFolderByID(e echo.Context) error {
	return folder_handlers.NewGetById(
		e,
		folderController.FolderService,
		folderController.BookmarkService,
		folderController.NoteService,
		folderController.AuthService,
	).Handle().JSON()
}

// @Summary Create a new Folder
// @Description Create a new Folder by Authorization Header
//
// @ID          CreateFolder
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       CreateFolderRequest body         repository.CreateFolderParams  true "Create Folder Request"
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders [post]
func (folderController FolderController) CreateFolder(e echo.Context) error {
	return folder_handlers.NewCreateHandler(
		e,
		*folderController.ValidatorService,
		folderController.FolderService,
		folderController.AuthService,
	).Handle().JSON()
}

// @Summary Delete a Folder
// @Description Delete a Folder by Authorization Header and tries to cascade delete
//
// @ID          DeleteFolder
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       folder_id           path         string                         true "Folder ID"                default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders/{folder_id} [delete]
func (folderController FolderController) DeleteFolder(e echo.Context) error {
	return folder_handlers.NewDeleteHandler(
		e,
		folderController.FolderService,
		folderController.AuthService,
	).Handle().JSON()
}

// @Summary Move a Folder
// @Description Move a Folder by Authorization Header
//
// @ID          MoveFolder
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       MoveFolderRequest body         requests.MoveFolderRequest true "Move Folder Request"
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     403                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders [patch]
func (folderController FolderController) MoveFolder(e echo.Context) error {
	return folder_handlers.NewMoveHandler(
		e,
		*folderController.ValidatorService,
		folderController.FolderService,
		folderController.AuthService,
	).Handle().JSON()
}

func (folderController FolderController) Attatch(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	api := e.Group("/api" + folderController.Name)
	api.GET("", folderController.GetRootFolders, authMiddleware)
	api.GET("/:folder_id", folderController.GetFolderByID, authMiddleware)
	api.POST("", folderController.CreateFolder, authMiddleware)
	api.PATCH("", folderController.MoveFolder, authMiddleware)
	api.DELETE("/:folder_id", folderController.DeleteFolder, authMiddleware)
}
