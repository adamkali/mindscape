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
	Services         *services.Registrar
}

func (uc FolderController) ControllerName() string {
	return uc.Name
}

func BuildFolderController(p *services.Registrar) FolderController {
	return FolderController{
		Name:             "/folders",
		Config:           p.Config,
		Services:         p,
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
// @Security BearerAuth
// @Success     200                 {object}     responses.FoldersResponse
// @Failure     401                 {object}     responses.FoldersResponse
// @Failure     404                 {object}     responses.FoldersResponse
// @Failure     500                 {object}     responses.FoldersResponse
// @Router      /folders [get]
func (folderController FolderController) GetRootFolders(e echo.Context) error {
	return folder_handlers.NewGetRootHandler(
		e,
		folderController.Services.FolderService,
		folderController.Services.BookmarkService,
		folderController.Services.AuthService,
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
// @Security BearerAuth
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders/{folder_id} [get]
func (folderController FolderController) GetFolderByID(e echo.Context) error {
	return folder_handlers.NewGetById(
		e,
		folderController.Services.FolderService,
		folderController.Services.BookmarkService,
		folderController.Services.AuthService,
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
// @Security BearerAuth
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders [post]
func (folderController FolderController) CreateFolder(e echo.Context) error {
	return folder_handlers.NewCreateHandler(
		e,
		*folderController.Services.ValidatorService,
		folderController.Services.FolderService,
		folderController.Services.AuthService,
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
// @Security BearerAuth
// @Success     200                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders/{folder_id} [delete]
func (folderController FolderController) DeleteFolder(e echo.Context) error {
	return folder_handlers.NewDeleteHandler(
		e,
		folderController.Services.FolderService,
		folderController.Services.AuthService,
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
// @Security BearerAuth
// @Success     200                 {object}     responses.FolderResponse
// @Failure     401                 {object}     responses.FolderResponse
// @Failure     403                 {object}     responses.FolderResponse
// @Failure     404                 {object}     responses.FolderResponse
// @Failure     500                 {object}     responses.FolderResponse
// @Router      /folders [patch]
func (folderController FolderController) MoveFolder(e echo.Context) error {
	return folder_handlers.NewMoveHandler(
		e,
		*folderController.Services.ValidatorService,
		folderController.Services.FolderService,
		folderController.Services.AuthService,
	).Handle().JSON()
}

// @Summary Update a Folder
// @Description Update a Folder's name and description by Authorization Header
//
// @ID          UpdateFolder
// @Tags        Folders
// @Accept      json
// @Produce     json
// @Param       folder_id              path         string                              true "Folder ID"                default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       UpdateFolderRequest    body         requests.UpdateFolderRequest        true "Update Folder Request"
// @Security BearerAuth
// @Success     200                    {object}     responses.FolderResponse
// @Failure     400                    {object}     responses.FolderResponse
// @Failure     401                    {object}     responses.FolderResponse
// @Failure     403                    {object}     responses.FolderResponse
// @Failure     404                    {object}     responses.FolderResponse
// @Failure     500                    {object}     responses.FolderResponse
// @Router      /folders/{folder_id} [put]
func (folderController FolderController) UpdateFolder(e echo.Context) error {
	return folder_handlers.NewUpdateHandler(
		e,
		*folderController.Services.ValidatorService,
		folderController.Services.FolderService,
		folderController.Services.AuthService,
	).Handle().JSON()
}

func (folderController FolderController) ShareFolderWithHousehold(e echo.Context) error {
	return folder_handlers.ShareFolderWithHouseHoldMemersHandlerJsonHandler(e, folderController.Services)
}

func (folderController FolderController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api" + folderController.Name)
	api.GET("", folderController.GetRootFolders, middlewares...)
	api.GET("/:folder_id", folderController.GetFolderByID, middlewares...)
	api.POST("", folderController.CreateFolder, middlewares...)
	api.PATCH("", folderController.MoveFolder, middlewares...)
	api.PUT("/:folder_id", folderController.UpdateFolder, middlewares...)
	api.DELETE("/:folder_id", folderController.DeleteFolder, middlewares...)
}
