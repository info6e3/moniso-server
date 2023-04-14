package services

import (
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
	"time"
)

type pointService struct {
	repoPointTypes repository.PointTypes
	repoPoints     repository.Points
}

func newPointService(repoPointTypes repository.PointTypes, repoPoints repository.Points) *pointService {
	return &pointService{
		repoPointTypes: repoPointTypes,
		repoPoints:     repoPoints,
	}
}

func (s *pointService) New(ownerId int, typeId int, value int, date time.Time, description string) (*domain.Point, error) {
	// Проверка существования типа
	pointType, err := s.repoPointTypes.Get(typeId)
	if err != nil {
		return nil, err
	}

	// Проверка принадлежности типа
	if pointType.Owner != ownerId {
		return nil, domain.NewAppError("Тип точки принадлежит не вам")
	}

	if value < pointType.Min {
		return nil, domain.NewAppError("Значение точки меньше допустимого")
	}

	if value > pointType.Max {
		return nil, domain.NewAppError("Значение точки больше допустимого")
	}

	return &domain.Point{
		Type:        typeId,
		Description: description,
		Value:       value,
		Date:        date,
	}, nil
}

func (s *pointService) Get(id int) (*domain.Point, error) {
	return s.repoPoints.Get(id)
}

func (s *pointService) Add(point *domain.Point) (int, error) {
	// Замещение точки с такой же датой и типом
	similarPoint, err := s.repoPoints.GetByTypeDate(point.Type, point.Date)
	if err == nil {
		return s.repoPoints.Update(similarPoint.Id, point)
	}

	return s.repoPoints.Add(point)
}

func (s *pointService) GetAllByType(pointTypeId int) ([]domain.Point, error) {
	return s.repoPoints.GetAllByType(pointTypeId)
}
