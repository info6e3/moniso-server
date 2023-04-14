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

type Point struct {
	ServicePoint     services.Point
	ServicePointType services.PointType
}

func (h Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			u, err := url.Parse(r.URL.String())
			if err != nil {
				http.Error(w, "Не удалось прочитать запрос", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)

			// Получаение точки по id
			id, err := strconv.Atoi(u.Query().Get("id"))
			if err == nil {
				point, err := h.ServicePoint.Get(id)
				if err != nil {
					http.Error(w, "Точки с таким id нет", http.StatusBadRequest)
					return
				}
				pointType, err := h.ServicePointType.Get(point.Type)
				if err != nil || pointType.Owner != payload.UserId {
					http.Error(w, "Типа с таким id нет или он принадлежит не вам", http.StatusBadRequest)
					return
				}
				jsonData, _ := json.Marshal(point)
				w.Write(jsonData)
				return
			}

			// Получаение точек по типу
			pointTypeId, err := strconv.Atoi(u.Query().Get("type"))
			if err == nil {
				pointType, err := h.ServicePointType.Get(pointTypeId)
				if err != nil || pointType.Owner != payload.UserId {
					http.Error(w, "Типа с таким id нет или он принадлежит не вам", http.StatusBadRequest)
					return
				}

				points, err := h.ServicePoint.GetAllByType(pointTypeId)
				if err != nil {
					http.Error(w, "Точек с таким типом нет", http.StatusBadRequest)
					return
				}

				jsonData, _ := json.Marshal(points)
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

			var reqPoint domain.Point
			err = json.Unmarshal(body, &reqPoint)
			if err != nil {
				http.Error(w, "Не удалось декодировать JSON тело запроса", http.StatusBadRequest)
				return
			}

			payload := r.Context().Value("payload").(*domain.Payload)

			err = validation.Point(&reqPoint)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			
			point, err := h.ServicePoint.New(payload.UserId, reqPoint.Type, reqPoint.Value, reqPoint.Date, reqPoint.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = h.ServicePoint.Add(point)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			jsonData, _ := json.Marshal(point)
			w.Write(jsonData)
		}
	default:
		http.Error(w, "Необрабатываемый метод запросв", http.StatusMethodNotAllowed)
	}
}
