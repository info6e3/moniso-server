package postgres

import (
	"github.com/jmoiron/sqlx"
	"moniso-server/internal/domain"
	"moniso-server/internal/repository/postgres/db_domain"
)

type RefreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepo(db *sqlx.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		db: db,
	}
}

func (r *RefreshTokenRepo) RemoveTable() {
	r.db.MustExec(`DROP TABLE IF EXISTS refresh_tokens`)
}

func (r *RefreshTokenRepo) CreateTable() {
	schema := `
		CREATE TABLE refresh_tokens (
			id SERIAL PRIMARY KEY,
			owner INTEGER REFERENCES users (id) NOT NULL,
			token TEXT NOT NULL
		)`
	r.db.MustExec(schema)
}

func (r *RefreshTokenRepo) GetByOwner(ownerId int) (*domain.RefreshToken, error) {
	var dbRefreshToken db_domain.DbRefreshToken
	err := r.db.Get(&dbRefreshToken, `SELECT * FROM refresh_tokens WHERE owner = $1`, ownerId)
	if err != nil {
		return nil, err
	}

	owner, err := getUser(r.db, "owner", dbRefreshToken.Owner)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshToken{
		Id:    dbRefreshToken.Id,
		Owner: *owner,
		Token: dbRefreshToken.Token,
	}, nil
}

func (r *RefreshTokenRepo) Add(token *domain.RefreshToken) (int, error) {
	rows, err := r.db.NamedQuery(`INSERT INTO refresh_tokens (owner, token)
	VALUES (:owner, :token) RETURNING id`, &token)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err := rows.StructScan(token)
		if err != nil {
			return 0, err
		}
	}

	return token.Id, nil
}
