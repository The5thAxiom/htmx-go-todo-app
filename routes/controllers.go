package routes

import (
	"database/sql"
	"fmt"
	"log"
)

func getAllTodoLists(db *sql.DB) ([]TodoList, error) {
	var todoLists []TodoList

	todoListRows, todoListError := db.Query("SELECT id, title FROM TodoList;")
	if todoListError != nil {
		log.Printf("getAllTodoLists: todoListError: %s\n", todoListError)
		return nil, fmt.Errorf("getAllTodoLists: %s", todoListError)
	}
	defer todoListRows.Close()

	for todoListRows.Next() {
		var todoList TodoList
		if err := todoListRows.Scan(&todoList.Id, &todoList.Title); err != nil {
			log.Printf("getAllTodoLists: TodoListRowsScanErr: %s\n", err)
			return nil, fmt.Errorf("getAllTodoLists: scanning todoListRows %s", err)
		}

		todoRows, todoErr := db.Query(
			`SELECT id, task, ifnull(description, ''), completed
			FROM Todo
			WHERE todoListId = ?;`,
			todoList.Id,
		)
		if todoErr != nil {
			log.Printf("getAllTodoLists: todoErr: %s\n", todoErr)
			return nil, fmt.Errorf("getAllTodoLists: todoErr %s", todoErr)
		}
		defer todoRows.Close()

		for todoRows.Next() {
			var todo Todo
			if err := todoRows.Scan(&todo.Id, &todo.Task, &todo.Description, &todo.Completed); err != nil {
			log.Printf("getAllTodoLists: TodoRowsScanErr: %s\n", err)
				return nil, fmt.Errorf("getAllTodoLists: scanning todoRows: %s", err)
			}
			todoList.Todos = append(todoList.Todos, todo)
		}
		todoLists = append(todoLists, todoList)
	}
	return todoLists, nil
}