package firestore

import (
	"math/rand"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var slugAlphabet = []byte{'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'}

// A GameSlug represents a redirection from a game code ("slug") to a particular game.
type GameSlug struct {
	fs         *firestore.Client
	ref        *firestore.DocumentRef
	Path       string
	ID         string
	TargetGame *Game
}

// NewGameSlug creates a new, random GameSlug
func NewGameSlug(fs *firestore.Client) *GameSlug {
	indices := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})
	slugBytes := make([]byte, 4)
	for i, k := range indices[0:4] {
		slugBytes[i] = slugAlphabet[k]
	}
	slug := string(slugBytes)
	path := "GameSlug/" + slug
	return &GameSlug{
		fs:         fs,
		ref:        fs.Doc(path),
		Path:       path,
		ID:         slug,
		TargetGame: nil,
	}
}

// Load a game slug from firestore within the transaction. Returns whether the
// document for the slug was found.
func (slug *GameSlug) Load(t *firestore.Transaction) (bool, error) {
	doc, err := t.Get(slug.ref)
	if status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	targetRaw, err := doc.DataAt("target")
	if err != nil {
		// current document is malformed, treat as unfound
		return false, nil
	}
	target, ok := targetRaw.(string)
	if !ok {
		// current document is malformed, treat as unfound
		return false, nil
	}
	slug.TargetGame = &Game{
		fs:   slug.fs,
		ref:  slug.fs.Doc(target),
		Path: target,
	}
	return true, nil
}

// IsAvailable checks if a GameSlug is available within the given transaction
func (slug *GameSlug) IsAvailable(t *firestore.Transaction) (bool, error) {
	slugExists, err := slug.Load(t)
	if err != nil {
		return false, err
	}
	if !slugExists {
		// redirector not yet used
		return true, nil
	}
	gameExists, err := slug.TargetGame.Load(t)
	if !gameExists {
		// current redirect target is invalid
		return true, nil
	}
	// is current target game finished?
	return slug.TargetGame.GameState == "finished", nil
}

// Save saves the contents of the slug back to firestore
func (slug *GameSlug) Save(t *firestore.Transaction) error {
	return t.Set(slug.ref, map[string]string{
		"target": slug.TargetGame.Path,
	})
}
