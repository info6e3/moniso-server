package services

import (
	"golang.org/x/crypto/bcrypt"
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
	"time"
)

type authService struct {
	repoUsers  repository.Users
	repoTokens repository.RefreshTokens
}

func newAuthService(repoUsers repository.Users, repoTokens repository.RefreshTokens) *authService {
	return &authService{
		repoUsers:  repoUsers,
		repoTokens: repoTokens,
	}
}

func (s *authService) Registration(user *domain.User) error {
	_, err := getUserByLogin(s.repoUsers, user.Login)
	if err == nil {
		return domain.NewAppError("Такой логин уже зарегестрирован")
	}

	hashBytePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashPassword := string(hashBytePassword)
	user.Password = hashPassword
	_, err = addUser(s.repoUsers, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Authorization(user *domain.User) (string, string, error) {
	dbUser, err := getUserByLogin(s.repoUsers, user.Login)
	if err != nil {
		return "", "", domain.NewAppError("Такого пользователя нет")
	}
	user.Id = dbUser.Id

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", "", domain.NewAppError("Неверный пароль")
	}

	refreshToken, err := generateToken(user, 12*time.Hour)
	if err != nil {
		return "", "", err
	}

	_, err = saveRefreshToken(s.repoTokens, refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := generateToken(user, time.Hour)
	if err != nil {
		return "", "", err
	}

	return refreshToken.Value, accessToken.Value, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	// Проверка существования токена в бд
	err := s.repoTokens.Exists(refreshToken)
	if err != nil {
		return "", err
	}

	payload, err := validateTokenReturningPayload(refreshToken)
	if err != nil {
		return "", err
	}

	// Проверка существования юзера из токена в бд
	user, err := s.repoUsers.Get(payload.UserId)
	if err != nil {
		return "", err
	}

	accessToken, err := generateToken(user, time.Hour)
	if err != nil {
		return "", err
	}

	return accessToken.Value, nil
}

func (s *authService) ValidateTokenReturningPayload(token string) (*domain.Payload, error) {
	return validateTokenReturningPayload(token)
}
