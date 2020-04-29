package vkbot

import (
	"chillit-vk-bot/internal/app/places"
	"errors"
	"net/http"
)

// Start starts vk bot
func Start(config *Config, placesStore places.PlacesStoreClient) error {
	if config == nil {
		return errors.New("could not start bot config is nil")
	}

	bot := newWebhookBot(config.GroupID, config.Confirmation, placesStore)
	return http.ListenAndServe(config.Host, bot)
}
