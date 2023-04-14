package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"moniso-server/internal/repository"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=moniso sslmode=disable")
	if err != nil {
		log.Println("Can't connect to database moniso")
		log.Println(err)
	} else {
		log.Printf("Connected to db: %s\n", "moniso")
	}
	defer db.Close()
	repositories := repository.NewRepositories(db)
	repositories.RemoveAllRepos()
	repositories.CreateAllRepos()
}
