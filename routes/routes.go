package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func (s *Server) AddRoutes() {
	s.Mux.HandleFunc("/", s.index)
}

func (s *Server) index(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.go.tmpl"))
	todoLists, err := getAllTodoLists(s.Db)
	if err != nil {
		fmt.Fprint(res, err.Error())
		return
	}
	data := PageData {
		Title: "Todos",
		TodoLists: todoLists,
		ErrorMessage: err,
	}

	tmpl.Execute(res, data)
}

