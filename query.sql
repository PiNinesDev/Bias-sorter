--- USER ---
-- name: NewUser :one
INSERT INTO user (name, password_hash)
VALUES (?,?)
RETURNING *;

-- GetUserByUsername :one
SELECT * FROM user
WHERE name = (?);

-- GetUserByID :one
SELECT * FROM user
WHERE id = (?);

-- name: DeactivateUser :exec
UPDATE user 
SET is_active = 0
WHERE id = (?);

--- SESSION ---
-- name: NewSession :one
INSERT INTO session (token, user_id, expiry)
VALUES (?,?,?)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM session
WHERE token = (?);

-- name: DeleteSession :exec
DELETE FROM session
WHERE token = (?);

-- name: DeleteExpiredSessions :exec
DELETE FROM session
WHERE expiry < (?);

--- QUIZ ---
-- name: NewQuiz :one
INSERT INTO quiz (name, user_id)
VALUES (?, ?)
RETURNING *;

-- name: GetRectentQuizzes :many
SELECT * FROM quiz
WHERE is_active = 1
ORDER BY created_on DESC limit (?) OFFSET (?);

-- name: GetQuizByID :one
SELECT * FROM quiz
WHERE id = (?) AND is_active = 1;;

-- name: GetQuizByUserId :many
SELECT * FROM quiz
WHERE user_id = (?) AND is_active = 1
ORDER BY created_on DESC limit (?) OFFSET (?);

-- name: DeactivateQuiz :exec
UPDATE quiz 
SET is_active = 0
WHERE id = (?) and user_id = (?);

--- ENTRY ---
-- name: NewEntry :one
INSERT INTO entry (name, quiz_id)
VALUES (?, ?)
RETURNING *;

-- name: DeactivateEntry :exec
UPDATE entry 
SET is_active = 0
WHERE id = (?);

-- name: GetQuizEntries :many
SELECT * FROM entry
WHERE quiz_id = (?) AND is_active = 1;

--- SORTER ---
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

