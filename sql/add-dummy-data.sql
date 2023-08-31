INSERT INTO TodoList (title) VALUES ('Hostel Room');
INSERT INTO TodoList (title) VALUES ('Studies');
INSERT INTO TodoList (title) VALUES ('Projects');

INSERT INTO Todo (
    task,
    description,
    completed,
    todoListId
) VALUES (
    'Clean bed',
    'put clean bedsheets',
    FALSE,
    1
);

INSERT INTO Todo (
    task,
    description,
    completed,
    todoListId
) VALUES (
    'Buy facewash',
    NULL,
    TRUE,
    1
);

INSERT INTO Todo (
    task,
    description,
    completed,
    todoListId
) VALUES (
    'Complete AI DA-1',
    NULL,
    TRUE,
    2
);

INSERT INTO Todo (
    task,
    description,
    completed,
    todoListId
) VALUES (
    'Prepare PPT for SIN Review 1',
    'Change 2 slides from TARP PPT',
    FALSE,
    2
);