package vkbot

import (
	"errors"
	"net/http"
)

// Start starts vk bot
func Start(config *Config) error {
	if config == nil {
		return errors.New("could not start bot config is nil")
	}

	bot := newWebhookBot(config.GroupID, config.Confirmation)
	return http.ListenAndServe(config.Host, bot)
}
