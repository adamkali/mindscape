package controllers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/models/handlers/bookmark_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type BookmarkController struct {
	Name             string
	Config           *configuration.Configuration
	AuthService      services.IAuthService
	UserService      services.IUserService
	FolderService    services.IFolderService
	BookmarkService  services.IBookmarkService
	NoteService      services.INoteService
	ValidatorService *services.ValidatorService
}

func (uc BookmarkController) ControllerName() string {
	return uc.Name	
}

func BuildBookmarkController(p *services.Registrar) BookmarkController {
	return BookmarkController {
		Name:             "/bookmarks",
		Config:           p.Config,
		AuthService:      p.AuthService,
		UserService:      p.UserService,
		FolderService:    p.FolderService,
		BookmarkService:  p.BookmarkService,
		NoteService:      p.NoteService,
		ValidatorService: p.ValidatorService,
	}
}

// @Summary Create a new Bookmark  
// @Description Create a new Bookmark by Authorization Header 
//
// @ID          CreateBookmark
// @Tags        Bookmarks 
// @Accept      json
// @Produce     json
// @Param       CreateBookmarkRequest body         repository.CreateBookmarkParams  true "CreateBookmarkRequest"
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     BookmarkResponse
// @Failure     401                 {object}     BookmarkResponse
// @Failure     404                 {object}     BookmarkResponse
// @Failure     500                 {object}     BookmarkResponse
// @Router      /bookmarks [post]
func (c BookmarkController) Create(e echo.Context) error {
	return bookmark_handlers.NewCreateHandler(
		e,
		*c.ValidatorService,
		c.BookmarkService,
		c.AuthService,
	).
		Handle().JSON()
}

// @Summary Get Bookmarks By Folder ID
// @Description Get all Bookmarks by Authorization Header and by the 
// @Description ParentFolderId [parent_id]
//
// @ID          GetBookmarks
// @Tags        Bookmarks
// @Accept      json
// @Produce     json
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Param       parent_id           path         string                         true "Parent Folder ID"                default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     BookmarksResponse
// @Failure     401                 {object}     BookmarksResponse
// @Failure     404                 {object}     BookmarksResponse
// @Failure     500                 {object}     BookmarksResponse
// @Router      /bookmarks/folder/{parent_id} [get]
func (c BookmarkController) GetByFolder(e echo.Context) error {
	return bookmark_handlers.NewGetFolderHandler(
		e,
		c.BookmarkService,
		c.AuthService,
	).Handle().JSON()
}

// @Summary     Delete a Bookmark
// @Description Delete a Bookmark by Authorization Header, and my a 
// @Description ParentFolderId [parent_id]. 
//
// @ID          DeleteBookmark 
// @Tags        Bookmarks
// @Accept      json
// @Produce     json
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Param       bookmark_id         path         string                         true "Parent Folder ID"         default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     BookmarksResponse
// @Failure     401                 {object}     BookmarksResponse
// @Failure     404                 {object}     BookmarksResponse
// @Failure     500                 {object}     BookmarksResponse
// @Router      /bookmarks/folder/{bookmark_id} [delete]
func (c BookmarkController) Delete(e echo.Context) error {
	return bookmark_handlers.NewDeleteHandler(
		e,
		c.BookmarkService,
		c.AuthService,
	).Handle().JSON()
}

// @Summary     Move a Bookmark
// @Description Move a Bookmark by Authorization Header, and my a 
// @Description ParentFolderId [parent_id]. 
//
// @ID          MoveBookmark 
// @Tags        Bookmarks
// @Accept      json
// @Produce     json
// @Param       MoveBookmarkRequest body         requests.MoveBookmarkRequest true "Move Bookmark Request"
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Success     200                 {object}     BookmarksResponse
// @Failure     401                 {object}     BookmarksResponse
// @Failure     404                 {object}     BookmarksResponse
// @Failure     500                 {object}     BookmarksResponse
// @Router      /bookmarks [patch]
func (c BookmarkController) Move(e echo.Context) error {
	return bookmark_handlers.NewMoveHandler(
		e,
		*c.ValidatorService,
		c.BookmarkService,
		c.AuthService,
	).Handle().JSON()
}

// @Summary     Update a Bookmark
// @Description Update a Bookmark's name and link by Authorization Header
//
// @ID          UpdateBookmark
// @Tags        Bookmarks
// @Accept      json
// @Produce     json
// @Param       UpdateBookmarkRequest body         requests.UpdateBookmarkRequest true "Update Bookmark Request"
// @Param       Authorization       header       string                         true "Authorization Header"     default("Bearer token")
// @Param       bookmark_id         path         string                         true "Bookmark ID"              default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     BookmarkResponse
// @Failure     400                 {object}     BookmarkResponse
// @Failure     401                 {object}     BookmarkResponse
// @Failure     403                 {object}     BookmarkResponse
// @Failure     404                 {object}     BookmarkResponse
// @Failure     500                 {object}     BookmarkResponse
// @Router      /bookmarks/{bookmark_id} [put]
func (c BookmarkController) UpdateBookmark(e echo.Context) error {
	return bookmark_handlers.NewUpdateHandler(
		e,
		*c.ValidatorService,
		c.BookmarkService,
		c.AuthService,
	).Handle().JSON()
}

func (c BookmarkController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api" + c.Name)
	api.POST("", c.Create, middlewares...)
	api.GET("/folder/:parent_id", c.GetByFolder, middlewares...)
	api.PUT("/:bookmark_id", c.UpdateBookmark, middlewares...)
	api.PATCH("", c.Move, middlewares...)
	api.DELETE("/folder/:bookmark_id", c.Delete, middlewares...)
}
