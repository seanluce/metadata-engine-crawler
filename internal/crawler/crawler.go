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

	fmt.Printf("Starting crawl of %q with %d workers\n", cfg.Root, cfg.Workers)
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

				rec, err := buildRecord(cfg.Root, entry)
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

func buildRecord(root string, entry WalkEntry) (models.FileRecord, error) {
	info, err := entry.Entry.Info()
	if err != nil {
		return models.FileRecord{}, err
	}

	path := filepath.ToSlash(entry.Path)
	name := info.Name()
	isDir := entry.Entry.IsDir()

	parentPath := filepath.ToSlash(filepath.Dir(entry.Path))
	if parentPath == "." {
		parentPath = ""
	}
	// Root entry: set parent_path to empty string
	if filepath.Clean(entry.Path) == filepath.Clean(root) {
		parentPath = ""
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
