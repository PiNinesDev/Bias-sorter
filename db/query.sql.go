// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entry 
WHERE id = (?)
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEntry, id)
	return err
}

const deleteQuiz = `-- name: DeleteQuiz :exec
DELETE FROM quiz 
WHERE id = (?)
`

func (q *Queries) DeleteQuiz(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteQuiz, id)
	return err
}

const findAllQuizes = `-- name: FindAllQuizes :many
SELECT id, name FROM quiz
`

// QUIZ --
func (q *Queries) FindAllQuizes(ctx context.Context) ([]Quiz, error) {
	rows, err := q.db.QueryContext(ctx, findAllQuizes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Quiz
	for rows.Next() {
		var i Quiz
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuiz = `-- name: GetQuiz :one
SELECT id, name FROM quiz
WHERE id = (?)
`

func (q *Queries) GetQuiz(ctx context.Context, id int64) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, getQuiz, id)
	var i Quiz
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getQuizEntries = `-- name: GetQuizEntries :many
SELECT id, quiz_id, name FROM entry
WHERE quiz_id = (?)
`

func (q *Queries) GetQuizEntries(ctx context.Context, quizID int64) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getQuizEntries, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(&i.ID, &i.QuizID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSorter = `-- name: GetSorter :one
SELECT id, type, quiz_id, state FROM sorter
WHERE id = (?)
`

func (q *Queries) GetSorter(ctx context.Context, id int64) (Sorter, error) {
	row := q.db.QueryRowContext(ctx, getSorter, id)
	var i Sorter
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.QuizID,
		&i.State,
	)
	return i, err
}

const newEntry = `-- name: NewEntry :one
INSERT INTO entry (name, quiz_id)
VALUES (?, ?)
RETURNING id, quiz_id, name
`

type NewEntryParams struct {
	Name   string
	QuizID int64
}

// ENTRY --
func (q *Queries) NewEntry(ctx context.Context, arg NewEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, newEntry, arg.Name, arg.QuizID)
	var i Entry
	err := row.Scan(&i.ID, &i.QuizID, &i.Name)
	return i, err
}

const newQuiz = `-- name: NewQuiz :one
INSERT INTO quiz (name)
VALUES (?)
RETURNING id, name
`

func (q *Queries) NewQuiz(ctx context.Context, name string) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, newQuiz, name)
	var i Quiz
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const newSorter = `-- name: NewSorter :one
INSERT INTO sorter (type, quiz_id, state)
VALUES (?, ?, ?)
RETURNING id, type, quiz_id, state
`

type NewSorterParams struct {
	Type   int64
	QuizID int64
	State  sql.NullString
}

func (q *Queries) NewSorter(ctx context.Context, arg NewSorterParams) (Sorter, error) {
	row := q.db.QueryRowContext(ctx, newSorter, arg.Type, arg.QuizID, arg.State)
	var i Sorter
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.QuizID,
		&i.State,
	)
	return i, err
}

const storeSorterState = `-- name: StoreSorterState :exec
UPDATE sorter 
SET state = (?)
WHERE id = (?)
`

type StoreSorterStateParams struct {
	State sql.NullString
	ID    int64
}

func (q *Queries) StoreSorterState(ctx context.Context, arg StoreSorterStateParams) error {
	_, err := q.db.ExecContext(ctx, storeSorterState, arg.State, arg.ID)
	return err
}
