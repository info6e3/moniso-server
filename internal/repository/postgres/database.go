package postgres

import "github.com/jmoiron/sqlx"

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) Create() {
	d.db.MustExec(`CREATE DATABASE moniso`)
}

func (d *Database) Remove() {
	d.db.MustExec(`DROP DATABASE moniso`)
}
