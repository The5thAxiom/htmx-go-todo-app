CREATE TABLE TodoList (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL
);

CREATE TABLE Todo (
    id INTEGER PRIMARY KEY,
    task TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL,
    todoListId INTEGER,
    FOREIGN KEY (todoListId) REFERENCES TodoList(id) ON DELETE CASCADE
);