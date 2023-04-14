package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"moniso-server/internal/config"
	"moniso-server/internal/delivery/http_api"
	"moniso-server/internal/repository"
	"moniso-server/internal/services"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func Run() error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover")
			fmt.Println(r)
		}
	}()

	conf := config.New()

	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=moniso sslmode=disable")
	if err != nil {
		log.Println("Can't connect to database moniso")
		log.Println(err)
	} else {
		log.Printf("Connected to db: %s\n", "moniso")
	}
	defer db.Close()

	repositories := repository.NewRepositories(db)

	services := services.NewServices(repositories)

	apiServer := http_api.NewApiServer(conf.Api)
	err = apiServer.Start(services)
	if err != nil {
		log.Println(err)
	}

	return nil
}
