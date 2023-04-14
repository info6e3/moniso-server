package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/pkg/params_query_string"
)

type PointTypesRepo struct {
	db *sqlx.DB
}

func NewPointTypesRepo(db *sqlx.DB) *PointTypesRepo {
	return &PointTypesRepo{
		db: db,
	}
}

func (r *PointTypesRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS point_types`)

}

func (r *PointTypesRepo) CreateTable() {
	schema := `
		CREATE TABLE point_types (
			id SERIAL PRIMARY KEY,
			type INTEGER NOT NULL,
			title TEXT NOT NULL,
			owner INTEGER REFERENCES users (id) NOT NULL,
			min INTEGER NOT NULL,
			max INTEGER NOT NULL
		)`
	r.db.MustExec(schema)
}

func getPointType(db *sqlx.DB, params map[string]any) (*domain.PointType, error) {
	if len(params) == 0 {
		return nil, domain.NewAppError("Нет параметров для поиска")
	}

	query, values := params_query_string.Generate("point_types", params)

	var pointType domain.PointType
	err := db.Get(&pointType, query, values...)
	if err != nil {
		return nil, err
	}

	return &pointType, nil
}

func (r *PointTypesRepo) Get(id int) (*domain.PointType, error) {
	return getPointType(r.db, map[string]any{
		"id": id,
	})
}

func (r *PointTypesRepo) Add(pointType *domain.PointType) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO point_types (type, title, owner, min, max)
	VALUES (:type, :title, :owner, :min, :max) RETURNING id`, &pointType)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.StructScan(pointType)
		if err != nil {
			return 0, err
		}
	}

	return pointType.Id, nil
}

func (r *PointTypesRepo) GetAllByOwner(userId int) ([]domain.PointType, error) {
	pointTypes := make([]domain.PointType, 0)
	err := r.db.Select(&pointTypes, `SELECT * FROM point_types WHERE owner = $1`, userId)
	if err != nil {
		return nil, err
	}
	if len(pointTypes) == 0 {
		return nil, domain.NewAppError("No results")
	}

	return pointTypes, nil
}

func (r *PointTypesRepo) Remove(id int) (int, error) {
	var pointType domain.PointType
	err := r.db.Get(&pointType, `DELETE FROM point_types WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}

	return pointType.Id, nil
}
