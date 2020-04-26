package vkbot

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type webhookBot struct {
	logger           *logrus.Logger
	router           *router
	groupID          int
	confirmationCode string
}

func newWebhookBot(groupID int, confirmation string) *webhookBot {
	bot := &webhookBot{
		logger:           logrus.New(),
		router:           newRouter(),
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
}

func (b *webhookBot) confirmationHandler() handlerFunc {
	return func(req *request, resp reponse) {
		resp.Write([]byte(b.confirmationCode))
	}
}

func (b *webhookBot) messageHandler() handlerFunc {
	return func(req *request, resp reponse) {

	}
}
