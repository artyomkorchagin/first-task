package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/artyomkorchagin/first-task/internal/config"
	orderpostgresql "github.com/artyomkorchagin/first-task/internal/repository/postgres/order"
	orderservice "github.com/artyomkorchagin/first-task/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//	@title			WB First Task
//	@version		1.0

//	@contact.name	Artyom Korchagin
//	@contact.email	artyomkorchagin333@gmail.com

//	@host		localhost:3000
//	@BasePath	/

func main() {
	db, err := sql.Open(config.GetDriver(), config.GetDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
	repo := orderpostgresql.NewRepository(db)
	service := orderservice.NewService(repo)
	fmt.Print(service)
}
