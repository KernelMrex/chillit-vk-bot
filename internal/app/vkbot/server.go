package vkbot

import (
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/vkapi"
	"chillit-vk-bot/internal/app/vkdialogs"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type webhookBot struct {
	logger           *logrus.Logger
	router           *router
	vkAPIClient      *vkapi.Client
	placesStore      places.PlacesStoreClient
	groupID          int
	confirmationCode string
	dialogs          vkdialogs.Dialogs
}

func newWebhookBot(groupID int, confirmation string, vkAPIClient *vkapi.Client, placesStore places.PlacesStoreClient, dialogs vkdialogs.Dialogs) *webhookBot {
	bot := &webhookBot{
		logger:           logrus.New(),
		router:           newRouter(),
		vkAPIClient:      vkAPIClient,
		placesStore:      placesStore,
		groupID:          groupID,
		confirmationCode: confirmation,
		dialogs:          dialogs,
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
			return
		}

		// Logging
		b.logger.Infof("Got message '%s' from user id '%d'", messageAct.Message.Text, messageAct.Message.From)

		// Handling default message
		if messageAct.Message.Payload == "" {
			messageText, err := b.dialogs.GetText("greeting")
			if err != nil {
				b.logger.Errorf("could not load dialog: %s", err)
				return
			}

			kb := &vkapi.Keyboard{
				OneTime: false,
				Buttons: []*vkapi.KeyboardRow{
					{
						&vkapi.KeyboardButton{
							Action: &vkapi.ButtonAction{
								Type:  "text",
								Label: "Город 1",
							},
							Color: "secondary",
						},
						&vkapi.KeyboardButton{
							Action: &vkapi.ButtonAction{
								Type:  "text",
								Label: "Город 2",
							},
							Color: "secondary",
						},
					},
					{
						&vkapi.KeyboardButton{
							Action: &vkapi.ButtonAction{
								Type:  "text",
								Label: "Город 3",
							},
							Color: "secondary",
						},
						&vkapi.KeyboardButton{
							Action: &vkapi.ButtonAction{
								Type:  "text",
								Label: "Город 4",
							},
							Color: "secondary",
						},
					},
					{
						&vkapi.KeyboardButton{
							Action: &vkapi.ButtonAction{
								Type:    "text",
								Label:   "Связаться со администратором",
								Payload: "{\"button\": \"1\"}",
							},
							Color: "secondary",
						},
					},
				},
			}

			if err := b.vkAPIClient.SendMessage(&vkapi.Message{
				RecieverID: messageAct.Message.From,
				Text:       messageText,
				Keyboard:   kb,
			}); err != nil {
				b.logger.Errorf("Could not send message id '%d' reason: %v", messageAct.Message.From, err)
			}
		}

		// TODO: payload handling

		resp.Write([]byte("ok"))
	}
}
