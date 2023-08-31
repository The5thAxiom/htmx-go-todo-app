package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var tmpl *template.Template

var db *sql.DB

type Todo struct {
	Id int
	Task string
	Completed bool
}

type TodoList struct {
	Id int
	Title string
	Todos []Todo
}

type PageData struct {
	Title     string
	TodoLists []TodoList
	ErrorMessage error
}

func getAllTodoLists() ([]TodoList, error) {
	var todoLists []TodoList

	todoListRows, todoListError := db.Query("SELECT id, title FROM TodoList;")
	if todoListError != nil {
		return nil, fmt.Errorf("getAllTodoLists: %s", todoListError)
	}
	defer todoListRows.Close()

	for todoListRows.Next() {
		var todoList TodoList
		if err := todoListRows.Scan(&todoList.Id, &todoList.Title); err != nil {
			return nil, fmt.Errorf("getAllTodoLists: scanning todoListRows %s", err)
		}

		todoRows, todoErr := db.Query(
			`SELECT id, task, completed
			FROM Todo
			WHERE todoListId = ?;`,
			todoList.Id,
		)
		if todoErr != nil {
			return nil, fmt.Errorf("getAllTodoLists: todoErr %s", todoErr)
		}
		for todoRows.Next() {
			var todo Todo
			if err := todoRows.Scan(&todo.Id, &todo.Task, &todo.Completed); err != nil {
				return nil, fmt.Errorf("getAllTodoLists: scanning todoRows: %s", err)
			}
			todoList.Todos = append(todoList.Todos, todo)
		}
		todoLists = append(todoLists, todoList)
	}
	return todoLists, nil
}

func index(res http.ResponseWriter, req *http.Request) {
	todoLists, err := getAllTodoLists()
	data := PageData {
		Title: "Todos",
		TodoLists: todoLists,
		ErrorMessage: err,
	}

	tmpl.Execute(res, data)
}

func main() {
	tmpl = template.Must(template.ParseFiles("templates/index.gohtml"))
	fs := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", index)

	var dbName = "Todos.db"

	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connected to db: %s\n", dbName)
	defer db.Close()

	const PORT = 42069
	fmt.Printf("Starting a server on port %d\n", PORT)
	serverError := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	log.Fatal(serverError)
}
