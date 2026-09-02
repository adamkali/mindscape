package schemas

import (
	"strings"
	"testing"
)

// The Coolify schema stays embedded but must be held out of the default
// install list unless explicitly opted in via config `features.coolify`.
func TestCoolifySchemaGating(t *testing.T) {
	t.Cleanup(func() { SetCoolifyEnabled(false) })

	hasCoolify := func(t *testing.T) bool {
		t.Helper()
		storage, err := EmbeddedScemas()
		if err != nil {
			t.Fatalf("EmbeddedScemas: %v", err)
		}
		for _, schema := range storage.GetAll() {
			if strings.HasPrefix(schema.Type, "coolify") {
				return true
			}
		}
		return false
	}

	SetCoolifyEnabled(false)
	if hasCoolify(t) {
		t.Error("coolify schema registered by default; expected it to be opt-in")
	}

	SetCoolifyEnabled(true)
	if !hasCoolify(t) {
		t.Error("coolify schema missing after SetCoolifyEnabled(true)")
	}
}
