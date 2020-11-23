package api

import (
	"encoding/json"
	"fmt"
)

type actionDecoder struct {
	CreateGame *CreateGameAction
	JoinGame   *JoinGameAction
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
	if decoder.JoinGame != nil {
		return decoder.JoinGame, nil
	}
	return nil, fmt.Errorf("unknown action")
}
