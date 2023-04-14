package services

import (
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
)

type pointTypeService struct {
	repoUsers      repository.Users
	repoPointTypes repository.PointTypes
}

func newPointTypeService(repoUsers repository.Users, repoPointTypes repository.PointTypes) *pointTypeService {
	return &pointTypeService{
		repoUsers:      repoUsers,
		repoPointTypes: repoPointTypes,
	}
}

func (s *pointTypeService) New(pType domain.PType, title string, ownerId int, min int, max int) (*domain.PointType, error) {
	// Проверка существования пользователя
	_, err := s.repoUsers.Get(ownerId)
	if err != nil {
		return nil, err
	}

	// Проверка типо по перечислению возможных типов
	if !pType.Valid() {
		return nil, domain.NewAppError("Неизвестный тип")
	}

	// Проверка диапазона типа
	if pType == domain.Number {
		if min >= max {
			return nil, domain.NewAppError("Максимальное значение больше минимального или равно ему")
		}
	} else {
		min = 0
		max = 1
	}

	return &domain.PointType{
		Type:  pType,
		Title: title,
		Owner: ownerId,
		Min:   min,
		Max:   max,
	}, nil
}

func (s *pointTypeService) Get(id int) (*domain.PointType, error) {
	return s.repoPointTypes.Get(id)
}

func (s *pointTypeService) Add(pointType *domain.PointType) (int, error) {
	return s.repoPointTypes.Add(pointType)
}

func (s *pointTypeService) GetAllByOwner(userId int) ([]domain.PointType, error) {
	return s.repoPointTypes.GetAllByOwner(userId)
}
