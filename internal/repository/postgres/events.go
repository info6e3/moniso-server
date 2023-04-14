package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/pkg/params_query_string"
)

type EventsRepo struct {
	db *sqlx.DB
}

func NewEventsRepo(db *sqlx.DB) *EventsRepo {
	return &EventsRepo{
		db: db,
	}
}

func (r *EventsRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS events`)
}

func (r *EventsRepo) CreateTable() {
	schema := `
		CREATE TABLE events (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			owner INTEGER REFERENCES users (id) NOT NULL,
			description TEXT,
			date DATE NOT NULL
		)`
	r.db.MustExec(schema)
}

func getEvent(db *sqlx.DB, params map[string]any) (*domain.Event, error) {
	if len(params) == 0 {
		return nil, domain.NewAppError("Нет параметров для поиска")
	}

	query, values := params_query_string.Generate("events", params)

	var event domain.Event
	err := db.Get(&event, query, values...)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventsRepo) Get(id int) (*domain.Event, error) {
	return getEvent(r.db, map[string]any{
		"id": id,
	})
}

func (r *EventsRepo) Add(event *domain.Event) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO events (title, owner, description, date)
	VALUES (:title, :owner, :description, :date) RETURNING id`, &event)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.StructScan(event)
		if err != nil {
			return 0, err
		}
	}

	return event.Id, nil
}

func (r *EventsRepo) Remove(id int) (int, error) {
	var event domain.Event
	err := r.db.Get(&event, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return event.Id, nil
}
