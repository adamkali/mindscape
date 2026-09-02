package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// SentinelClient talks to Coolify's Sentinel agent, which is the only thing on
// a Coolify host that reports live resource usage — the Coolify REST API
// exposes container `limits_*` but no actual CPU or memory readings.
//
// Sentinel (v0.0.22) serves exactly four metric routes:
//
//	GET /api/health           -> "ok"
//	GET /api/cpu/current      -> {"percent": 14.6, "time": "..."}
//	GET /api/cpu/history      -> [{"time": "...", "percent": "14.60"}, ...]
//	GET /api/memory/current   -> {"total": ..., "used": ..., "usedPercent": ...}
//	GET /api/memory/history   -> [ ...same shape... ]
//
// There is no disk route, which is why storage comes from DiskUsage instead.
// Note that history reports `percent` as a string while current reports it as
// a number; only the current endpoints are used here.
type SentinelClient struct {
	httpClient     *http.Client
	defaultTimeout time.Duration
	token          string
	baseUrl        string
}

type SentinelClientOption func(*SentinelClient)

// SentinelCPU is GET /api/cpu/current.
type SentinelCPU struct {
	Percent float64 `json:"percent"`
	Time    string  `json:"time"`
}

// SentinelMemory is GET /api/memory/current. Byte counts are absolute.
type SentinelMemory struct {
	Time        string  `json:"time"`
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

func NewSentinelClient(token string, baseUrl string, opts ...SentinelClientOption) *SentinelClient {
	client := &SentinelClient{
		httpClient:     &http.Client{},
		defaultTimeout: 10 * time.Second,
		token:          token,
		// Sentinel is normally addressed without a trailing slash; tolerate one
		// so a pasted config value does not produce `//api/cpu/current`.
		baseUrl: strings.TrimRight(strings.TrimSpace(baseUrl), "/"),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func SentinelWithHTTPClient(httpClient *http.Client) SentinelClientOption {
	return func(c *SentinelClient) {
		c.httpClient = httpClient
	}
}

func SentinelWithTimeout(timeout time.Duration) SentinelClientOption {
	return func(c *SentinelClient) {
		c.defaultTimeout = timeout
	}
}

// BaseUrl reports the resolved Sentinel address, so a handler can tell the
// caller which host the numbers came from.
func (c *SentinelClient) BaseUrl() string { return c.baseUrl }

func (c *SentinelClient) doRequest(ctx context.Context, path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf(
				"Sentinel at %s did not respond within %s: it listens on the Coolify host and is not published publicly by default",
				c.baseUrl, c.defaultTimeout,
			)
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("request canceled: %w", err)
		}
		return nil, fmt.Errorf("could not reach Sentinel at %s: %w", c.baseUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		// Sentinel authenticates in middleware ahead of routing, so a bad token
		// turns every path into a 401 — including ones that do not exist.
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf(
				"Sentinel rejected the token (401): set the widget's sentinelToken to the TOKEN value from the coolify-sentinel container",
			)
		}
		return nil, fmt.Errorf(
			"Sentinel API error: %d (%s)",
			resp.StatusCode, strings.TrimSpace(string(detail)),
		)
	}

	return io.ReadAll(resp.Body)
}

// HealthCheck hits /api/health, the one route Sentinel serves unauthenticated.
func (c *SentinelClient) HealthCheck(ctx echo.Context) error {
	body, err := c.doRequest(ctx.Request().Context(), "/api/health")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(body)) != "ok" {
		return fmt.Errorf("sentinel health check failed: %s", strings.TrimSpace(string(body)))
	}
	return nil
}

// CurrentCPU reads instantaneous CPU utilisation as a percentage.
func (c *SentinelClient) CurrentCPU(ctx echo.Context) (*SentinelCPU, error) {
	body, err := c.doRequest(ctx.Request().Context(), "/api/cpu/current")
	if err != nil {
		return nil, err
	}

	var cpu SentinelCPU
	if err := json.Unmarshal(body, &cpu); err != nil {
		return nil, fmt.Errorf("failed to parse sentinel CPU JSON: %w", err)
	}
	return &cpu, nil
}

// CurrentMemory reads instantaneous memory usage in bytes plus a percentage.
func (c *SentinelClient) CurrentMemory(ctx echo.Context) (*SentinelMemory, error) {
	body, err := c.doRequest(ctx.Request().Context(), "/api/memory/current")
	if err != nil {
		return nil, err
	}

	var mem SentinelMemory
	if err := json.Unmarshal(body, &mem); err != nil {
		return nil, fmt.Errorf("failed to parse sentinel memory JSON: %w", err)
	}
	return &mem, nil
}
