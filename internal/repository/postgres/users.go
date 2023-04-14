package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/pkg/params_query_string"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS users`)
}

func (r *UsersRepo) CreateTable() {
	schema := `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			login TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			username TEXT NOT NULL
		)`
	r.db.MustExec(schema)
}

func getUser(db *sqlx.DB, params map[string]any) (*domain.User, error) {
	if len(params) == 0 {
		return nil, domain.NewAppError("Нет параметров для поиска")
	}

	query, values := params_query_string.Generate("users", params)

	var user domain.User
	err := db.Get(&user, query, values...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepo) Get(id int) (*domain.User, error) {
	return getUser(r.db, map[string]any{
		"id": id,
	})
}

func (r *UsersRepo) Add(user *domain.User) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO users (login, password, username)
	VALUES (:login, :password, :username) RETURNING id`, &user)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			return 0, err
		}
	}
	return user.Id, nil
}

func (r *UsersRepo) GetByLogin(login string) (*domain.User, error) {
	return getUser(r.db, map[string]any{
		"login": login,
	})
}
