package validation

import (
	"moniso-server/internal/domain"
)

func PointType(pointType *domain.PointType) error {
	if pointType == nil {
		return domain.NewAppError("Ошибка валидации")
	}
	if !pointType.Type.Valid() {
		return domain.NewAppError("Ошибка валидации: неверный тип")
	}
	if pointType.Title == "" {
		return domain.NewAppError("Ошибка валидации: не указано название")
	}
	if pointType.Owner == 0 {
		return domain.NewAppError("Ошибка валидации: не указан владелец")
	}
	if pointType.Min >= pointType.Max {
		return domain.NewAppError("Ошибка валидации: неправильный диапазон значений")
	}

	return nil
}
