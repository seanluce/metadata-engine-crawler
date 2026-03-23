//go:build windows

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
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		t := info.ModTime()
		return Timestamps{Created: t, Modified: t, Accessed: t}
	}
	return Timestamps{
		Created:  time.Unix(0, stat.CreationTime.Nanoseconds()).UTC(),
		Modified: time.Unix(0, stat.LastWriteTime.Nanoseconds()).UTC(),
		Accessed: time.Unix(0, stat.LastAccessTime.Nanoseconds()).UTC(),
	}
}
