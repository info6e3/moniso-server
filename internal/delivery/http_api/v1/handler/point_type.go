package handler

import (
	"encoding/json"
	"io"
	"log"
	"moniso-server/internal/delivery/http_api/v1/handler/validation"
	"moniso-server/internal/domain"
	"moniso-server/internal/services"
	"net/http"
	"net/url"
	"strconv"
)

type PointType struct {
	Service services.PointType
}

func (h PointType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			u, err := url.Parse(r.URL.String())
			if err != nil {
				http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)

			// Получаение типа по id
			id, err := strconv.Atoi(u.Query().Get("id"))
			if err == nil {
				pointType, err := h.Service.Get(id)
				if err != nil || payload.UserId != pointType.Owner {
					http.Error(w, "Такого типа нет или он принадлежит не вам", http.StatusBadRequest)
					return
				}
				jsonData, _ := json.Marshal(pointType)
				w.Write(jsonData)
				return
			}

			// Получаение типов по owner если параметр - all
			value := u.Query().Get("all")
			if value != "" {
				pointTypes, err := h.Service.GetAllByOwner(payload.UserId)
				if err != nil {
					http.Error(w, "Таких точек нет", http.StatusBadRequest)
				}
				jsonData, _ := json.Marshal(pointTypes)
				w.Write(jsonData)
				return
			}

			http.Error(w, "Необрабатываемые поля запроса", http.StatusBadRequest)
		}
	case "POST":
		{
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
				return
			}

			var reqPointType domain.PointType
			err = json.Unmarshal(body, &reqPointType)
			if err != nil {
				log.Println(err)
				http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)
			reqPointType.Owner = payload.UserId

			err = validation.PointType(&reqPointType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			pointType, err := h.Service.New(reqPointType.Type, reqPointType.Title, reqPointType.Owner, reqPointType.Min, reqPointType.Max)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = h.Service.Add(pointType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonData, _ := json.Marshal(pointType)
			w.Write(jsonData)
		}
	default:
		http.Error(w, "Необрабатываемый метод запросв", http.StatusMethodNotAllowed)
	}
}
