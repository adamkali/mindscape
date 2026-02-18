package controllers

import (
	handlers "github.com/adamkali/mindscape/models/handlers/apikey_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type ApiKeyController struct {
	Name      string
	Registrar *services.Registrar
}

func (c ApiKeyController) ControllerName() string {
	return c.Name
}

func BuildApiKeyController(registrar *services.Registrar) ApiKeyController {
	return ApiKeyController{
		Name:      "/apikeys",
		Registrar: registrar,
	}
}

// @Summary Create a new API Key
// @Description Create a new API Key for programmatic access
//
// @ID          CreateApiKey
// @Tags        ApiKeys
// @Accept      json
// @Produce     json
// @Param       Authorization         header       string                          true "auth header"     default(Bearer token)
// @Param       CreateApiKeyRequest   body         CreateApiKeyRequest             true "CreateApiKeyRequest"
// @Success     200                   {object}     responses.ApiKeyResponse
// @Failure     400                   {object}     responses.ApiKeyResponse
// @Failure     401                   {object}     responses.ApiKeyResponse
// @Failure     500                   {object}     responses.ApiKeyResponse
// @Router      /apikeys              [post]
func (c ApiKeyController) Create(e echo.Context) error {
	return handlers.CreateHandlerJsonHandler(e, *c.Registrar)
}

// @Summary List API Keys
// @Description List all API Keys for the current user
//
// @ID          ListApiKeys
// @Tags        ApiKeys
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Success     200                 {object}     responses.ApiKeysResponse
// @Failure     401                 {object}     responses.ApiKeysResponse
// @Failure     500                 {object}     responses.ApiKeysResponse
// @Router      /apikeys            [get]
func (c ApiKeyController) List(e echo.Context) error {
	return handlers.ListHandlerJsonHandler(e, *c.Registrar)
}

// @Summary Delete an API Key
// @Description Delete an API Key by its ID
//
// @ID          DeleteApiKey
// @Tags        ApiKeys
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       keyId               path         string                         true "API Key ID"
// @Success     200                 {object}     responses.StringResponse
// @Failure     400                 {object}     responses.StringResponse
// @Failure     401                 {object}     responses.StringResponse
// @Failure     500                 {object}     responses.StringResponse
// @Router      /apikeys/{keyId}    [delete]
func (c ApiKeyController) Delete(e echo.Context) error {
	return handlers.DeleteHandlerJsonHandler(e, *c.Registrar)
}

func (c ApiKeyController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api" + c.Name)
	api.POST("", c.Create, middlewares...)
	api.GET("", c.List, middlewares...)
	api.DELETE("/:keyId", c.Delete, middlewares...)
}
