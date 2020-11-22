package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A Game encapsulates interacting with a Game document in firestore.
type Game struct {
	fs         *firestore.Client
	ref        *firestore.DocumentRef
	Path       string
	HostUserID string
	GameState  string
}

// NewGame creates a document for a new game hosted by the user
func NewGame(fs *firestore.Client, hostUserID string) *Game {
	path := "Game/" + xid.New().String()
	return &Game{
		fs:         fs,
		ref:        fs.Doc(path),
		Path:       path,
		HostUserID: hostUserID,
		GameState:  "starting",
	}
}

// Load a game from firestore within the transaction. Returns whether the
// document for the game was found.
func (game *Game) Load(t *firestore.Transaction) (bool, error) {
	doc, err := t.Get(game.ref)
	if status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	hostUserIDRaw, err := doc.DataAt("hostUserId")
	if err != nil {
		// current document is malformed, treat as unfound
		return false, nil
	}
	hostUserID, ok := hostUserIDRaw.(string)
	if !ok {
		// current document is malformed, treat as unfound
		return false, nil
	}
	gameStateRaw, err := doc.DataAt("gameState")
	if err != nil {
		// current document is malformed, treat as unfound
		return false, nil
	}
	gameState, ok := gameStateRaw.(string)
	if !ok {
		// current document is malformed, treat as unfound
		return false, nil
	}
	game.HostUserID = hostUserID
	game.GameState = gameState
	return true, nil
}

// Save saves the contents of the game back to firestore
func (game *Game) Save(t *firestore.Transaction) error {
	return t.Set(game.ref, map[string]string{
		"hostUserId": game.HostUserID,
		"gameState":  game.GameState,
	})
}
