package routes

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *Server) AddEndpoints() {
	s.Mux.HandleFunc("/api/todos", s.handleTodo)
	s.Mux.HandleFunc("/api/todoLists", s.handleTodoList)
	s.Mux.HandleFunc("/api/todos/complete", s.completeTodo)
	s.Mux.HandleFunc("/api/todos/uncomplete", s.uncompleteTodo)
	s.Mux.HandleFunc("/api/todos/add", s.addTodo)
	s.Mux.HandleFunc("/api/todoLists/add", s.addTodoList)
}

func (s *Server) handleTodo (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "DELETE":
		todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("todoId"))
		if todoIdErr != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoId: must be an integer")
			return
		}

		_, queryErr := s.Db.Exec("DELETE FROM Todo WHERE id=?;", todoId)
		if queryErr != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoId: %d\n", todoId)
			return
		}

		res.WriteHeader(http.StatusAccepted)
	default:
		res.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) handleTodoList (res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "DELETE":
		todoListId, todoListIdErr := strconv.Atoi(req.URL.Query().Get("todoListId"))
		if todoListIdErr != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(res, "Invalid todoId: must be an integer")
			return
		}

		_, queryErr := s.Db.Exec("DELETE FROM TodoList WHERE id=?;", todoListId)
		if queryErr != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Invalid todoListId: %d\n", todoListId)
			return
		}

		res.WriteHeader(http.StatusAccepted)
	default:
		res.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) completeTodo(res http.ResponseWriter, req *http.Request) {
	todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("todoId"))
	if todoIdErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, "Invalid todoId: must be an integer")
		return
	}
	_, queryErr := s.Db.Exec("UPDATE Todo SET completed=TRUE WHERE id=?;", todoId)
	if queryErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid todoId: %d\n", todoId)
		return
	}

	fmt.Fprintf(
		res,
		`<button
			hx-get="/api/todos/uncomplete?todoId=%d"
			hx-trigger="click"
			hx-swap="outerHTML"
		>X</button>`,
		todoId,
	)
}

func (s *Server) uncompleteTodo(res http.ResponseWriter, req *http.Request) {
	todoId, todoIdErr := strconv.Atoi(req.URL.Query().Get("todoId"))
	if todoIdErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, "Invalid todoId: must be an integer")
		return
	}
	_, queryErr := s.Db.Exec("UPDATE Todo SET completed=FALSE WHERE id=?;", todoId)
	if queryErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid todoId: %d\n", todoId)
		return
	}

	fmt.Fprintf(
		res,
		`<button
			hx-get="/api/todos/complete?todoId=%d"
			hx-trigger="click"
			hx-swap="outerHTML"
		>&nbsp;</button>`,
		todoId,
	)
}

func (s *Server) addTodo(res http.ResponseWriter, req *http.Request) {
	todoListId, todoListIderr := strconv.Atoi(req.URL.Query().Get("todoListId"))
	if todoListIderr != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, "Invalid todoListId: must be an integer")
		return
	}

	task := req.FormValue("task")
	description := req.FormValue("description")

	insert, insertError := s.Db.Exec(`
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
	if insertError != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "Invalid todoListId: %d\n", todoListId)
		return
	}
	todoId, _ := insert.LastInsertId()

	fmt.Fprintf(
		res,
		`<li>
		<div class="todo-task">
			<button
				hx-get="/api/todos/complete?todoId=%d"
				hx-trigger="click"
				hx-swap="outerHTML"
			>&nbsp;</button>%s
			<button
				hx-delete="/api/todos?todoId=%d"
				hx-swap="delete"
				hx-trigger="click"
				hx-target="closest li"
				hx-confirm="Are you sure you want to delete this task?"
			>Delete</button>
		</div>
		<div class="todo-description">%s</div>
	</li>`,
		todoId, task, todoId, description,
	)
}

func (s *Server) addTodoList (res http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	insert, insertErr := s.Db.Exec(`INSERT INTO TodoList (title) VALUES (?);`, title)
	if insertErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, "A database insertion error occurred")
		return
	}

	todoListId, _ := insert.LastInsertId()

	fmt.Fprintf(res,`
		<li>
			%s
			<button
				hx-delete="/api/todoLists?todoListId=%d"
				hx-swap="delete"
				hx-trigger="click"
				hx-target="closest li"
				hx-confirm="Are you sure you want to delete this todoList? This will also delete all tasks in the list."
			>Delete</button>
			<ul id="todoList-%d">
			</ul>
			<form
				class="new-todo-form"
				hx-post="/api/todos/add?todoListId=%d"
				hx-swap="beforeend"
				hx-target="#todoList-%d"
				hx-on::after-request="this.reset()"
			>
				<input type="text" placeholder="Task" name="task"/>
				<textarea placeholder="description" name="description"></textarea>
				<button>Add Todo to list <i>%s</i></button>
			</form>
		</li>`, title, todoListId, todoListId, todoListId, todoListId, title,
	)
}