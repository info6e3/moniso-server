package handler

import (
	"encoding/json"
	"io"
	"moniso-server/internal/delivery/http_api/v1/handler/validation"
	"moniso-server/internal/domain"
	"moniso-server/internal/services"
	"net/http"
	"net/url"
	"strconv"
)

type Event struct {
	Service services.Event
}

func (h Event) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			u, err := url.Parse(r.URL.String())
			if err != nil {
				http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)

			id, err := strconv.Atoi(u.Query().Get("id"))
			if err != nil {
				http.Error(w, "Не удалось прочитать параметр id", http.StatusBadRequest)
				return
			}

			event, err := h.Service.Get(id)
			if err != nil || event.Owner != payload.UserId {
				http.Error(w, "События таким id нет или он принадлежит не вам", http.StatusBadRequest)
				return
			}

			jsonData, _ := json.Marshal(event)
			w.Write(jsonData)
		}
	case "POST":
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
				return
			}

			var reqEvent domain.Event
			err = json.Unmarshal(body, &reqEvent)
			if err != nil {
				http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)
			reqEvent.Owner = payload.UserId

			err = validation.Event(&reqEvent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			event, err := h.Service.New(reqEvent.Title, reqEvent.Owner, reqEvent.Date, reqEvent.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = h.Service.Add(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonData, _ := json.Marshal(event)
			w.Write(jsonData)
		}
	default:
		http.Error(w, "Необрабатываемый метод запросв", http.StatusMethodNotAllowed)
	}
}
