package routes

import (
	"database/sql"
	"net/http"
)

type Server struct {
	Db *sql.DB
	Mux *http.ServeMux
}

type Todo struct {
	Id        int
	Task      string
	Completed bool
}

type TodoList struct {
	Id    int
	Title string
	Todos []Todo
}

type PageData struct {
	Title        string
	TodoLists    []TodoList
	ErrorMessage error
}