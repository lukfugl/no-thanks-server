package firestore

import (
	"context"
	"lukfugl/no-thanks-server/config"

	"cloud.google.com/go/firestore"
)

// MakeClient creates and returns a firestore client according to the provided
// configuration
func MakeClient(ctx context.Context, cfg *config.FirestoreConfig) (*firestore.Client, error) {
	return firestore.NewClient(ctx, cfg.ProjectID, cfg.ServiceAccountKey)
}
