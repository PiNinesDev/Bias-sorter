-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = ON;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE quiz (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAP,
	is_active BOOLEAN NOT NULL DEFAULT 1,
	FOREIGN KEY (user_id) REFERENCES user (id),
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE entry (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	quiz_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAP,
	is_active BOOLEAN NOT NULL DEFAULT 1
	FOREIGN KEY (quiz_id) REFERENCES quiz(id) ON DELETE CASCADE,
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAP,
	is_active BOOLEAN NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE session (
	token TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expiry TIMESTAMP NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE entry;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE session;
-- +goose StatementEnd
