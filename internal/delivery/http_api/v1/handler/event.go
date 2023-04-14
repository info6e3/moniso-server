package handlers

import "net/http"

type eventHandler struct{}

func newEventHandler() *eventHandler {
	return &eventHandler{}
}

func (h *eventHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *eventHandler) Add(w http.ResponseWriter, r *http.Request) {

}
