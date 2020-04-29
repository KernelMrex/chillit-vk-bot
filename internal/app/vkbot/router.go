package vkbot

import (
	"encoding/json"
	"net/http"
)

const (
	actionNewMessage   = "message_new"
	actionConfirmation = "confirmation"
)

type router struct {
	routes map[string]*route
}

func newRouter() *router {
	return &router{
		make(map[string]*route),
	}
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyReader := r.Body
	defer bodyReader.Close()

	requestContent := struct {
		Type    string           `json:"type"`
		Object  *json.RawMessage `json:"object"`
		GroupID int              `json:"group_id"`
	}{}
	if err := json.NewDecoder(bodyReader).Decode(&requestContent); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	route, ok := router.routes[requestContent.Type]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	requestStruct := request{
		Object: requestContent.Object,
	}

	responseStruct := responseStruct{
		Body:  make([]byte, 0),
		Error: nil,
	}

	route.Handler(&requestStruct, &responseStruct)

	if responseStruct.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseStruct.Body)
}

func (router *router) HandleFunc(action string, handler handlerFunc) {
	router.routes[action] = &route{
		Handler: handler,
	}
}
