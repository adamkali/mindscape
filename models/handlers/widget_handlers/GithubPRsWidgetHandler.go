package widget_handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type GithubPRsWidgetConfig struct {
	GithubToken           string   `json:"githubToken"`
	Repos                 []string `json:"repos"`
	FilterReviewRequested bool     `json:"filterReviewRequested"`
	FilterAuthored        bool     `json:"filterAuthored"`
	FilterMentioned       bool     `json:"filterMentioned"`
	FilterAssigned        bool     `json:"filterAssigned"`
	RefreshInterval       int      `json:"refreshInterval"`
}

type GithubPRsWidgetHandler struct {
	ctx             echo.Context
	code            int
	err             error
	widget          *repository.UserWidget
	reviewRequested []clients.PullRequest
	authored        []clients.PullRequest
	mentioned       []clients.PullRequest
	assigned        []clients.PullRequest
}

func NewGithubPRsWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *GithubPRsWidgetHandler {
	return &GithubPRsWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
	}
}

func GithubPRsWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewGithubPRsWidgetHandler(ctx, widget)
	return handler.Handle().JSON()
}

func (h *GithubPRsWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GithubPRsWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GithubPRsWidgetHandler) Code() int                            { return h.code }
func (h *GithubPRsWidgetHandler) Error() error                         { return h.err }
func (h *GithubPRsWidgetHandler) Data() any                            { return nil }

func (h *GithubPRsWidgetHandler) Handle() handlers.IHandler {
	var config GithubPRsWidgetConfig
	if err := json.Unmarshal(h.widget.Config, &config); err != nil {
		return handlers.Lock(h, 400, err)
	}

	// Every filter below is scoped with `@me`, which GitHub only resolves for an
	// authenticated request. Without a token the search comes back 401 and the
	// widget renders four empty lists, so say so plainly instead.
	token := strings.TrimSpace(config.GithubToken)
	if token == "" {
		return handlers.Lock(h, 400, fmt.Errorf(
			"no GitHub personal access token configured for this widget: add a classic PAT with the 'repo' scope",
		))
	}

	// If no filter is explicitly enabled, default to showing all three main sections.
	showAll := !config.FilterReviewRequested && !config.FilterAuthored &&
		!config.FilterMentioned && !config.FilterAssigned

	ghClient := clients.NewGitHubClient(token, clients.WithTimeout(30*time.Second))

	if showAll || config.FilterReviewRequested {
		items, err := ghClient.FetchPRs(h.ctx, "review-requested:@me", config.Repos)
		if err != nil {
			return handlers.Lock(h, 502, err)
		}
		h.reviewRequested = items
	}

	if showAll || config.FilterAuthored {
		items, err := ghClient.FetchPRs(h.ctx, "author:@me", config.Repos)
		if err != nil {
			return handlers.Lock(h, 502, err)
		}
		h.authored = items
	}

	if showAll || config.FilterMentioned {
		items, err := ghClient.FetchPRs(h.ctx, "mentions:@me", config.Repos)
		if err != nil {
			return handlers.Lock(h, 502, err)
		}
		h.mentioned = items
	}

	if config.FilterAssigned {
		items, err := ghClient.FetchPRs(h.ctx, "assignee:@me", config.Repos)
		if err != nil {
			return handlers.Lock(h, 502, err)
		}
		h.assigned = items
	}

	return h
}

func (h *GithubPRsWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewGithubPRsWidgetResponse().Successful(
			h.ctx,
			h.reviewRequested,
			h.authored,
			h.mentioned,
			h.assigned,
		)
	}
	return responses.NewGithubPRsWidgetResponse().Fail(h.ctx, h.code, h.err)
}
