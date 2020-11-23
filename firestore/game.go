package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A Player describes a player in a Game
type Player struct {
	UserID string
	Name   string
}

// A Game encapsulates interacting with a Game document in firestore.
type Game struct {
	fs         *firestore.Client
	ref        *firestore.DocumentRef
	Path       string
	HostUserID string
	GameState  string
	Players    []Player
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
	playersRaw, err := doc.DataAt("players")
	if err != nil {
		// current document is malformed, treat as unfound
		return false, fmt.Errorf("playersRaw from doc.DataAt failed")
	}
	playersMap, ok := playersRaw.(map[string]interface{})
	if !ok {
		// current document is malformed, treat as unfound
		return false, fmt.Errorf("playersRaw as map[string]interface{} failed")
	}
	players := []Player{}
	for userID, playerRaw := range playersMap {
		player, ok := playerRaw.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("playerRaw as map[string]interface{} failed")
		}
		name, ok := player["name"].(string)
		if !ok {
			return false, fmt.Errorf("player[name] as string failed")
		}
		players = append(players, Player{
			UserID: userID,
			Name:   name,
		})
	}

	game.HostUserID = hostUserID
	game.GameState = gameState
	game.Players = players
	return true, nil
}

// AddPlayer adds a player to a game
func (game *Game) AddPlayer(name string, userID string) {
	game.Players = append(game.Players, Player{
		UserID: userID,
		Name:   name,
	})
}

func serializePlayers(players []Player) map[string]interface{} {
	serialized := map[string]interface{}{}
	for _, player := range players {
		serialized[player.UserID] = map[string]string{
			"name": player.Name,
		}
	}
	return serialized
}

// Save saves the contents of the game back to firestore
func (game *Game) Save(t *firestore.Transaction) error {
	return t.Set(game.ref, map[string]interface{}{
		"hostUserId": game.HostUserID,
		"gameState":  game.GameState,
		"players":    serializePlayers(game.Players),
	})
}
