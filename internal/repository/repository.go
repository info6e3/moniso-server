package repository

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/internal/repository/postgres"
	"time"
)

type Database interface {
	Create()
	Remove()
}

type Users interface {
	CreateTable()
	RemoveTable()
	Get(id int) (*domain.User, error)
	Add(user *domain.User) (int, error)

	GetByLogin(login string) (*domain.User, error)
}

type RefreshTokens interface {
	CreateTable()
	RemoveTable()
	Get(id int) (*domain.Token, error)
	Add(token *domain.Token) (int, error)
	Exists(tokenValue string) error
}

type PointTypes interface {
	CreateTable()
	RemoveTable()
	Get(id int) (*domain.PointType, error)
	Add(pointType *domain.PointType) (int, error)
	Remove(id int) (int, error)

	GetAllByOwner(userId int) ([]domain.PointType, error)
}

type Points interface {
	CreateTable()
	RemoveTable()
	Get(id int) (*domain.Point, error)
	Add(point *domain.Point) (int, error)
	Remove(id int) (int, error)
	Update(id int, point *domain.Point) (int, error)

	GetByTypeDate(pointTypeId int, date time.Time) (*domain.Point, error)
	GetAllByType(pointTypeId int) ([]domain.Point, error)
}

type Events interface {
	CreateTable()
	RemoveTable()
	Get(id int) (*domain.Event, error)
	Add(point *domain.Event) (int, error)
	Remove(id int) (int, error)
}

type Repositories struct {
	Database      Database
	Users         Users
	RefreshTokens RefreshTokens
	PointTypes    PointTypes
	Points        Points
	Events        Events
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Database:      postgres.NewDatabase(db),
		Users:         postgres.NewUsersRepo(db),
		RefreshTokens: postgres.NewRefreshTokensRepo(db),
		PointTypes:    postgres.NewPointTypesRepo(db),
		Points:        postgres.NewPointsRepo(db),
		Events:        postgres.NewEventsRepo(db),
	}
}

func (r *Repositories) CreateDatabase() {
	r.Database.Create()
}

func (r *Repositories) CreateAllRepos() {
	r.Users.CreateTable()
	r.RefreshTokens.CreateTable()
	r.PointTypes.CreateTable()
	r.Points.CreateTable()
	r.Events.CreateTable()
}

func (r *Repositories) RemoveAllRepos() {
	r.Events.RemoveTable()
	r.Points.RemoveTable()
	r.PointTypes.RemoveTable()
	r.RefreshTokens.RemoveTable()
	r.Users.RemoveTable()
}
