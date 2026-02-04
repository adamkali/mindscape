package widget_handlers

import (
	"encoding/json"
	"time"

	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type GithubProfileWidgetHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   *clients.GitHubProfile
}

func NewGithubProfileWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *GithubProfileWidgetHandler {
	return &GithubProfileWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   nil,
	}
}

func GithubProfileWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := &GithubProfileWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
	}
	return handler.Handle().JSON()
}

func (h *GithubProfileWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GithubProfileWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GithubProfileWidgetHandler) Code() int                            { return h.code }
func (h *GithubProfileWidgetHandler) Error() error                         { return h.err }
func (h *GithubProfileWidgetHandler) Data() any                            { return h.data }

func (h *GithubProfileWidgetHandler) Handle() handlers.IHandler {
	githubWidgetConfigBytes := h.widget.Config

	var githubWidgetConfig GithubWidgetConfig
	err := json.Unmarshal(githubWidgetConfigBytes, &githubWidgetConfig)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	githubClient := clients.NewGitHubClient(githubWidgetConfig.PersonalAccessToken, clients.WithTimeout(30*time.Second))

	githubProfile, err := githubClient.FetchProfile(h.ctx, githubWidgetConfig.Username)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	h.data = githubProfile

	return h
}

func (h *GithubProfileWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewGithubProfileWidgetResponse().Successful(h.ctx, h.data)
	} else {
		return responses.NewGithubProfileWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
