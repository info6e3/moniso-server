package services

import (
	"github.com/golang-jwt/jwt"
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
	"strconv"
	"time"
)

func generateToken(user *domain.User, duration time.Duration) (*domain.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(user.Id),
	})

	tokenString, err := token.SignedString([]byte("mykey"))
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		Owner: user.Id,
		Value: tokenString,
	}, nil
}

func validateTokenReturningPayload(tokenString string) (*domain.Payload, error) {
	claims := jwt.StandardClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewAppError("Неправильный метод авторизации")
		}

		return []byte("mykey"), nil
	})
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, err
	}

	return &domain.Payload{
		UserId: id,
	}, nil
}

func saveRefreshToken(repo repository.RefreshTokens, token *domain.Token) (int, error) {
	tokenId, err := repo.Add(token)
	if err != nil {
		return 0, err
	}

	return tokenId, nil
}
