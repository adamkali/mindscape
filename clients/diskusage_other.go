//go:build !linux

package clients

import (
	"fmt"
	"runtime"
)

// DiskUsage is unsupported outside Linux. Mindscape ships in a Linux
// container; this exists so the package still builds elsewhere rather than
// pretending to have a number.
func DiskUsage(path string) (*Disk, error) {
	return nil, fmt.Errorf("disk usage is not supported on %s", runtime.GOOS)
}
