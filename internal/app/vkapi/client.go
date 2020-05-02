package vkapi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	version string = "5.103"

	endpoint          string = "https://api.vk.com/method/"
	methodMessageSend string = "messages.send"
)

// Client ...
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient ...
func NewClient(token string) *Client {
	rand.Seed(time.Now().UnixNano())
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Millisecond * 2000,
		},
		token: token,
	}
}

// SendMessage sends message of type Message
func (c *Client) SendMessage(m *Message) error {
	m.RandomID = rand.Int63()

	methodParams, err := m.urlEncode()
	if err != nil {
		return fmt.Errorf("error while encoding method params: %v", err)
	}

	resp, err := c.httpClient.Get(
		fmt.Sprintf("%s?%s&access_token=%s&v=%s", endpoint+methodMessageSend, methodParams, c.token, version),
	)
	if err != nil {
		return fmt.Errorf("error while executing request: %v", err)
	}

	var parsedResp Response
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return fmt.Errorf("error while parsing response: %v", err)
	}

	if parsedResp.Error != nil {
		return fmt.Errorf("error returned from vk: %d %v", parsedResp.Error.Code, parsedResp.Error.Message)
	}

	return nil
}
