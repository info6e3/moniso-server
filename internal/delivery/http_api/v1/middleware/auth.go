package middleware

import (
	"context"
	"moniso-server/internal/services"
	"net/http"
	"strings"
)

type Auth struct {
	Service services.Auth
	Next    http.Handler
}

func (m Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Вы не авторизированы", http.StatusUnauthorized)
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		http.Error(w, "Неверный хедер авторизации", http.StatusUnauthorized)
		return
	}

	accessToken := authHeaderParts[1]
	payload, err := m.Service.ValidateTokenReturningPayload(accessToken)
	if err != nil {
		http.Error(w, "Неверный токен авторизации", http.StatusUnauthorized)
		return
	}

	// Создаем новый контекст и добавляем значение
	ctx := context.WithValue(r.Context(), "payload", payload)

	// Передаем обновленный контекст следующему middleware
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}
