package vkbot

import (
	"encoding/json"
)

type handlerFunc func(req *request, resp reponse)

type route struct {
	Handler handlerFunc
}

type request struct {
	Object *json.RawMessage
}

type reponse interface {
	Write([]byte) (int, error)
	SetError(error)
}

type responseStruct struct {
	Body  []byte
	Error error
}

func (r *responseStruct) Write(written []byte) (int, error) {
	r.Body = append(r.Body, written...)
	return len(written), nil
}

func (r *responseStruct) SetError(isError error) {
	r.Error = isError
}
