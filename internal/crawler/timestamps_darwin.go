//go:build darwin

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
	return Timestamps{
		Created:  time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec).UTC(),
		Modified: time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec).UTC(),
		Accessed: time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec).UTC(),
	}
}
