package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// PlexClient fetches data from a Plex Media Server via server-side proxy.
//
// Design decision — server-side proxy:
//   - The X-Plex-Token is sent only in server-to-Plex requests; it never appears in
//     browser JS or localStorage.
//   - Poster art is returned as signed thumb URLs ({serverUrl}{thumb}?X-Plex-Token={token})
//     so the browser can load images directly. The token therefore appears in image
//     request URLs visible in browser network tools, but it is not stored client-side.
//     A full image-proxy endpoint would remove this exposure; that is left for a future
//     iteration if required by the deployment threat model.
type PlexClient struct {
	httpClient     *http.Client
	apiToken       string
	serverUrl      string
	defaultTimeout time.Duration
}

type PlexClientOption func(*PlexClient)

func WithPlexTimeout(timeout time.Duration) PlexClientOption {
	return func(c *PlexClient) {
		c.defaultTimeout = timeout
	}
}

func NewPlexClient(serverUrl, apiToken string, opts ...PlexClientOption) *PlexClient {
	client := &PlexClient{
		httpClient:     &http.Client{},
		apiToken:       apiToken,
		serverUrl:      serverUrl,
		defaultTimeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// PlexMetadata represents a single media item from the Plex /library/recentlyAdded endpoint.
type PlexMetadata struct {
	Title     string `json:"title"`
	Thumb     string `json:"thumb"`
	AddedAt   int64  `json:"addedAt"`
	Type      string `json:"type"`
	Year      int    `json:"year"`
	ParentTitle string `json:"parentTitle"` // show title for episodes
}

type plexMediaContainer struct {
	Metadata []PlexMetadata `json:"Metadata"`
}

type plexResponse struct {
	MediaContainer plexMediaContainer `json:"MediaContainer"`
}

// FetchRecentlyAdded retrieves recently-added items from the Plex server.
// maxItems controls the page size; pass 0 to use the Plex server default.
func (c *PlexClient) FetchRecentlyAdded(ctx echo.Context, maxItems int) ([]PlexMetadata, error) {
	url := fmt.Sprintf("%s/library/recentlyAdded", c.serverUrl)
	if maxItems > 0 {
		url = fmt.Sprintf("%s?X-Plex-Container-Size=%d", url, maxItems)
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request().Context(), c.defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create plex request: %w", err)
	}

	req.Header.Set("X-Plex-Token", c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("plex request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plex API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read plex response: %w", err)
	}

	var plexResp plexResponse
	if err := json.Unmarshal(body, &plexResp); err != nil {
		return nil, fmt.Errorf("failed to parse plex JSON: %w", err)
	}

	items := plexResp.MediaContainer.Metadata

	// Attach signed thumb URLs so the browser can load poster art directly.
	// The X-Plex-Token is embedded in the URL query string — see PlexClient
	// design-decision comment above regarding the trade-off.
	for i := range items {
		if items[i].Thumb != "" {
			items[i].Thumb = fmt.Sprintf("%s%s?X-Plex-Token=%s", c.serverUrl, items[i].Thumb, c.apiToken)
		}
	}

	return items, nil
}
