package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type GitHubClient struct {
	httpClient          *http.Client
	personalAccessToken string
	defaultTimeout      time.Duration

}

var (
	baseColor   = "#10106c"
	targetColor = "#1e1efc"
)

type GitHubProfile struct {
	Login       string `json:"login"`
	ID          int    `json:"id"`
	AvatarURL   string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	PublicGists int    `json:"public_gists"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CommitAuthor struct {
	Date string `json:"date"`
}

type CommitData struct {
	Author CommitAuthor `json:"author"`
}

type CommitItem struct {
	SHA    string     `json:"sha"`
	Commit CommitData `json:"commit"`
}

type SearchResponse struct {
	Items []CommitItem `json:"items"`
}

type DateEntry struct {
	Date    string
	Count   int
	Percent int
	Color   string
}

type FinalDateEntry struct {
	DateEntries map[string]*DateEntry
	SortedDates []string // Sorted keys from oldest to newest
}

// GetSortedEntries returns DateEntry slice in sorted order (oldest to newest)
func (f *FinalDateEntry) GetSortedEntries() []*DateEntry {
	entries := make([]*DateEntry, 0, len(f.SortedDates))
	for _, date := range f.SortedDates {
		entries = append(entries, f.DateEntries[date])
	}
	return entries
}

// GetReverseSortedEntries returns DateEntry slice in reverse order (newest to oldest)
func (f *FinalDateEntry) GetReverseSortedEntries() []*DateEntry {
	entries := make([]*DateEntry, 0, len(f.SortedDates))
	for i := len(f.SortedDates) - 1; i >= 0; i-- {
		entries = append(entries, f.DateEntries[f.SortedDates[i]])
	}
	return entries
}

// Iterate calls fn for each entry in sorted order (oldest to newest)
func (f *FinalDateEntry) Iterate(fn func(date string, entry *DateEntry)) {
	for _, date := range f.SortedDates {
		fn(date, f.DateEntries[date])
	}
}
type ClientOption func(*GitHubClient)

// WithTimeout sets a custom default timeout for requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *GitHubClient) {
		c.defaultTimeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *GitHubClient) {
		c.httpClient = httpClient
	}
}

func NewGitHubClient(pat string, opts ...ClientOption) *GitHubClient {
	client := &GitHubClient{
		httpClient:          &http.Client{},
		personalAccessToken: pat,
		defaultTimeout:      30 * time.Second, // Default timeout
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *GitHubClient) doRequest(ctx context.Context, url string, acceptHeader string) ([]byte, error) {
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", acceptHeader)
	if c.personalAccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.personalAccessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Check if error is due to context cancellation or timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out: %w", err)
		}
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("request canceled: %w", err)
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// FetchProfile fetches a user's GitHub profile
func (c *GitHubClient) FetchProfile(ctx echo.Context, username string) (*GitHubProfile, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	body, err := c.doRequest(ctx.Request().Context(), url, "application/vnd.github.v3+json")
	if err != nil {
		return nil, err
	}

	var profile GitHubProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &profile, nil
}


// FetchCommits fetches commit history for a user
func (c *GitHubClient) FetchCommits(ctx echo.Context, username string, pages int) (*FinalDateEntry, error) {
	commits := make(map[string]string)
	urlTemplate := "https://api.github.com/search/commits?q=author:%s&sort=author-date&order=desc&page=%d"

	for i := 1; i <= pages; i++ {
		// Check if context is already cancelled before making request
		select {
		case <-ctx.Request().Context().Done():
			return nil, fmt.Errorf("operation cancelled: %w", ctx.Request().Context().Err())
		default:
		}

		url := fmt.Sprintf(urlTemplate, username, i)
		body, err := c.doRequest(ctx.Request().Context(), url, "application/vnd.github.cloak-preview")
		if err != nil {
			// If context cancelled/timed out, return immediately
			if ctx.Request().Context().Err() != nil {
				return nil, fmt.Errorf("operation cancelled after %d pages: %w", i-1, ctx.Request().Context().Err())
			}
			fmt.Printf("Warning: failed to fetch page %d: %v\n", i, err)
			continue
		}

		var searchResp SearchResponse
		if err := json.Unmarshal(body, &searchResp); err != nil {
			fmt.Printf("Warning: failed to parse page %d: %v\n", i, err)
			continue
		}

		for _, item := range searchResp.Items {
			commits[item.SHA] = item.Commit.Author.Date
		}
	}

	return processCommits(commits), nil
}

func processCommits(commits map[string]string) *FinalDateEntry {
	datesDict := make(map[string]int)
	finalDict := FinalDateEntry{
		DateEntries: make(map[string]*DateEntry),
		SortedDates: []string{},
	}

	for _, dateStr := range commits {
		date := strings.Split(dateStr, "T")[0]
		datesDict[date]++
		finalDict.DateEntries[date] = &DateEntry{
			Date:  date,
			Count: datesDict[date],
		}
	}

	maxCommitCount := 0
	for _, count := range datesDict {
		if count > maxCommitCount {
			maxCommitCount = count
		}
	}

	for date, count := range datesDict {
		percent := 0
		if maxCommitCount > 0 {
			percent = int((float64(count) / float64(maxCommitCount)) * 100)
		}

		finalDict.DateEntries[date].Percent = percent
		finalDict.DateEntries[date].Color = blendHex(baseColor, targetColor, percent)
	}

	// Populate SortedDates from DateEntries keys and sort them
	for date := range finalDict.DateEntries {
		finalDict.SortedDates = append(finalDict.SortedDates, date)
	}
	sort.Strings(finalDict.SortedDates)

	return &finalDict
}

func blendHex(hex1, hex2 string, percent int) string {
	hex1 = strings.TrimPrefix(hex1, "#")
	hex2 = strings.TrimPrefix(hex2, "#")

	r1, _ := parseInt(hex1[0:2], 16)
	g1, _ := parseInt(hex1[2:4], 16)
	b1, _ := parseInt(hex1[4:6], 16)

	r2, _ := parseInt(hex2[0:2], 16)
	g2, _ := parseInt(hex2[2:4], 16)
	b2, _ := parseInt(hex2[4:6], 16)

	t := float64(percent) / 100.0
	r := int(float64(r1) + float64(r2-r1)*t)
	g := int(float64(g1) + float64(g2-g1)*t)
	b := int(float64(b1) + float64(b2-b1)*t)

	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func parseInt(s string, base int) (int64, error) {
	var result int64
	for _, c := range s {
		result *= int64(base)
		if c >= '0' && c <= '9' {
			result += int64(c - '0')
		} else if c >= 'a' && c <= 'f' {
			result += int64(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			result += int64(c - 'A' + 10)
		}
	}
	return result, nil
}
