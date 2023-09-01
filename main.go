package main

import (
	"database/sql"
	"fmt"
	"htmx-go-todo-app/routes"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "Todos.db"
const port = 42069

func main() {
	db, dbErr := sql.Open("sqlite3", dbName)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connected to db: %s\n", dbName)
	defer db.Close()

	server := routes.Server {
		Db: db,
		Mux: http.NewServeMux(),
	}

	fs := http.FileServer(http.Dir("./static"))
	server.Mux.Handle("/static/", http.StripPrefix("/static/", fs))
	server.AddRoutes()
	server.AddEndpoints()

	fmt.Printf("Starting a server on port %d\n", port)
	serverError := http.ListenAndServe(fmt.Sprintf(":%d", port), server.Mux)
	log.Fatal(serverError)
}
