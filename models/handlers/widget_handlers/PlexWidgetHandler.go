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

type PlexWidgetConfig struct {
	ServerUrl       string   `json:"serverUrl"`
	ApiToken        string   `json:"apiToken"`
	Libraries       []string `json:"libraries"`
	MaxItems        int      `json:"maxItems"`
	RefreshInterval int      `json:"refreshInterval"`
}

type PlexWidgetHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   []clients.PlexMetadata
}

func NewPlexWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *PlexWidgetHandler {
	return &PlexWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   nil,
	}
}

func PlexWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewPlexWidgetHandler(ctx, widget)
	return handler.Handle().JSON()
}

func (h *PlexWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *PlexWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *PlexWidgetHandler) Code() int                            { return h.code }
func (h *PlexWidgetHandler) Error() error                         { return h.err }
func (h *PlexWidgetHandler) Data() any                            { return h.data }

func (h *PlexWidgetHandler) Handle() handlers.IHandler {
	var config PlexWidgetConfig
	if err := json.Unmarshal(h.widget.Config, &config); err != nil {
		return handlers.Lock(h, 400, err)
	}

	maxItems := config.MaxItems
	if maxItems <= 0 {
		maxItems = 10
	}

	plexClient := clients.NewPlexClient(
		config.ServerUrl,
		config.ApiToken,
		clients.WithPlexTimeout(30*time.Second),
	)

	items, err := plexClient.FetchRecentlyAdded(h.ctx, maxItems)
	if err != nil {
		return handlers.Lock(h, 502, err)
	}

	h.data = items
	return h
}

func (h *PlexWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewPlexRecentlyAddedResponse().Successful(h.ctx, h.data)
	}
	return responses.NewPlexRecentlyAddedResponse().Fail(h.ctx, h.code, h.err)
}
