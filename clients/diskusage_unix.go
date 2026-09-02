//go:build linux

package clients

import (
	"fmt"
	"syscall"
)

// DiskUsage reports filesystem usage for the mount that contains path.
//
// This reads the *Mindscape host's* filesystem, not the remote Coolify
// server's: neither the Coolify REST API nor Sentinel v0.0.22 exposes a disk
// metric, so there is nothing to ask over the network. The numbers are only
// meaningful when Mindscape runs on the same machine as the Coolify server,
// which is why the response labels its source.
func DiskUsage(path string) (*Disk, error) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return nil, fmt.Errorf("could not stat filesystem at %q: %w", path, err)
	}

	blockSize := uint64(st.Bsize)
	total := st.Blocks * blockSize
	// Bavail is what a non-root process may actually use; Bfree includes the
	// reserved blocks root can dip into. `df` computes capacity from Bavail, so
	// match it rather than reporting a percentage users cannot reconcile.
	available := st.Bavail * blockSize
	used := (st.Blocks - st.Bfree) * blockSize

	var usedPercent float64
	if denom := used + available; denom > 0 {
		usedPercent = (float64(used) / float64(denom)) * 100
	}

	return &Disk{
		Path:        path,
		Total:       total,
		Used:        used,
		Available:   available,
		UsedPercent: usedPercent,
	}, nil
}
