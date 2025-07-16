-- name: CreateBook :one
INSERT INTO
    books (title, author, price)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: GetBook :one
SELECT * FROM books WHERE id = $1;

-- name: ListBooks :many
SELECT * FROM books ORDER BY created_at DESC;

-- name: UpdateBook :one
UPDATE books
SET
    title = $2,
    author = $3,
    price = $4
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1;