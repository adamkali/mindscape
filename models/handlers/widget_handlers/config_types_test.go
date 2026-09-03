package widget_handlers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The whole point of StringList: widgets stored before diskPath became a list
// hold a bare string, and a plain []string field would reject them outright.
func TestStringListAcceptsLegacySingleString(t *testing.T) {
	var config CoolifyWidgetConfig

	err := json.Unmarshal([]byte(`{"diskPath":"/"}`), &config)

	assert.NoError(t, err)
	assert.Equal(t, []string{"/"}, config.DiskPaths.Clean())
}

func TestStringListAcceptsList(t *testing.T) {
	var config CoolifyWidgetConfig

	err := json.Unmarshal([]byte(`{"diskPath":["/","/home","/mnt/data"]}`), &config)

	assert.NoError(t, err)
	assert.Equal(t, []string{"/", "/home", "/mnt/data"}, config.DiskPaths.Clean())
}

// An absent or null key must not error; the handler falls back to "/".
func TestStringListAbsentAndNull(t *testing.T) {
	for _, payload := range []string{`{}`, `{"diskPath":null}`} {
		var config CoolifyWidgetConfig

		err := json.Unmarshal([]byte(payload), &config)

		assert.NoError(t, err, payload)
		assert.Empty(t, config.DiskPaths.Clean(), payload)
	}
}

func TestStringListCleanTrimsDropsEmptyAndDedupes(t *testing.T) {
	var config CoolifyWidgetConfig

	err := json.Unmarshal(
		[]byte(`{"diskPath":["  /  ","","/home","/",  "/home  "]}`),
		&config,
	)

	assert.NoError(t, err)
	// Order preserved, whitespace trimmed, blanks dropped, repeats collapsed —
	// a duplicate would otherwise render as a second identical tile.
	assert.Equal(t, []string{"/", "/home"}, config.DiskPaths.Clean())
}

func TestStringListRejectsWrongType(t *testing.T) {
	var config CoolifyWidgetConfig

	err := json.Unmarshal([]byte(`{"diskPath":42}`), &config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "string or an array of strings")
}

// The other Coolify config keys must survive alongside the list field.
func TestCoolifyConfigRoundTripWithLegacyDiskPath(t *testing.T) {
	var config CoolifyWidgetConfig

	err := json.Unmarshal([]byte(`{
		"baseUrl":"https://coolify.example.com/api/v1",
		"personalAccessToken":"tok",
		"serverUuid":"abc",
		"sentinelUrl":"http://coolify-sentinel:8888",
		"sentinelToken":"stok",
		"diskPath":"/",
		"applicationIds":[],
		"serviceIds":[]
	}`), &config)

	assert.NoError(t, err)
	assert.Equal(t, "https://coolify.example.com/api/v1", config.BaseUrl)
	assert.Equal(t, "http://coolify-sentinel:8888", config.SentinelUrl)
	assert.Equal(t, []string{"/"}, config.DiskPaths.Clean())
}
