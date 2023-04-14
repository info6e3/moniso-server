package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/pkg/params_query_string"
	"time"
)

type PointsRepo struct {
	db *sqlx.DB
}

func NewPointsRepo(db *sqlx.DB) *PointsRepo {
	return &PointsRepo{
		db: db,
	}
}

func (r *PointsRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS points`)
}

func (r *PointsRepo) CreateTable() {
	schema := `
		CREATE TABLE points (
			id SERIAL PRIMARY KEY,
			type  INTEGER REFERENCES point_types (id) NOT NULL,
			value INTEGER NOT NULL,
			description TEXT,
			date DATE NOT NULL
		)`
	r.db.MustExec(schema)
}

func getPoint(db *sqlx.DB, params map[string]any) (*domain.Point, error) {
	if len(params) == 0 {
		return nil, domain.NewAppError("Нет параметров для поиска")
	}

	query, values := params_query_string.Generate("points", params)

	var point domain.Point
	err := db.Get(&point, query, values...)
	if err != nil {
		return nil, err
	}

	return &point, nil
}

func (r *PointsRepo) Get(id int) (*domain.Point, error) {
	return getPoint(r.db, map[string]any{
		"id": id,
	})
}

func (r *PointsRepo) GetByTypeDate(pointTypeId int, date time.Time) (*domain.Point, error) {
	return getPoint(r.db, map[string]any{
		"type": pointTypeId,
		"date": date,
	})
}

func (r *PointsRepo) Add(point *domain.Point) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO points (type, value, description, date)
	VALUES (:type, :value, :description, :date) RETURNING id`, point)
	if err != nil {
		return 0, err
	}
	// Взять айди из ответа и добавить в точку
	for rows.Next() {
		err := rows.StructScan(point)
		if err != nil {
			return 0, err
		}
	}

	return point.Id, nil
}

func (r *PointsRepo) GetAllByType(pointTypeId int) ([]domain.Point, error) {
	points := make([]domain.Point, 0)
	err := r.db.Select(&points, `SELECT * FROM points WHERE type = $1`, pointTypeId)
	if err != nil {
		return nil, err
	}
	if len(points) == 0 {
		return nil, domain.NewAppError("No results")
	}

	return points, nil
}

func (r *PointsRepo) Remove(id int) (int, error) {
	var point domain.Point
	err := r.db.Get(&point, `DELETE FROM points WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return point.Id, nil
}

func (r *PointsRepo) Update(id int, point *domain.Point) (int, error) {
	point.Id = id
	_, err := r.db.NamedQuery(`UPDATE points SET type=:type, value=:value, description=:description, date=:date WHERE id=:id`, point)
	if err != nil {
		return 0, err
	}

	return point.Id, nil
}
