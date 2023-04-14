package handlers

import "net/http"

type pointHandler struct{}

func newPointHandler() *pointHandler {
	return &pointHandler{}
}

func (h *pointHandler) GetType(w http.ResponseWriter, r *http.Request) {

}

func (h *pointHandler) AddType(w http.ResponseWriter, r *http.Request) {

}

func (h *pointHandler) GetPoint(w http.ResponseWriter, r *http.Request) {

}

func (h *pointHandler) AddPoint(w http.ResponseWriter, r *http.Request) {

}

func (h *pointHandler) GetByTypes(w http.ResponseWriter, r *http.Request) {

}
