package api

import (
	"encoding/json"
	"fmt"
)

type actionDecoder struct {
	Echo *EchoAction `json:"echo"`
}

func decodeAction(bytes []byte) (Action, error) {
	var decoder actionDecoder
	err := json.Unmarshal(bytes, &decoder)
	if err != nil {
		return nil, err
	}
	if decoder.Echo != nil {
		return decoder.Echo, nil
	}
	return nil, fmt.Errorf("unknown action")
}
