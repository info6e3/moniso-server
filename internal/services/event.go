package services

import (
	"moniso-server/internal/domain"
	"moniso-server/internal/repository"
	"time"
)

type eventService struct {
	repoUsers  repository.Users
	repoEvents repository.Events
}

func newEventService(repoUsers repository.Users, repoEvents repository.Events) *eventService {
	return &eventService{
		repoUsers:  repoUsers,
		repoEvents: repoEvents,
	}
}

func (s *eventService) New(title string, ownerId int, date time.Time, description string) (*domain.Event, error) {
	// Проверка существования пользователя
	_, err := s.repoUsers.Get(ownerId)
	if err != nil {
		return nil, err
	}

	return &domain.Event{
		Title:       title,
		Owner:       ownerId,
		Description: description,
		Date:        date,
	}, nil
}

func (s *eventService) Get(id int) (*domain.Event, error) {
	return s.repoEvents.Get(id)
}

func (s *eventService) Add(pointType *domain.Event) (int, error) {
	return s.repoEvents.Add(pointType)
}
