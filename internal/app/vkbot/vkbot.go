package vkbot

import (
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/templates"
	"chillit-vk-bot/internal/app/vkapi"
	"chillit-vk-bot/internal/app/vkdialogs"
	"fmt"
	"net/http"
)

// Start starts vk bot
func Start(config *Config, placesStore places.PlacesStoreClient, templatesStorage *templates.Storage) error {
	if config == nil {
		return fmt.Errorf("bot config is nil")
	}

	dialogs, err := vkdialogs.LoadFromFile(config.DialogsPath)
	if err != nil {
		return fmt.Errorf("could not load vk dialogs: %v", err)
	}

	bot := newWebhookBot(
		config.GroupID,
		config.Confirmation,
		vkapi.NewClient(config.Token),
		placesStore,
		dialogs,
		templatesStorage,
	)

	return http.ListenAndServe(config.Host, bot)
}
