package vkbot

import (
	"bytes"
	"chillit-vk-bot/internal/app/places"
	"chillit-vk-bot/internal/app/templates"
	"chillit-vk-bot/internal/app/vkapi"
	"chillit-vk-bot/internal/app/vkdialogs"
	"context"
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
	templatesStorage *templates.Storage
}

func newWebhookBot(
	groupID int,
	confirmation string,
	vkAPIClient *vkapi.Client,
	placesStore places.PlacesStoreClient,
	dialogs vkdialogs.Dialogs,
	templatesStorage *templates.Storage,
) *webhookBot {
	bot := &webhookBot{
		logger:           logrus.New(),
		router:           newRouter(),
		vkAPIClient:      vkAPIClient,
		placesStore:      placesStore,
		groupID:          groupID,
		confirmationCode: confirmation,
		dialogs:          dialogs,
		templatesStorage: templatesStorage,
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

		// Request cities from storage
		citiesResp, err := b.placesStore.GetCities(context.Background(), &places.GetCitiesRequest{
			Amount: 5,
			Offset: 0,
		})

		// Build buttons for vk keyboard
		buttons := make([]*vkapi.KeyboardRow, 0)
		for _, city := range citiesResp.Cities {
			buttons = append(buttons, &vkapi.KeyboardRow{
				&vkapi.KeyboardButton{
					Action: &vkapi.ButtonAction{
						Type:  "text",
						Label: city.GetTitle(),
						Payload: fmt.Sprintf(
							"{\"button\": \"city\", \"object\": {\"title\": \"%s\", \"id\": %d, \"offset\": 0}}",
							city.GetTitle(),
							city.GetId(),
						),
					},
					Color: "primary",
				},
			})
		}

		// Adding admin button
		buttons = append(buttons, &vkapi.KeyboardRow{
			&vkapi.KeyboardButton{
				Action: &vkapi.ButtonAction{
					Type:    "text",
					Label:   "Связаться со администратором",
					Payload: "{\"button\": \"admin\"}",
				},
				Color: "secondary",
			},
		})

		// Build keyboard
		kb := &vkapi.Keyboard{
			OneTime: false,
			Buttons: buttons,
		}

		// Send reply
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
		case "random":
			b.handleRandomButton(mo, &pld)(req, resp)
		}
	}
}

func (b *webhookBot) handleCityButton(mo *messageObject, pld *payload) handlerFunc {
	type cityButtonPayload struct {
		Title string `json:"title"`
		ID    int    `json:"id"`
	}

	return func(req *request, resp reponse) {
		var payload cityButtonPayload
		if err := json.Unmarshal([]byte(pld.Object), &payload); err != nil {
			b.logger.Errorf("could not unmarshal payload '%s'. error: %v", string(pld.Object), err)
			return
		}

		feelingLuckyText, err := b.dialogs.GetText("feeling_lucky")
		if err != nil {
			b.logger.Errorf("could not load dialog: %s", err)
			return
		}

		returnText, err := b.dialogs.GetText("return")
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
							Label: feelingLuckyText,
							Payload: fmt.Sprintf(
								"{\"button\": \"random\", \"object\": {\"city_name\": \"%s\", \"city_id\": %d}}",
								payload.Title,
								payload.ID,
							),
						},
						Color: "primary",
					},
				},
				{
					&vkapi.KeyboardButton{
						Action: &vkapi.ButtonAction{
							Type:  "text",
							Label: returnText,
						},
						Color: "secondary",
					},
				},
			},
		}

		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       fmt.Sprintf("Вы выбрали город '%v'", payload.Title),
			Keyboard:   kb,
		}); err != nil {
			b.logger.Errorf("could not send message id '%d' reason: %v", mo.From, err)
		}
	}
}

func (b *webhookBot) handleAdminButton(mo *messageObject, pld *payload) handlerFunc {
	return func(req *request, resp reponse) {
		text, err := b.dialogs.GetText("contact_admin")
		if err != nil {
			b.logger.Errorf("error handling admin button push for id '%d' error: %v", mo.From, err)
		}

		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       text,
		}); err != nil {
			b.logger.Errorf("could not send message id '%d' reason: %v", mo.From, err)
		}
	}
}

func (b *webhookBot) handleRandomButton(mo *messageObject, pld *payload) handlerFunc {
	type randomButtonPayload struct {
		CityID   int    `json:"city_id"`
		CityName string `json:"city_name"`
	}

	return func(req *request, resp reponse) {
		var payload randomButtonPayload
		if err := json.Unmarshal([]byte(pld.Object), &payload); err != nil {
			b.logger.Errorf("could not unmarshal payload '%s' error: %v", string(pld.Object), err)
			return
		}

		plStoreResp, err := b.placesStore.GetRandomPlaceByCityName(context.Background(), &places.GetRandomPlaceByCityNameRequest{
			CityName: payload.CityName,
		})
		if err != nil {
			b.logger.Errorf("could not get place info from placesstore for city '%s' error: %v", payload.CityName, err)
			return
		}

		msgTmpl, err := b.templatesStorage.Get("random_city")
		if err != nil {
			b.logger.Errorf("could not get template for message error: %v", err)
			return
		}

		var messageTextBuf bytes.Buffer
		if err := msgTmpl.Execute(&messageTextBuf, struct {
			Title       string
			Description string
			Address     string
		}{
			Title:       plStoreResp.GetPlace().GetTitle(),
			Description: plStoreResp.GetPlace().GetDescription(),
			Address:     plStoreResp.GetPlace().GetAddress(),
		}); err != nil {
			b.logger.Errorf("could not execute template for message error: %v", err)
			return
		}

		if err := b.vkAPIClient.SendMessage(&vkapi.Message{
			RecieverID: mo.From,
			Text:       string(messageTextBuf.Bytes()),
		}); err != nil {
			b.logger.Errorf("could not send message id '%d' reason: %v", mo.From, err)
		}
	}
}
