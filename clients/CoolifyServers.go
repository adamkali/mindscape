package clients

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

// CoolifyServer is the subset of Coolify's Server model the metrics widget
// needs. Coolify returns a much wider object; everything else is deploy-time
// plumbing that has no place on a dashboard.
type CoolifyServer struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	ProxyType   string `json:"proxy_type"`
}

// CoolifySentinelSettings mirrors GET /servers/{uuid}/sentinel. SentinelToken
// and SentinelCustomURL are only populated when the personal access token
// carries the read:sensitive (or root) ability — without it Coolify silently
// omits both fields rather than failing the request.
//
// Two traps here, both confirmed against a live v4.3.14 instance:
//
//   - SentinelCustomURL is the Coolify address Sentinel *pushes metrics to*,
//     not the address Sentinel listens on. It is useless for reaching Sentinel
//     and equals the agent's own PUSH_ENDPOINT.
//   - SentinelToken is whatever Coolify has on record, which drifts from the
//     agent's real TOKEN if the container was not recreated after a rotation.
//     A stale value fails every Sentinel call with 401.
//
// Treat both as hints. Explicit widget config is the dependable path.
type CoolifySentinelSettings struct {
	IsSentinelEnabled bool   `json:"is_sentinel_enabled"`
	IsMetricsEnabled  bool   `json:"is_metrics_enabled"`
	SentinelToken     string `json:"sentinel_token"`
	// SentinelCustomURL is Sentinel's push target, not its listen address.
	SentinelCustomURL string `json:"sentinel_custom_url"`
	RefreshRateSecs   int    `json:"sentinel_metrics_refresh_rate_seconds"`
	HistoryDays       int    `json:"sentinel_metrics_history_days"`
	PushIntervalSecs  int    `json:"sentinel_push_interval_seconds"`
}

// GetCoolifyServers lists the servers the token's team can see.
func (c *CoolifyClient) GetCoolifyServers(ctx echo.Context) ([]CoolifyServer, error) {
	url := fmt.Sprintf("%s/servers", c.baseUrl)
	body, err := c.doRequest(ctx.Request().Context(), url, "application/json")
	if err != nil {
		return nil, err
	}

	var servers []CoolifyServer
	if err := json.Unmarshal(body, &servers); err != nil {
		return nil, fmt.Errorf("failed to parse servers JSON: %w", err)
	}
	return servers, nil
}

// GetCoolifyServerSentinel reads a server's Sentinel settings, which is how the
// widget discovers where to reach Sentinel when the config does not say.
func (c *CoolifyClient) GetCoolifyServerSentinel(
	ctx echo.Context,
	serverUUID string,
) (*CoolifySentinelSettings, error) {
	url := fmt.Sprintf("%s/servers/%s/sentinel", c.baseUrl, serverUUID)
	body, err := c.doRequest(ctx.Request().Context(), url, "application/json")
	if err != nil {
		return nil, err
	}

	var settings CoolifySentinelSettings
	if err := json.Unmarshal(body, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse sentinel settings JSON: %w", err)
	}
	return &settings, nil
}
