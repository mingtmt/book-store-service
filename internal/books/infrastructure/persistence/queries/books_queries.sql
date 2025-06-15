-- name: CreateBook :one
INSERT INTO books (id, title, author, price)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1;
-- name: ListBooks :many
SELECT *
FROM books
ORDER BY created_at DESC;