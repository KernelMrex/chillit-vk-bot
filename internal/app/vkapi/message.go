package vkapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Message structure of message data in request
type Message struct {
	RecieverID int
	RandomID   int64
	Text       string
	Keyboard   *Keyboard
}

func (m *Message) urlEncode() (string, error) {
	params := url.Values{}
	params.Add("user_id", strconv.Itoa(m.RecieverID))
	params.Add("random_id", strconv.FormatInt(m.RandomID, 10))
	params.Add("message", m.Text)
	if m.Keyboard != nil {
		var err error
		encodedKB, err := json.Marshal(m.Keyboard)
		if err != nil {
			return "", fmt.Errorf("error while encoding keyboard: %v", err)
		}
		params.Add("keyboard", string(encodedKB))
	}
	return params.Encode(), nil
}
