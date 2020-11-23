package api

import (
	"context"
	"encoding/json"

	model "lukfugl/no-thanks-server/firestore"

	"cloud.google.com/go/firestore"
)

// A CreateGameAction is an Action that creates a new Game, owned by the
// user, along with Redirector pointing at it.
type CreateGameAction struct {
	PlayerName string
}

// Execute creates the game
func (a *CreateGameAction) Execute(ctx context.Context, fs *firestore.Client, userID string) ([]byte, error) {
	var slug *model.GameSlug
	var err error
	txErr := fs.RunTransaction(ctx, func(_ctx context.Context, t *firestore.Transaction) error {
		available := false
		for !available {
			slug = model.NewGameSlug(fs)
			available, err = slug.IsAvailable(t)
			if err != nil {
				return err
			}
		}

		game := model.NewGame(fs, userID)
		game.AddPlayer(a.PlayerName, userID)
		err = game.Save(t)
		if err != nil {
			return err
		}

		slug.TargetGame = game
		return slug.Save(t)
	})
	if txErr != nil {
		return []byte{}, txErr
	}
	return json.Marshal(map[string]string{
		"slug": slug.ID,
	})
}
