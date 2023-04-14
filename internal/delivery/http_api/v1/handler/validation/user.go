package validation

import (
	"moniso-server/internal/domain"
)

func User(user *domain.User) error {
	if user == nil {
		return domain.NewAppError("Ошибка валидации")
	}
	if user.Login == "" {
		return domain.NewAppError("Ошибка валидации: не указан логин")
	}
	if user.Password == "" {
		return domain.NewAppError("Ошибка валидации: не указан пароль")
	}
	if user.Username == "" {
		return domain.NewAppError("Ошибка валидации: не указан юзернейм")
	}

	return nil
}
