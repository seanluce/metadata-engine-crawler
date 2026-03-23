//go:build !windows && !darwin

package crawler

import (
	"os"
	"syscall"
	"time"
)

type Timestamps struct {
	Created  time.Time
	Modified time.Time
	Accessed time.Time
}

func getTimestamps(info os.FileInfo) Timestamps {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t := info.ModTime()
		return Timestamps{Created: t, Modified: t, Accessed: t}
	}
	mod := time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).UTC()
	acc := time.Unix(stat.Atim.Sec, stat.Atim.Nsec).UTC()
	return Timestamps{
		Created:  mod, // Linux does not expose birthtime; use mtime as proxy
		Modified: mod,
		Accessed: acc,
	}
}
