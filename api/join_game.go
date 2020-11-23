package api

import (
	"context"
	"fmt"

	model "lukfugl/no-thanks-server/firestore"

	"cloud.google.com/go/firestore"
)

// A JoinGameAction is an Action that adds a player to an existing game, identified by slug
type JoinGameAction struct {
	PlayerName string
	Slug       string
}

// Execute adds the player to the game
func (a *JoinGameAction) Execute(ctx context.Context, fs *firestore.Client, userID string) ([]byte, error) {
	txErr := fs.RunTransaction(ctx, func(_ctx context.Context, t *firestore.Transaction) error {
		slug := model.FindGameSlug(fs, a.Slug)
		found, err := slug.Load(t)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("unknown slug")
		}
		game := slug.TargetGame
		found, err = game.Load(t)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("slug references unknown game")
		}
		game.AddPlayer(a.PlayerName, userID)
		return game.Save(t)
	})
	return nil, txErr
}
