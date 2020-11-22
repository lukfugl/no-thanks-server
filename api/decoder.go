package api

import (
	"encoding/json"
	"fmt"
)

type actionDecoder struct {
	CreateGame *CreateGameAction
}

func decodeAction(bytes []byte) (Action, error) {
	var decoder actionDecoder
	err := json.Unmarshal(bytes, &decoder)
	if err != nil {
		return nil, err
	}
	if decoder.CreateGame != nil {
		return decoder.CreateGame, nil
	}
	return nil, fmt.Errorf("unknown action")
}
