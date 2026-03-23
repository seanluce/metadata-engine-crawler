package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seanluce/metadata-engine/crawler/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "files"

type Client struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

func New(uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	c, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	database := c.Database(dbName)
	coll := database.Collection(collectionName)

	if err := ensureIndexes(ctx, coll); err != nil {
		return nil, err
	}

	return &Client{client: c, db: database, coll: coll}, nil
}

func ensureIndexes(ctx context.Context, coll *mongo.Collection) error {
	// Drop the old unique-path index if it exists (schema migration to historical model).
	// The default name for {path:1} unique is "path_1".
	if _, err := coll.Indexes().DropOne(ctx, "path_1"); err != nil {
		// NamespaceNotFound (26) or IndexNotFound (27) are both fine — index just doesn't exist yet.
		var cmdErr mongo.CommandError
		if !errors.As(err, &cmdErr) || (cmdErr.Code != 26 && cmdErr.Code != 27) {
			return fmt.Errorf("drop old path index: %w", err)
		}
	}

	indexes := []mongo.IndexModel{
		{
			// Compound unique: one snapshot per path per crawl run.
			// Also covers single-field queries on path (compound index prefix rule).
			Keys:    bson.D{{Key: "path", Value: 1}, {Key: "crawled_at", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("path_crawled_at"),
		},
		{
			Keys: bson.D{{Key: "parent_path", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_dir", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "name", Value: "text"}},
			Options: options.Index().SetName("name_text"),
		},
	}
	_, err := coll.Indexes().CreateMany(ctx, indexes)
	return err
}

// Insert writes a new snapshot for rec. It first copies custom_metadata forward
// from the most recent previous snapshot for the same path so user annotations
// are never lost across crawl runs.
func (c *Client) Insert(ctx context.Context, rec models.FileRecord) error {
	// Look up the latest existing snapshot for this path to carry metadata forward.
	var prev struct {
		CustomMetadata map[string]string `bson:"custom_metadata"`
	}
	err := c.coll.FindOne(
		ctx,
		bson.D{{Key: "path", Value: rec.Path}},
		options.FindOne().
			SetSort(bson.D{{Key: "crawled_at", Value: -1}}).
			SetProjection(bson.D{{Key: "custom_metadata", Value: 1}}),
	).Decode(&prev)

	switch {
	case err == nil && len(prev.CustomMetadata) > 0:
		rec.CustomMetadata = prev.CustomMetadata
	case err == nil || errors.Is(err, mongo.ErrNoDocuments):
		rec.CustomMetadata = map[string]string{}
	default:
		return fmt.Errorf("copy-forward query: %w", err)
	}

	_, err = c.coll.InsertOne(ctx, rec)
	return err
}

func (c *Client) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = c.client.Disconnect(ctx)
}
