package validation

import (
	"moniso-server/internal/domain"
	"time"
)

func Event(event *domain.Event) error {
	if event == nil {
		return domain.NewAppError("Ошибка валидации")
	}
	if event.Title == "" {
		return domain.NewAppError("Ошибка валидации: не указано название")
	}
	if event.Owner == 0 {
		return domain.NewAppError("Ошибка валидации: не указан владелец")
	}
	var defaultTime = time.Time{}
	if event.Date.String() == defaultTime.String() {
		return domain.NewAppError("Ошибка валидации: не указана дата")
	}

	return nil
}
