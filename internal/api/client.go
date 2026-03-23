package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/seanluce/metadata-engine/crawler/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type ingestPayload struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	ParentPath string `json:"parent_path"`
	IsDir      bool   `json:"is_dir"`
	Size       int64  `json:"size"`
	Extension  string `json:"extension"`
	MimeType   string `json:"mime_type"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	AccessedAt string `json:"accessed_at"`
	CrawledAt  string `json:"crawled_at"`
	Mode       string `json:"mode"`
}

func (c *Client) Insert(ctx context.Context, rec models.FileRecord) error {
	payload := ingestPayload{
		Name:       rec.Name,
		Path:       rec.Path,
		ParentPath: rec.ParentPath,
		IsDir:      rec.IsDir,
		Size:       rec.Size,
		Extension:  rec.Extension,
		MimeType:   rec.MimeType,
		CreatedAt:  rec.CreatedAt.Format(time.RFC3339Nano),
		ModifiedAt: rec.ModifiedAt.Format(time.RFC3339Nano),
		AccessedAt: rec.AccessedAt.Format(time.RFC3339Nano),
		CrawledAt:  rec.CrawledAt.Format(time.RFC3339Nano),
		Mode:       rec.Mode,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/files/ingest", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("api returned %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Disconnect() {
	// no-op for HTTP client
}
