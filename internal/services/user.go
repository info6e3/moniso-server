package services

import (
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
)

type userService struct {
	repo repository.Users
}

func newUserService(repo repository.Users) *userService {
	return &userService{
		repo: repo,
	}
}

func addUser(repo repository.Users, user *domain.User) (int, error) {
	return repo.Add(user)
}

func getUserByLogin(repo repository.Users, login string) (*domain.User, error) {
	return repo.GetByLogin(login)
}
