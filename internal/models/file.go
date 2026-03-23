package models

import (
	"time"
)

type FileRecord struct {
	Name           string            `bson:"name"`
	Path           string            `bson:"path"`
	ParentPath     string            `bson:"parent_path"`
	IsDir          bool              `bson:"is_dir"`
	Size           int64             `bson:"size"`
	Extension      string            `bson:"extension"`
	MimeType       string            `bson:"mime_type"`
	CreatedAt      time.Time         `bson:"created_at"`
	ModifiedAt     time.Time         `bson:"modified_at"`
	AccessedAt     time.Time         `bson:"accessed_at"`
	CrawledAt      time.Time         `bson:"crawled_at"`
	Mode           string            `bson:"mode"`
	CustomMetadata map[string]string `bson:"custom_metadata"`
}
