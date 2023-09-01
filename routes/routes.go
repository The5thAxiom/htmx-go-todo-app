package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (s *Server) AddRoutes() {
	s.Mux.HandleFunc("/", s.index)
}

func (s *Server) index(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.go.tmpl"))
	todoLists, todoListsError := getAllTodoLists(s.Db)
	if todoListsError != nil {
		log.Printf("index: todoListsError: %s\n", todoListsError)
		fmt.Fprint(res, todoListsError.Error())
		return
	}
	data := PageData {
		Title: "Todos",
		TodoLists: todoLists,
		ErrorMessage: todoListsError,
	}

	tmpl.Execute(res, data)
}

