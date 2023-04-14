package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/pkg/params_query_string"
)

type RefreshTokensRepo struct {
	db *sqlx.DB
}

func NewRefreshTokensRepo(db *sqlx.DB) *RefreshTokensRepo {
	return &RefreshTokensRepo{
		db: db,
	}
}

func (r *RefreshTokensRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS refresh_tokens`)
}

func (r *RefreshTokensRepo) CreateTable() {
	schema := `
		CREATE TABLE refresh_tokens (
			id SERIAL PRIMARY KEY,
			owner INTEGER REFERENCES users (id) NOT NULL,
			value TEXT NOT NULL
		)`
	r.db.MustExec(schema)
}

func getRefreshToken(db *sqlx.DB, params map[string]any) (*domain.Token, error) {
	if len(params) == 0 {
		return nil, domain.NewAppError("Нет параметров для поиска")
	}

	query, values := params_query_string.Generate("refresh_tokens", params)

	var token domain.Token
	err := db.Get(&token, query, values...)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokensRepo) Get(id int) (*domain.Token, error) {
	return getRefreshToken(r.db, map[string]any{
		"id": id,
	})
}

func (r *RefreshTokensRepo) Exists(value string) error {
	_, err := getRefreshToken(r.db, map[string]any{
		"value": value,
	})
	return err
}

func (r *RefreshTokensRepo) Add(token *domain.Token) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO refresh_tokens (owner, value)
	VALUES  (:owner, :value) RETURNING id`, &token)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.StructScan(&token)
		if err != nil {
			return 0, err
		}
	}

	return token.Id, nil
}
