-- QUIZ --
-- name: FindAllQuizes :many
SELECT * FROM quiz;

-- name: GetQuiz :one
SELECT * FROM quiz
WHERE id = (?);

-- name: NewQuiz :one
INSERT INTO quiz (name)
VALUES (?)
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM quiz 
WHERE id = (?);

-- ENTRY --
-- name: NewEntry :one
INSERT INTO entry (name, quiz_id)
VALUES (?, ?)
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entry 
WHERE id = (?);

-- name: GetQuizEntries :many
SELECT * FROM entry
WHERE quiz_id = (?);

-- name: GetSorter :one
SELECT * FROM sorter
WHERE id = (?);

-- name: NewSorter :one
INSERT INTO sorter (type, quiz_id, state)
VALUES (?, ?, ?)
RETURNING *;

-- name: StoreSorterState :exec
UPDATE sorter 
SET state = (?)
WHERE id = (?);

