package widget_handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// StringList is a config value that accepts either a single JSON string or an
// array of them.
//
// It exists so a field can grow from one value to many without breaking rows
// already in the database: `json.Unmarshal` of "/" into a plain []string fails
// with "cannot unmarshal string into []string", which would 400 every widget
// configured before the field became a list.
type StringList []string

func (s *StringList) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*s = nil
		return nil
	}

	if trimmed[0] == '[' {
		var list []string
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return err
		}
		*s = list
		return nil
	}

	var single string
	if err := json.Unmarshal(trimmed, &single); err != nil {
		return fmt.Errorf("expected a string or an array of strings: %w", err)
	}
	*s = StringList{single}
	return nil
}

// Clean trims each entry, drops the empty ones, and removes duplicates while
// preserving order. Duplicates matter here because two identical paths would
// otherwise render as two identical tiles.
func (s StringList) Clean() []string {
	seen := make(map[string]struct{}, len(s))
	cleaned := make([]string, 0, len(s))

	for _, entry := range s {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		cleaned = append(cleaned, trimmed)
	}

	return cleaned
}
