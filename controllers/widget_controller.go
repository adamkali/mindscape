package controllers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/db/repository"
	handlers "github.com/adamkali/mindscape/models/handlers/widget_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

func (uc WidgetController) ControllerName() string {
	return uc.Name
}

func BuildWidgetController(p *services.Registrar) WidgetController {
	return WidgetController{
		Name:             "/widgets",
		Config:           p.Config,
		AuthService:      p.AuthService,
		MinioService:     p.MinioService,
		WidgetService:    p.WidgetService,
		RedisService:     p.RedisService,
		ValidatorService: p.ValidatorService,
	}
}

type WidgetController struct {
	Name             string
	Config           *configuration.Configuration
	AuthService      services.IAuthService
	WidgetService    services.IWidgetService
	MinioService     services.IMinioService
	RedisService     services.IRedisService
	ValidatorService *services.ValidatorService
}

// @Summary Get Widget Schemas
// @Description Get all embeded widget schemas from the Schema
// @Description Storage
//
// @ID          GetWidgetSchemas
// @Tags        Widgets
// @Produce     json
// @Success     200                 {object}     responses.WidgetsResponse
// @Failure     404                 {object}     responses.WidgetsResponse
// @Router      /widgets/schemas    [get]
func (wc WidgetController) ReadSchemas(
	ctx echo.Context,
) error {
	return handlers.NewReadHandler(
		ctx,
		wc.WidgetService,
	).Handle().JSON()
}

// @Summary Get Widget Schema by ID
// @Description Get embeded widget schema from the Schema
// @Description Storage by its identifier
//
// @ID          GetWidgetSchemaByID
// @Tags        Widgets
// @Produce     json
// @Param       schema_id           path         string                         true "Widget Schema Id"
// @Success     200                 {object}     responses.WidgetResponse
// @Failure     400                 {object}     responses.WidgetResponse
// @Failure     404                 {object}     responses.WidgetResponse
// @Router      /widgets/schemas/{schema_id}    [get]
func (wc WidgetController) GetSchemaByID(
	ctx echo.Context,
) error {
	return handlers.NewGetByIdHandler(
		ctx,
		wc.WidgetService,
	).Handle().JSON()
}

// @Summary Get a Users Widgets
// @Description Get a Users Widgets by their auth token
// @Description and return the list of widgets associated
// @Description with the user account in the request params.
//
// @ID          GetUserWidgets
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Success     200                 {object}     responses.UserWidgetsResponse
// @Failure     401                 {object}     responses.UserWidgetsResponse
// @Failure     403                 {object}     responses.UserWidgetsResponse
// @Failure     500                 {object}     responses.UserWidgetsResponse
// @Router      /widgets            [get]
func (wc WidgetController) Read(ctx echo.Context) error {
	return handlers.NewReadUserWidgetsHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

// @Summary Get a Users Widget
// @Description Get a Users Widget by their auth token
// @Description and a path parameter and return the
// @Description widget from the database.
//
// @ID          GetUserWidget
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     responses.UserWidgetResponse
// @Failure     401                 {object}     responses.UserWidgetResponse
// @Failure     403                 {object}     responses.UserWidgetResponse
// @Failure     500                 {object}     responses.UserWidgetResponse
// @Router      /widgets/{user_widget_id}    [get]
func (wc WidgetController) ReadById(ctx echo.Context) error {
	return handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

// @Summary Get a Users Github Widget
// @Description Get a Users Github Widget by their auth token
// @Description and a path parameter to return the GithubWidgetData
// @Description This is a special widget that needs authorization outside of
// @Description the mindscape so we use the github api to get the data.
//
// @ID          GetGithubWidgetData
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     GithubResponse
// @Failure     401                 {object}     GithubResponse
// @Failure     401                 {object}     GithubResponse
// @Failure     403                 {object}     GithubResponse
// @Failure     500                 {object}     GithubResponse
// @Router      /widgets/github/{user_widget_id}    [get]
func (wc WidgetController) GithubWidget(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	// use the Widget from ReadUserWidgetHandler
	return handlers.GithubWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get a Users Github Profile Widget Data
// @Description Get a Users Github Profile Widget by their auth token
// @Description and a path parameter to return only the profile data.
// @Description This is a fast endpoint that returns profile info quickly.
//
// @ID          GetGithubProfileWidgetData
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     GithubProfileResponse
// @Failure     401                 {object}     GithubProfileResponse
// @Failure     403                 {object}     GithubProfileResponse
// @Failure     500                 {object}     GithubProfileResponse
// @Router      /widgets/{user_widget_id}/github/profile    [get]
func (wc WidgetController) GithubProfileWidget(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.GithubProfileWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get a Users Github Commits Widget Data
// @Description Get a Users Github Commits Widget by their auth token
// @Description and a path parameter to return only the commits data.
// @Description This endpoint fetches commit history and may take longer.
//
// @ID          GetGithubCommitsWidgetData
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     GithubCommitsResponse
// @Failure     401                 {object}     GithubCommitsResponse
// @Failure     403                 {object}     GithubCommitsResponse
// @Failure     500                 {object}     GithubCommitsResponse
// @Router      /widgets/{user_widget_id}/github/commits    [get]
func (wc WidgetController) GithubCommitsWidget(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.GithubCommitWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get GitHub PRs Widget Data
// @Description Fetches open pull requests from GitHub (review-requested, authored, mentioned, assigned)
// @Description using the GitHub Issues Search API with a server-side PAT proxy.
// @Description The PAT requires 'repo' read scope for private repositories.
//
// @ID          GetGithubPRsWidgetData
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     responses.GithubPRsWidgetResponse
// @Failure     401                 {object}     responses.GithubPRsWidgetResponse
// @Failure     403                 {object}     responses.GithubPRsWidgetResponse
// @Failure     502                 {object}     responses.GithubPRsWidgetResponse
// @Router      /widgets/{user_widget_id}/github/prs    [get]
func (wc WidgetController) GithubPRsWidget(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.GithubPRsWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get Plex Recently Added
// @Description Fetches recently-added media items from the configured Plex server.
// @Description The Plex API token is kept server-side; this is a server-side proxy endpoint.
// @Description Poster art is returned as signed thumb URLs (serverUrl + thumb + X-Plex-Token query param).
//
// @ID          GetPlexRecentlyAdded
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     responses.PlexRecentlyAddedResponse
// @Failure     401                 {object}     responses.PlexRecentlyAddedResponse
// @Failure     403                 {object}     responses.PlexRecentlyAddedResponse
// @Failure     502                 {object}     responses.PlexRecentlyAddedResponse
// @Router      /widgets/{user_widget_id}/plex/recently-added    [get]
func (wc WidgetController) PlexRecentlyAdded(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.PlexWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Add a Users Widget
// @Description Add a Users Widget by their auth token.
// @Description The config is defined by the configuration parameters as defined by the user and the schema.
//
// @ID          AddUserWidget
// @Tags        Widgets
// @Produce     json
// @Accept      json
// @Security BearerAuth
// @Param       AddUserWidgetRequest body        requests.AddUserWidgetRequst      true "Add Widget Request"
// @Success     200                 {object}     responses.UserWidgetResponse
// @Failure     400                 {object}     responses.UserWidgetResponse
// @Failure     401                 {object}     responses.UserWidgetResponse
// @Failure     403                 {object}     responses.UserWidgetResponse
// @Failure     500                 {object}     responses.UserWidgetResponse
// @Router      /widgets            [post]
func (wc WidgetController) AddWidget(ctx echo.Context) error {
	return handlers.NewAddWidgetHandler(
		ctx,
		*wc.ValidatorService,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

//region COOLIFY
// @Summary Get a Users Coolify Applications
// @Description Get a Users Coolify Applications by their auth token
//
// @ID          GetUserCoolifyApplications
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     CoolifyWidgetApplicationResponse
// @Failure     401                 {object}     CoolifyWidgetApplicationResponse
// @Failure     403                 {object}     CoolifyWidgetApplicationResponse
// @Failure     500                 {object}     CoolifyWidgetApplicationResponse
// @Router      /widgets/{user_widget_id}/coolify/applications    [get]
func (wc WidgetController) CoolifyWidgetApplications(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}
	return handlers.CoolifyWidgetApplicationsJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Start a Coolify Application
// @Description Start a Coolify Application by app UUID
//
// @ID          StartCoolifyApplication
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       app_uuid            path         string                         true "Application UUID"
// @Success     200                 {object}     CoolifyActionResponse
// @Failure     400                 {object}     CoolifyActionResponse
// @Failure     401                 {object}     CoolifyActionResponse
// @Failure     403                 {object}     CoolifyActionResponse
// @Failure     500                 {object}     CoolifyActionResponse
// @Router      /widgets/{user_widget_id}/coolify/applications/{app_uuid}/start [post]
func (wc WidgetController) StartCoolifyApplication(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	appUUID := ctx.Param("app_uuid")
	return handlers.CoolifyApplicationStartJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
		appUUID,
	)
}

// @Summary Stop a Coolify Application
// @Description Stop a Coolify Application by app UUID
//
// @ID          StopCoolifyApplication
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       app_uuid            path         string                         true "Application UUID"
// @Success     200                 {object}     CoolifyActionResponse
// @Failure     400                 {object}     CoolifyActionResponse
// @Failure     401                 {object}     CoolifyActionResponse
// @Failure     403                 {object}     CoolifyActionResponse
// @Failure     500                 {object}     CoolifyActionResponse
// @Router      /widgets/{user_widget_id}/coolify/applications/{app_uuid}/stop [post]
func (wc WidgetController) StopCoolifyApplication(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	appUUID := ctx.Param("app_uuid")
	return handlers.CoolifyApplicationStopJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
		appUUID,
	)
}

// @Summary Restart a Coolify Application
// @Description Restart a Coolify Application by app UUID
//
// @ID          RestartCoolifyApplication
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Param       app_uuid            path         string                         true "Application UUID"
// @Success     200                 {object}     CoolifyActionResponse
// @Failure     400                 {object}     CoolifyActionResponse
// @Failure     401                 {object}     CoolifyActionResponse
// @Failure     403                 {object}     CoolifyActionResponse
// @Failure     500                 {object}     CoolifyActionResponse
// @Router      /widgets/{user_widget_id}/coolify/applications/{app_uuid}/restart [post]
func (wc WidgetController) RestartCoolifyApplication(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	appUUID := ctx.Param("app_uuid")
	return handlers.CoolifyApplicationRestartJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
		appUUID,
	)
}

// @Summary Get a Users Coolify Services 
// @Description Get a Users Coolify Services by their auth token
//
// @ID          GetUserCoolifyServices 
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     CoolifyWidgetServiceResponse
// @Failure     401                 {object}     CoolifyWidgetServiceResponse
// @Failure     403                 {object}     CoolifyWidgetServiceResponse
// @Failure     500                 {object}     CoolifyWidgetServiceResponse
// @Router      /widgets/{user_widget_id}/coolify/services [get]
func (wc WidgetController) CoolifyWidgetServices(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()
	
	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.CoolifyWidgetServicesJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get a Users Coolify Status
// @Description Get the combined Coolify application and service status for a widget
//
// @ID          GetUserCoolifyStatus
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     CoolifyWidgetResponse
// @Failure     400                 {object}     CoolifyWidgetResponse
// @Failure     401                 {object}     CoolifyWidgetResponse
// @Failure     403                 {object}     CoolifyWidgetResponse
// @Failure     500                 {object}     CoolifyWidgetResponse
// @Router      /widgets/{user_widget_id}/coolify [get]
func (wc WidgetController) CoolifyWidgetStatus(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.CoolifyWidgetJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}

// @Summary Get Coolify Server Metrics
// @Description Get CPU, memory and storage usage for the widget's Coolify server. CPU and memory come from the Sentinel agent; storage is read from the Mindscape host filesystem. Sources that fail are reported in data.warnings rather than failing the request.
//
// @ID          GetUserCoolifyMetrics
// @Tags        Widgets
// @Produce     json
// @Security BearerAuth
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     CoolifyWidgetMetricsResponse
// @Failure     400                 {object}     CoolifyWidgetMetricsResponse
// @Failure     401                 {object}     CoolifyWidgetMetricsResponse
// @Failure     403                 {object}     CoolifyWidgetMetricsResponse
// @Failure     500                 {object}     CoolifyWidgetMetricsResponse
// @Router      /widgets/{user_widget_id}/coolify/metrics [get]
func (wc WidgetController) CoolifyWidgetMetrics(ctx echo.Context) error {
	widget := handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle()

	if widget.Error() != nil {
		return widget.JSON()
	}

	return handlers.CoolifyWidgetMetricsJsonHandler(
		ctx,
		widget.Data().(*repository.UserWidget),
	)
}
//endregion


func (wc WidgetController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api" + wc.Name)
	api.GET("/schemas", wc.ReadSchemas)
	api.GET("/schemas/:schema_id", wc.GetSchemaByID)

	api.GET("", wc.Read, middlewares...)
	api.GET("/:user_widget_id", wc.ReadById, middlewares...)
	api.GET("/github/:user_widget_id", wc.GithubWidget, middlewares...)
	api.GET("/:user_widget_id/github/profile", wc.GithubProfileWidget, middlewares...)
	api.GET("/:user_widget_id/github/commits", wc.GithubCommitsWidget, middlewares...)
	api.GET("/:user_widget_id/github/prs", wc.GithubPRsWidget, middlewares...)
	api.GET("/:user_widget_id/plex/recently-added", wc.PlexRecentlyAdded, middlewares...)
	api.GET("/:user_widget_id/coolify", wc.CoolifyWidgetStatus, middlewares...)
	api.GET("/:user_widget_id/coolify/metrics", wc.CoolifyWidgetMetrics, middlewares...)
	api.GET("/:user_widget_id/coolify/applications", wc.CoolifyWidgetApplications, middlewares...)
	api.GET("/:user_widget_id/coolify/services", wc.CoolifyWidgetServices, middlewares...)
	api.POST("/:user_widget_id/coolify/applications/:app_uuid/start", wc.StartCoolifyApplication, middlewares...)
	api.POST("/:user_widget_id/coolify/applications/:app_uuid/stop", wc.StopCoolifyApplication, middlewares...)
	api.POST("/:user_widget_id/coolify/applications/:app_uuid/restart", wc.RestartCoolifyApplication, middlewares...)
	api.POST("", wc.AddWidget, middlewares...)
}
