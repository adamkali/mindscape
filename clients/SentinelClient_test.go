package clients

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// echoCtx builds the minimal echo.Context the client methods need, since they
// only ever reach for the request's context.
func echoCtx(t *testing.T) echo.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func TestSentinelClientCurrentCPU(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/cpu/current", r.URL.Path)
		assert.Equal(t, "Bearer sentinel-token", r.Header.Get("Authorization"))
		w.Write([]byte(`{"percent":14.655988977869095,"time":"1788379187033"}`))
	}))
	defer server.Close()

	cpu, err := NewSentinelClient("sentinel-token", server.URL).CurrentCPU(echoCtx(t))

	assert.NoError(t, err)
	assert.InDelta(t, 14.655988977869095, cpu.Percent, 0.000001)
	assert.Equal(t, "1788379187033", cpu.Time)
}

func TestSentinelClientCurrentMemory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/memory/current", r.URL.Path)
		w.Write([]byte(`{"time":"1788379187038","total":67333689344,"available":59373621248,"used":6977318912,"usedPercent":10.36,"free":43386630144}`))
	}))
	defer server.Close()

	mem, err := NewSentinelClient("t", server.URL).CurrentMemory(echoCtx(t))

	assert.NoError(t, err)
	assert.Equal(t, uint64(67333689344), mem.Total)
	assert.Equal(t, uint64(6977318912), mem.Used)
	assert.InDelta(t, 10.36, mem.UsedPercent, 0.001)
}

// A trailing slash in configuration must not produce `//api/cpu/current`.
func TestSentinelClientTrimsTrailingSlash(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Write([]byte(`{"percent":1,"time":"0"}`))
	}))
	defer server.Close()

	client := NewSentinelClient("t", server.URL+"/  ")
	_, err := client.CurrentCPU(echoCtx(t))

	assert.NoError(t, err)
	assert.Equal(t, "/api/cpu/current", gotPath)
}

// Sentinel authenticates in middleware ahead of routing, so a bad token turns
// every path into a 401. The error has to name the token, not the route.
func TestSentinelClientUnauthorizedNamesTheToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := NewSentinelClient("wrong", server.URL).CurrentCPU(echoCtx(t))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sentinelToken")
}

func TestSentinelClientHealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/health", r.URL.Path)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	assert.NoError(t, NewSentinelClient("t", server.URL).HealthCheck(echoCtx(t)))
}

func TestDiskUsageReportsRootFilesystem(t *testing.T) {
	disk, err := DiskUsage("/")

	assert.NoError(t, err)
	assert.Equal(t, "/", disk.Path)
	assert.Greater(t, disk.Total, uint64(0))
	// used + available is the `df` denominator and cannot exceed the raw total,
	// which also counts root's reserved blocks.
	assert.LessOrEqual(t, disk.Used+disk.Available, disk.Total)
	assert.GreaterOrEqual(t, disk.UsedPercent, 0.0)
	assert.LessOrEqual(t, disk.UsedPercent, 100.0)
}

func TestDiskUsageRejectsMissingPath(t *testing.T) {
	_, err := DiskUsage("/nonexistent-path-for-mindscape-test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not stat filesystem")
}

func TestGithubAPIErrorExplains401(t *testing.T) {
	err := githubAPIError(http.StatusUnauthorized, []byte(`{"message":"Bad credentials"}`))

	assert.Contains(t, err.Error(), "Bad credentials")
	assert.Contains(t, err.Error(), "expired, revoked, or mistyped")
}

func TestGithubAPIErrorFallsBackToRawBody(t *testing.T) {
	err := githubAPIError(http.StatusBadGateway, []byte("upstream boom"))

	assert.Contains(t, err.Error(), "502")
	assert.Contains(t, err.Error(), "upstream boom")
}

func TestCoolifyAPIErrorExplains401(t *testing.T) {
	err := coolifyAPIError(http.StatusUnauthorized, []byte(`{"message":"Unauthenticated."}`))

	assert.Contains(t, err.Error(), "Unauthenticated.")
	assert.True(t, strings.Contains(err.Error(), "expired"))
}
