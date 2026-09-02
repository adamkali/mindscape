package widget_handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetHandler struct {
	ctx         echo.Context
	code        int
	err         error
	widget      *repository.UserWidget
	applicatons []clients.CoolifyApplication
	services    []clients.CoolifyService
}

func NewCoolifyWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetHandler {
	return &CoolifyWidgetHandler{
		ctx:         ctx,
		code:        200,
		err:         nil,
		widget:      widget,
		applicatons: []clients.CoolifyApplication{},
		services:    []clients.CoolifyService{},
	}
}

func CoolifyWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewCoolifyWidgetHandler(ctx, widget)
	return handler.Handle().JSON()
}

type CoolifyWidgetConfig struct {
	BaseUrl             string   `json:"baseUrl"`
	PersonalAccessToken string   `json:"personalAccessToken"`
	ApplicationIds      []string `json:"applicationIds"`
	ServiceIds          []string `json:"serviceIds"`
	// Metrics settings. SentinelUrl/SentinelToken address Coolify's Sentinel
	// agent directly; leaving them blank makes the metrics handler try to
	// discover them through the Coolify API instead. ServerUuid disambiguates
	// when a token can see more than one server, and DiskPath selects the
	// filesystem reported as storage (default "/").
	SentinelUrl   string `json:"sentinelUrl"`
	SentinelToken string `json:"sentinelToken"`
	ServerUuid    string `json:"serverUuid"`
	DiskPath      string `json:"diskPath"`
}

func (h *CoolifyWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *CoolifyWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *CoolifyWidgetHandler) Code() int                            { return h.code }
func (h *CoolifyWidgetHandler) Error() error                         { return h.err }

func (h *CoolifyWidgetHandler) Data() any {
	return struct {
		Applicatons []clients.CoolifyApplication
		Services    []clients.CoolifyService
	}{
		Applicatons: h.applicatons,
		Services:    h.services,
	}
}

func (h *CoolifyWidgetHandler) Handle() handlers.IHandler {
	if h.widget.Config == nil {
		return handlers.Lock(h, 400, fmt.Errorf("coolify widget has no configuration"))
	}

	var config CoolifyWidgetConfig
	if err := json.Unmarshal(h.widget.Config, &config); err != nil {
		return handlers.Lock(h, 400, err)
	}

	if strings.TrimSpace(config.PersonalAccessToken) == "" {
		return handlers.Lock(h, 400, fmt.Errorf(
			"no Coolify personal access token configured for this widget",
		))
	}
	if strings.TrimSpace(config.BaseUrl) == "" {
		return handlers.Lock(h, 400, fmt.Errorf(
			"no Coolify base URL configured for this widget (for example https://coolify.example.com/api/v1)",
		))
	}

	client := clients.NewCoolifyClient(
		strings.TrimSpace(config.PersonalAccessToken),
		strings.TrimRight(strings.TrimSpace(config.BaseUrl), "/"),
	)

	applications, err := client.GetCoolifyApplications(h.ctx)
	if err != nil {
		return handlers.Lock(h, 502, err)
	}
	services, err := client.GetCoolifyServices(h.ctx)
	if err != nil {
		return handlers.Lock(h, 502, err)
	}

	// An empty id list means "show everything", matching the schema copy.
	if len(config.ApplicationIds) > 0 {
		applications = keepApplications(applications, config.ApplicationIds)
	}
	if len(config.ServiceIds) > 0 {
		services = keepServices(services, config.ServiceIds)
	}

	h.applicatons = applications
	h.services = services

	return h
}

// keepApplications and keepServices replace clients.FilterApplications and
// clients.FilterServices, which assign to their local slice parameter and so
// discard the result.
func keepApplications(
	apps []clients.CoolifyApplication,
	ids []string,
) []clients.CoolifyApplication {
	wanted := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		wanted[strings.TrimSpace(id)] = struct{}{}
	}

	kept := make([]clients.CoolifyApplication, 0, len(apps))
	for _, app := range apps {
		if _, ok := wanted[app.UUID]; ok {
			kept = append(kept, app)
		}
	}
	return kept
}

func keepServices(
	services []clients.CoolifyService,
	ids []string,
) []clients.CoolifyService {
	wanted := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		wanted[strings.TrimSpace(id)] = struct{}{}
	}

	kept := make([]clients.CoolifyService, 0, len(services))
	for _, service := range services {
		if _, ok := wanted[service.UUID]; ok {
			kept = append(kept, service)
		}
	}
	return kept
}

func (h *CoolifyWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetResponse().Successful(h.ctx, h.applicatons, h.services)
	} else {
		return responses.NewCoolifyWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
