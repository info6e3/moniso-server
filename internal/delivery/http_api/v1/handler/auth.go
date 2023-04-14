package handler

import (
	"encoding/json"
	"io"
	"moniso-server/internal/delivery/http_api/v1/handler/validation"
	"moniso-server/internal/domain"
	"moniso-server/internal/services"
	"net/http"
)

type authHandler struct {
	service services.User
}

func newAuthHandler(service services.User) *authHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод запроса должен быть POST", http.StatusMethodNotAllowed)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	var user domain.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
		return
	}

	err = validation.User(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.Add(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(&user)
	w.Write(jsonData)
}

func (h *authHandler) Authorization(w http.ResponseWriter, r *http.Request) {

}
