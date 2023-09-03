package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) AddEndpoints() {
	s.Mux.HandleFunc("/api/todos", s.handleTodo)
	s.Mux.HandleFunc("/api/todoLists", s.handleTodoList)
}

func (s *Server) handleTodo (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		todoListId, todoListIdErr := strconv.Atoi(req.URL.Query().Get("todoListId"))
		if todoListIdErr != nil {
			log.Printf("addTodo: todoListIdErr: %s\n", todoListIdErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoListId: must be an integer")
			return
		}
	
		task := req.FormValue("task")
		description := req.FormValue("description")
	
		insert, insertErr := s.Db.Exec(`
			INSERT INTO Todo (
				task,
				description,
				completed,
				todoListId
			) VALUES (
				?,
				?,
				FALSE,
				?
			);`, task, description, todoListId,
		)
		if insertErr != nil {
			log.Printf("addTodo: insertErr: %s\n", insertErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoListId: %d\n", todoListId)
			return
		}
		todoId, _ := insert.LastInsertId()
		todo := Todo {
			Id: int(todoId),
			Task: task,
			Description: description,
			Completed: false,
		}
		
		tmpl := template.Must(template.ParseFiles("templates/todo.go.html"))
		tmplErr := tmpl.Execute(res, todo)
		if tmplErr != nil {
			log.Printf("handleTodo: POST: templateError: %s\n", tmplErr)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PATCH":
		todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("todoId"))
		if todoIdErr != nil {
			log.Printf("toggleTodo: todoIdErr: %s\n", todoIdErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoId: must be an integer")
			return
		}

		setCompleted, setCompletedErr := strconv.ParseBool(req.URL.Query().Get("setCompleted"))
		if setCompletedErr != nil {
			log.Printf("toggleTodo: setCompletedErr: %s\n", setCompletedErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid setCompleted: must be a boolean")
			return
		}
		
		var completed string
		if setCompleted {
			completed = "TRUE"
		} else {
			completed = "FALSE"
		}

		_, queryErr := s.Db.Exec("UPDATE Todo SET completed=? WHERE id=?;", completed, todoId)
		if queryErr != nil {
			log.Printf("toggleTodo: queryErr: %s\n", queryErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoId: %d\n", todoId)
			return
		}
		
		todoRow := s.Db.QueryRow("SELECT id, task, description, completed FROM Todo WHERE id=?;", todoId)
		var todo Todo
		todoRow.Scan(&todo.Id, &todo.Task, &todo.Description, &todo.Completed)

		tmpl := template.Must(template.ParseFiles("templates/todo.go.html"))
		tmplErr := tmpl.Execute(res, todo)
		if tmplErr != nil {
			log.Printf("handleTodo: PATCH: templateError: %s\n", tmplErr)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "DELETE":
		todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("todoId"))
		if todoIdErr != nil {
			log.Printf("handleTodo: DELETE: todoIdErr: %s\n", todoIdErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoId: must be an integer")
			return
		}

		_, queryErr := s.Db.Exec("DELETE FROM Todo WHERE id=?;", todoId)
		if queryErr != nil {
			log.Printf("handleTodo: DELETE: queryErr: %s\n", queryErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoId: %d\n", todoId)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "handleTodo: %s: Invalid http method\n", req.Method)
	}
}

func (s *Server) handleTodoList (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		title := req.FormValue("title")
		insert, insertErr := s.Db.Exec(`INSERT INTO TodoList (title) VALUES (?);`, title)
		if insertErr != nil {
			log.Printf("handleTodoList: POST: %s\n", insertErr)
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, "A database insertion error occurred")
			return
		}

		todoListId, _ := insert.LastInsertId()
		todoList := TodoList {
			Id: int(todoListId),
			Title: title,
			Todos: []Todo {},
		}

		tmpl := template.Must(template.ParseFiles("templates/todoList.go.html", "templates/todo.go.html"))
		tmplErr := tmpl.ExecuteTemplate(res, "todoList.go.html", todoList)
		if tmplErr != nil {
			log.Printf("handleTodoList: POST: templateError: %s\n", tmplErr)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "DELETE":
		todoListId, todoListIdErr := strconv.Atoi(req.URL.Query().Get("todoListId"))
		if todoListIdErr != nil {
			log.Printf("handleTodoList: DELETE: todoListIdErr: %s\n", todoListIdErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoId: must be an integer")
			return
		}

		_, queryErr := s.Db.Exec("DELETE FROM TodoList WHERE id=?;", todoListId)
		if queryErr != nil {
			log.Printf("handleTodoList: DELETE: queryErr: %s\n", queryErr)
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoListId: %d\n", todoListId)
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "handleTodoList: %s: Invalid http method\n", req.Method)
	}
}
