package services

import (
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
	"time"
)

type Auth interface {
	Registration(user *domain.User) error
	Authorization(user *domain.User) (refreshToken string, accessToken string, err error)
	Refresh(refreshToken string) (accessToken string, err error)
	ValidateTokenReturningPayload(token string) (*domain.Payload, error)
}

type PointType interface {
	New(pType domain.PType, title string, ownerId int, min int, max int) (*domain.PointType, error)
	Get(id int) (*domain.PointType, error)
	Add(pointType *domain.PointType) (int, error)
	GetAllByOwner(userId int) ([]domain.PointType, error)
}

type Point interface {
	New(ownerId int, typeId int, value int, date time.Time, description string) (*domain.Point, error)
	Get(int) (point *domain.Point, err error)
	Add(point *domain.Point) (int, error)
	GetAllByType(pointTypeId int) ([]domain.Point, error)
}

type Event interface {
	New(title string, ownerId int, date time.Time, description string) (*domain.Event, error)
	Get(id int) (*domain.Event, error)
	Add(point *domain.Event) (int, error)
}

type Services struct {
	Auth      Auth
	PointType PointType
	Point     Point
	Event     Event
}

func NewServices(repositories *repository.Repositories) *Services {
	return &Services{
		Auth:      newAuthService(repositories.Users, repositories.RefreshTokens),
		PointType: newPointTypeService(repositories.Users, repositories.PointTypes),
		Point:     newPointService(repositories.PointTypes, repositories.Points),
		Event:     newEventService(repositories.Users, repositories.Events),
	}
}
