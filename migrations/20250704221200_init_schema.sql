-- +goose Up
-- +goose StatementBegin
CREATE TABLE quiz (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE entry (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	quiz_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	FOREIGN KEY (quiz_id) REFERENCES quiz(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE sorter (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type INTEGER NOT NULL,
	quiz_id INTEGER NOT NULL,
	state TEXT,
	FOREIGN KEY (quiz_id) REFERENCES quiz(id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE entry;
-- +goose StatementEnd
