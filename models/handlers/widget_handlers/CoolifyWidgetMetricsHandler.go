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

// sentinelDefaultPort is the port Sentinel listens on; it is not configurable
// through Coolify, so discovery can assume it.
const sentinelDefaultPort = 8888

// CoolifyWidgetMetricsHandler reports CPU, memory and storage for a Coolify
// server.
//
// The three numbers come from two different places because Coolify's REST API
// reports neither: CPU and memory come from the Sentinel agent on the host,
// and storage is read from the local filesystem since Sentinel has no disk
// route. Any source that fails contributes a warning instead of failing the
// whole request — a dashboard tile row should degrade, not go blank.
type CoolifyWidgetMetricsHandler struct {
	ctx         echo.Context
	code        int
	err         error
	widget      *repository.UserWidget
	serverName  string
	serverUUID  string
	sentinelURL string
	cpu         *clients.SentinelCPU
	memory      *clients.SentinelMemory
	storage     *clients.Disk
	warnings    []string
}

func NewCoolifyWidgetMetricsHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetMetricsHandler {
	return &CoolifyWidgetMetricsHandler{
		ctx:      ctx,
		code:     200,
		err:      nil,
		widget:   widget,
		warnings: []string{},
	}
}

func CoolifyWidgetMetricsJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewCoolifyWidgetMetricsHandler(ctx, widget)
	return handler.Handle().JSON()
}

func (h *CoolifyWidgetMetricsHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *CoolifyWidgetMetricsHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *CoolifyWidgetMetricsHandler) Code() int    { return h.code }
func (h *CoolifyWidgetMetricsHandler) Error() error { return h.err }
func (h *CoolifyWidgetMetricsHandler) Data() any    { return nil }

func (h *CoolifyWidgetMetricsHandler) warn(format string, args ...any) {
	h.warnings = append(h.warnings, fmt.Sprintf(format, args...))
}

func (h *CoolifyWidgetMetricsHandler) Handle() handlers.IHandler {
	var config CoolifyWidgetConfig
	if err := json.Unmarshal(h.widget.Config, &config); err != nil {
		return handlers.Lock(h, 400, err)
	}

	h.collectStorage(config)

	sentinelURL, sentinelToken := h.resolveSentinel(config)
	if sentinelURL == "" {
		// resolveSentinel has already recorded why.
		return h
	}
	h.sentinelURL = sentinelURL

	sentinel := clients.NewSentinelClient(
		sentinelToken,
		sentinelURL,
		clients.SentinelWithTimeout(10*time.Second),
	)

	if cpu, err := sentinel.CurrentCPU(h.ctx); err != nil {
		h.warn("CPU unavailable: %v", err)
	} else {
		h.cpu = cpu
	}

	if memory, err := sentinel.CurrentMemory(h.ctx); err != nil {
		h.warn("Memory unavailable: %v", err)
	} else {
		h.memory = memory
	}

	return h
}

// collectStorage reads the filesystem Mindscape itself runs on. The path is
// configurable so a deployment can point at the volume that actually matters.
func (h *CoolifyWidgetMetricsHandler) collectStorage(config CoolifyWidgetConfig) {
	path := strings.TrimSpace(config.DiskPath)
	if path == "" {
		path = "/"
	}

	disk, err := clients.DiskUsage(path)
	if err != nil {
		h.warn("Storage unavailable: %v", err)
		return
	}
	h.storage = disk
}

// resolveSentinel finds the Sentinel address and token, preferring explicit
// widget config and falling back to asking Coolify. It returns ("", "") when
// neither route works, having appended a warning saying which.
//
// Discovery is best-effort: Coolify's recorded sentinel_token drifts from the
// agent's real one, and Sentinel is usually bound to the server rather than
// published, so an explicitly configured sentinelUrl/sentinelToken pair is the
// only combination guaranteed to work.
func (h *CoolifyWidgetMetricsHandler) resolveSentinel(config CoolifyWidgetConfig) (string, string) {
	if url := strings.TrimSpace(config.SentinelUrl); url != "" {
		// Explicit config skips the Coolify round-trip entirely, so the server
		// name stays unknown; report the uuid we were given rather than nothing.
		h.serverUUID = strings.TrimSpace(config.ServerUuid)
		return url, strings.TrimSpace(config.SentinelToken)
	}

	baseUrl := strings.TrimSpace(config.BaseUrl)
	token := strings.TrimSpace(config.PersonalAccessToken)
	if baseUrl == "" || token == "" {
		h.warn("CPU and memory unavailable: set sentinelUrl and sentinelToken on the widget, or a Coolify baseUrl and personalAccessToken so they can be discovered")
		return "", ""
	}

	client := clients.NewCoolifyClient(token, baseUrl)

	server, err := h.pickServer(client, config)
	if err != nil {
		h.warn("CPU and memory unavailable: %v", err)
		return "", ""
	}
	h.serverName = server.Name
	h.serverUUID = server.UUID

	settings, err := client.GetCoolifyServerSentinel(h.ctx, server.UUID)
	if err != nil {
		h.warn("CPU and memory unavailable: could not read Sentinel settings for %q: %v", server.Name, err)
		return "", ""
	}
	if !settings.IsSentinelEnabled {
		h.warn("CPU and memory unavailable: Sentinel is disabled on server %q", server.Name)
		return "", ""
	}

	// Deliberately NOT settings.SentinelCustomURL: that is the Coolify address
	// Sentinel pushes metrics *to*, not the address it listens on. Sentinel
	// serves its API on port 8888 of the server itself, so derive it from the
	// server's IP instead.
	if settings.SentinelToken == "" {
		h.warn(
			"CPU and memory unavailable: Coolify returned no sentinel_token, which needs a personal access token with the read:sensitive ability - set sentinelUrl and sentinelToken on the widget instead",
		)
		return "", ""
	}

	ip := strings.TrimSpace(server.IP)
	if ip == "" || ip == "host.docker.internal" {
		// Coolify records its own host as host.docker.internal, an alias that
		// only resolves inside Coolify's containers.
		h.warn(
			"CPU and memory unavailable: server %q reports its address as %q, which does not resolve from here - set sentinelUrl on the widget (Sentinel listens on port 8888)",
			server.Name, server.IP,
		)
		return "", ""
	}

	return fmt.Sprintf("http://%s:%d", ip, sentinelDefaultPort), settings.SentinelToken
}

// pickServer honours an explicit serverUuid and otherwise takes the only
// server. It refuses to guess between several, since reporting the wrong host's
// CPU is worse than reporting none.
func (h *CoolifyWidgetMetricsHandler) pickServer(
	client *clients.CoolifyClient,
	config CoolifyWidgetConfig,
) (*clients.CoolifyServer, error) {
	servers, err := client.GetCoolifyServers(h.ctx)
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("Coolify returned no servers for this token")
	}

	wanted := strings.TrimSpace(config.ServerUuid)
	if wanted != "" {
		for i := range servers {
			if servers[i].UUID == wanted {
				return &servers[i], nil
			}
		}
		return nil, fmt.Errorf("no Coolify server with uuid %q", wanted)
	}

	if len(servers) > 1 {
		names := make([]string, 0, len(servers))
		for _, server := range servers {
			names = append(names, fmt.Sprintf("%s (%s)", server.Name, server.UUID))
		}
		return nil, fmt.Errorf(
			"this token can see %d servers, so set serverUuid on the widget to choose one: %s",
			len(servers), strings.Join(names, ", "),
		)
	}

	return &servers[0], nil
}

func (h *CoolifyWidgetMetricsHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetMetricsResponse().Successful(
			h.ctx,
			h.serverName,
			h.serverUUID,
			h.sentinelURL,
			h.cpu,
			h.memory,
			h.storage,
			h.warnings,
		)
	}
	return responses.NewCoolifyWidgetMetricsResponse().Fail(h.ctx, h.code, h.err)
}
