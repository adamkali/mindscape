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
type Data struct {
		Profile *clients.GitHubProfile  `json:"profile"`
		Commits *clients.FinalDateEntry `json:"commits"`
	}

type GithubWidgetHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   Data
}

func NewGithubWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *GithubWidgetHandler {
	return &GithubWidgetHandler{
		ctx:           ctx,
		code:          200,
		err:           nil,
		widget:        widget,
		data:          Data{},
	}
}

func GithubWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := &GithubWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
	}
	return handler.Handle().JSON()
}

type GithubWidgetConfig struct {
	Username            string `json:"username"`
	PersonalAccessToken string `json:"personalAccessToken"`
	ShowStats           bool   `json:"showStats"`
}

func (h *GithubWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GithubWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GithubWidgetHandler) Code() int                            { return h.code }
func (h *GithubWidgetHandler) Error() error                         { return h.err }
func (h *GithubWidgetHandler) Data() any                            { return h.data }

// Methods:
func (h *GithubWidgetHandler) Handle() handlers.IHandler {
	githubWidgetConfigBytes := h.widget.Config

	var githubWidgetConfig GithubWidgetConfig
	err := json.Unmarshal(githubWidgetConfigBytes, &githubWidgetConfig)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	githubClient := clients.NewGitHubClient(githubWidgetConfig.PersonalAccessToken, clients.WithTimeout(90*time.Second))

	githubProfile, err := githubClient.FetchProfile(h.ctx, githubWidgetConfig.Username)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	githubCommits, err := githubClient.FetchCommits(h.ctx, githubWidgetConfig.Username, 5)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	h.data = Data{
		Profile: githubProfile,
		Commits: githubCommits,
	} 

	return h
}

func (h *GithubWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewGithubWidgetResponse().Successful(h.ctx, h.data.Profile, h.data.Commits)
	} else {
		return responses.NewGithubWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
