package crawler

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seanluce/metadata-engine/crawler/internal/api"
	"github.com/seanluce/metadata-engine/crawler/internal/config"
	mimepkg "github.com/seanluce/metadata-engine/crawler/internal/mime"
	"github.com/seanluce/metadata-engine/crawler/internal/models"
)

func Run(cfg config.Config) error {
	client := api.New(cfg.ApiURL)

	// Default volume name to the root directory name
	volumeName := cfg.VolumeName
	if volumeName == "" {
		volumeName = filepath.Base(filepath.Clean(cfg.Root))
	}

	fmt.Printf("Starting crawl of %q as volume %q with %d workers\n", cfg.Root, volumeName, cfg.Workers)
	start := time.Now()

	entries := Walk(cfg.Root)

	var wg sync.WaitGroup
	var count atomic.Int64
	var errCount atomic.Int64

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for entry := range entries {
				if entry.Err != nil {
					fmt.Printf("WARN: walk error at %s: %v\n", entry.Path, entry.Err)
					errCount.Add(1)
					continue
				}

				rec, err := buildRecord(cfg.Root, volumeName, entry)
				if err != nil {
					fmt.Printf("WARN: stat error at %s: %v\n", entry.Path, err)
					errCount.Add(1)
					continue
				}

				if err := client.Insert(context.Background(), rec); err != nil {
					fmt.Printf("WARN: insert error at %s: %v\n", entry.Path, err)
					errCount.Add(1)
					continue
				}

				n := count.Add(1)
				if n%500 == 0 {
					fmt.Printf("  ... processed %d files\n", n)
				}
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Crawl complete: %d files in %v (%d errors)\n",
		count.Load(), time.Since(start).Round(time.Millisecond), errCount.Load())
	return nil
}

func buildRecord(root string, volumeName string, entry WalkEntry) (models.FileRecord, error) {
	info, err := entry.Entry.Info()
	if err != nil {
		return models.FileRecord{}, err
	}

	// Normalize everything to forward slashes first, then strip the
	// crawl root so stored paths are relative and consistent across platforms.
	cleanRoot := filepath.ToSlash(filepath.Clean(root))
	fullPath := filepath.ToSlash(entry.Path)

	// Build the relative path under the volume name.
	// e.g. root=C:\Users\sean\Downloads, volumeName=Downloads
	//   C:/Users/sean/Downloads       → /Downloads
	//   C:/Users/sean/Downloads/a.txt → /Downloads/a.txt
	relPath := strings.TrimPrefix(fullPath, cleanRoot)
	var path string
	if relPath == "" {
		// This is the root entry itself
		path = "/" + volumeName
	} else {
		path = "/" + volumeName + relPath
	}

	name := info.Name()
	if relPath == "" {
		// Root entry uses the volume name
		name = volumeName
	}
	isDir := entry.Entry.IsDir()

	// Parent path: strip the root and prepend the volume name
	parentFull := filepath.ToSlash(filepath.Dir(entry.Path))
	parentRel := strings.TrimPrefix(parentFull, cleanRoot)
	var parentPath string
	if fullPath == cleanRoot {
		// Root entry has no parent
		parentPath = ""
	} else if parentRel == "" {
		// Direct child of root
		parentPath = "/" + volumeName
	} else {
		parentPath = "/" + volumeName + parentRel
	}

	ext := ""
	if !isDir {
		ext = strings.ToLower(filepath.Ext(name))
	}

	timestamps := getTimestamps(info)

	return models.FileRecord{
		Name:       name,
		Path:       path,
		ParentPath: parentPath,
		IsDir:      isDir,
		Size:       info.Size(),
		Extension:  ext,
		MimeType:   mimepkg.TypeByExtension(ext),
		CreatedAt:  timestamps.Created,
		ModifiedAt: timestamps.Modified,
		AccessedAt: timestamps.Accessed,
		CrawledAt:  time.Now().UTC(),
		Mode:       info.Mode().String(),
	}, nil
}
