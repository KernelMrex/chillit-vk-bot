package vkbot

import (
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/vkapi"
	"chillit-vk-bot/internal/app/vkdialogs"
	"encoding/json"
	"fmt"
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

		if messageAct.Message.Payload != "" {
			// Handle payload
			b.messagePayloadHandler(messageAct.Message)(req, resp)
		} else {
			// Handle default message
			b.messageOnlyTextHandler(messageAct.Message)(req, resp)
		}

		resp.Write([]byte("ok"))
	}
}

func (b *webhookBot) messageOnlyTextHandler(mo *messageObject) handlerFunc {
	return func(req *request, resp reponse) {
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
							Type:    "text",
							Label:   "Казань",
							Payload: "{\"button\": \"city\", \"object\": {\"title\": \"казань\", \"offset\": 0}}",
						},
						Color: "primary",
					},
				},
				{
					&vkapi.KeyboardButton{
						Action: &vkapi.ButtonAction{
							Type:    "text",
							Label:   "Связаться со администратором",
							Payload: "{\"button\": \"admin\"}",
						},
						Color: "secondary",
					},
				},
			},
		}

		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       messageText,
			Keyboard:   kb,
		}); err != nil {
			b.logger.Errorf("Could not send message id '%d' reason: %v", mo.From, err)
		}
	}
}

func (b *webhookBot) messagePayloadHandler(mo *messageObject) handlerFunc {
	return func(req *request, resp reponse) {
		var pld payload
		if err := pld.Parse(mo.Payload); err != nil {
			b.logger.Errorf("could not handle message with payload: %v", err)
			return
		}

		switch pld.Button {
		case "city":
			b.handleCityButton(mo, &pld)(req, resp)
		case "admin":
			b.handleAdminButton(mo, &pld)(req, resp)
		}
	}
}

func (b *webhookBot) handleCityButton(mo *messageObject, pld *payload) handlerFunc {
	type cityButtonPayload struct {
		Title string `json:"title"`
	}

	return func(req *request, resp reponse) {
		var payload cityButtonPayload
		if err := json.Unmarshal([]byte(pld.Object), &payload); err != nil {
			b.logger.Errorf("could not unmarshal payload '%s'. error: %v", string(pld.Object), err)
			return
		}

		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       fmt.Sprintf("Вы выбрали город '%v'", payload.Title),
		}); err != nil {
			b.logger.Errorf("could not send message id '%d' reason: %v", mo.From, err)
		}

		// TODO: places preview
	}
}

func (b *webhookBot) handleAdminButton(mo *messageObject, pld *payload) handlerFunc {
	return func(req *request, resp reponse) {
		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       fmt.Sprintf("Администратор и программист в едином лице https://vk.com/versarter"),
		}); err != nil {
			b.logger.Errorf("could not send message id '%d' reason: %v", mo.From, err)
		}
	}
}
