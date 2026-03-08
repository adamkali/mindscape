package services

import "github.com/adamkali/mindscape/cmd/configuration"

type Registrar struct {
	Config           *configuration.Configuration
	ApiKeyService    IApiKeyService
	AuthService      IAuthService
	BookmarkService  IBookmarkService
	FolderService    IFolderService
	MinioService     IMinioService
	NoteService      INoteService
	RedisService     IRedisService
	TaskService      ITaskService
	UserService      IUserService
	WidgetService    IWidgetService
	ValidatorService *ValidatorService
}
