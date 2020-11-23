package api

import (
	"context"

	"cloud.google.com/go/firestore"
)

// An Action is a unit of execution for the API
type Action interface {
	Execute(ctx context.Context, fs *firestore.Client, userID string) ([]byte, error)
}
