package handler

import (
	"encoding/json"
	"io"
	"moniso-server/internal/delivery/http_api/v1/handler/validation"
	"moniso-server/internal/domain"
	"moniso-server/internal/services"
	"net/http"
)

type Registration struct {
	Service services.Auth
}

type Authorization struct {
	Service services.Auth
}

type Refresh struct {
	Service services.Auth
}

func (h Registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
				return
			}

			var reqUser domain.User
			err = json.Unmarshal(body, &reqUser)
			if err != nil {
				http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
				return
			}

			err = validation.User(&reqUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h.Service.Registration(&reqUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonData, _ := json.Marshal(&reqUser)
			w.Write(jsonData)
		}
	default:
		http.Error(w, "Необрабатываемый метод запросв", http.StatusMethodNotAllowed)
	}
}

func (h Authorization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
				return
			}

			var reqUser domain.User
			err = json.Unmarshal(body, &reqUser)
			if err != nil {
				http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
				return
			}

			err = validation.User(&reqUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			refreshToken, accessToken, err := h.Service.Authorization(&reqUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    refreshToken,
				MaxAge:   12 * 60 * 60,
				HttpOnly: true,
			})

			w.Write([]byte(accessToken))
		}
	default:
		http.Error(w, "Необрабатываемый метод запросв", http.StatusMethodNotAllowed)
	}
}

func (h Refresh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		{
			token, err := r.Cookie("token")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			accessToken, err := h.Service.Refresh(token.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Write([]byte(accessToken))
		}
	default:
		http.Error(w, "Необрабатываемый метод запроса", http.StatusMethodNotAllowed)
	}
}
