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

type GithubCommitWidgetHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   *clients.FinalDateEntry
}

func NewGithubCommitWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *GithubCommitWidgetHandler {
	return &GithubCommitWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   nil,
	}
}

func GithubCommitWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := &GithubCommitWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
	}
	return handler.Handle().JSON()
}

func (h *GithubCommitWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GithubCommitWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GithubCommitWidgetHandler) Code() int                            { return h.code }
func (h *GithubCommitWidgetHandler) Error() error                         { return h.err }
func (h *GithubCommitWidgetHandler) Data() any                            { return h.data }

func (h *GithubCommitWidgetHandler) Handle() handlers.IHandler {
	githubWidgetConfigBytes := h.widget.Config

	var githubWidgetConfig GithubWidgetConfig
	err := json.Unmarshal(githubWidgetConfigBytes, &githubWidgetConfig)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	githubClient := clients.NewGitHubClient(githubWidgetConfig.PersonalAccessToken, clients.WithTimeout(90*time.Second))

	githubCommits, err := githubClient.FetchCommits(h.ctx, githubWidgetConfig.Username, 5)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	h.data = githubCommits

	return h
}

func (h *GithubCommitWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewGithubCommitWidgetResponse().Successful(h.ctx, h.data)
	} else {
		return responses.NewGithubCommitWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
