package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg/v10"

	"todo/domain"
	"todo/handlers"
	"todo/postgres"
)

func main() {
	DB := postgres.New(&pg.Options{
		User:     "postgres",
		Password: "1234",
		Database: "todo_dev",
	})

	defer DB.Close()

	domainDB := domain.DB{
		UserRepo: postgres.NewUserRepo(DB),
		TodoRepo: postgres.NewTodoRepo(DB),
	}

	d := &domain.Domain{DB: domainDB}

	r := handlers.SetupRouter(d)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	// ListenAndServe starts an HTTP server with a given address and handler
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("cannot start server %v", err)
	}
}
