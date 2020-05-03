package vkbot

import (
	"encoding/json"
	"fmt"
)

type payload struct {
	Button string          `json:"button"`
	Object json.RawMessage `json:"object"`
}

func (p *payload) Parse(payloadString string) error {
	if err := json.Unmarshal([]byte(payloadString), p); err != nil {
		return fmt.Errorf("could not unmarshal payload '%s'. error: %v", payloadString, err)
	}
	return nil
}
