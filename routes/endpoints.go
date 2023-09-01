package routes

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) AddEndpoints() {
	s.Mux.HandleFunc("/api/todos/complete", s.completeTodo)
	s.Mux.HandleFunc("/api/todos/uncomplete", s.uncompleteTodo)
}

func (s *Server) completeTodo(res http.ResponseWriter, req *http.Request) {
	todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("id"))
	if todoIdErr != nil {
		res.WriteHeader((http.StatusBadRequest))
		fmt.Fprintln(res, "Invalid Todo id: must be an integer")
		return
	}
	_, queryErr := s.Db.Exec("UPDATE Todo SET completed=TRUE WHERE id=?;", todoId)
	if queryErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid Todo id: %d\n", todoId)
		return
	}

	fmt.Fprintf(
		res,
		`<button
			hx-get="/api/todos/uncomplete?id=%d"
			hx-trigger="click"
			hx-swap="outerHTML"
		>X</button>`,
		todoId,
	)
}

func (s *Server) uncompleteTodo(res http.ResponseWriter, req *http.Request) {
	todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("id"))
	if todoIdErr != nil {
		res.WriteHeader((http.StatusBadRequest))
		fmt.Fprintln(res, "Invalid Todo id: must be an integer")
		return
	}
	_, queryErr := s.Db.Exec("UPDATE Todo SET completed=FALSE WHERE id=?;", todoId)
	if queryErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid Todo id: %d\n", todoId)
		return
	}

	fmt.Fprintf(
		res,
		`<button
			hx-get="/api/todos/complete?id=%d"
			hx-trigger="click"
			hx-swap="outerHTML"
		>&nbsp;</button>`,
		todoId,
	)

}
