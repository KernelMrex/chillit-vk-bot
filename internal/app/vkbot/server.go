package vkbot

import (
	"chillit-vk-bot/internal/app/places"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type webhookBot struct {
	logger           *logrus.Logger
	router           *router
	placesStore      places.PlacesStoreClient
	groupID          int
	confirmationCode string
}

func newWebhookBot(groupID int, confirmation string, placesStore places.PlacesStoreClient) *webhookBot {
	bot := &webhookBot{
		logger:           logrus.New(),
		router:           newRouter(),
		placesStore:      placesStore,
		groupID:          groupID,
		confirmationCode: confirmation,
	}
	bot.configureRoutes()
	return bot
}

func (b *webhookBot) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.router.ServeHTTP(w, r)
}

func (b *webhookBot) configureRoutes() {
	b.router.HandleFunc(actionConfirmation, b.confirmationHandler())
	b.router.HandleFunc(actionNewMessage, b.messageHandler())
}

func (b *webhookBot) confirmationHandler() handlerFunc {
	return func(req *request, resp reponse) {
		resp.Write([]byte(b.confirmationCode))
	}
}

func (b *webhookBot) messageHandler() handlerFunc {
	return func(req *request, resp reponse) {
		var messageAct messageActionObject
		if err := json.Unmarshal(*req.Object, &messageAct); err != nil {
			b.logger.Errorf("could not unmarshal messageActionObject in messageHandler: '%v'", err)
			resp.Write([]byte("ok"))
			return
		}

		// TODO: message handling logic

		resp.Write([]byte("ok"))
	}
}
