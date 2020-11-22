package api

import (
	"context"

	"cloud.google.com/go/firestore"
)

// An EchoAction is an example ApiAction that simply logs the provided message
type EchoAction struct {
	Message string `json:"message"`
}

// Execute logs the provided message
func (a *EchoAction) Execute(ctx context.Context, fs *firestore.Client) error {
	messages := fs.Collection("Messages")
	_, _, err := messages.Add(ctx, map[string]interface{}{
		"message": a.Message,
	})
	return err
}
