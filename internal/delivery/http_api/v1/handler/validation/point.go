package validation

import (
	"moniso-server/internal/domain"
	"time"
)

func Point(point *domain.Point) error {
	if point == nil {
		return domain.NewAppError("Ошибка валидации")
	}
	if point.Type == 0 {
		return domain.NewAppError("Ошибка валидации: не указан тип")
	}
	var defaultTime = time.Time{}
	if point.Date.String() == defaultTime.String() {
		return domain.NewAppError("Ошибка валидации: не указана дата")
	}

	return nil
}
