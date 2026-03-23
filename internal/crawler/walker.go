package crawler

import (
	"io/fs"
	"path/filepath"
)

type WalkEntry struct {
	Path  string
	Entry fs.DirEntry
	Err   error
}

// Walk sends every entry from filepath.WalkDir into the returned channel.
// The channel is closed when the walk finishes or a top-level error occurs.
func Walk(root string) <-chan WalkEntry {
	ch := make(chan WalkEntry, 256)
	go func() {
		defer close(ch)
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			ch <- WalkEntry{Path: path, Entry: d, Err: err}
			return nil
		})
	}()
	return ch
}
