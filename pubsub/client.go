package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// Initialize initializes the pubsub client
func NewClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID)
}
