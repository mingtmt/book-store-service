-- name: RegisterUser :one
INSERT INTO
    auths (id, username, password)
VALUES ($1, $2, $3)
RETURNING
    *;
-- name: FindByUsername :one
SELECT * FROM auths WHERE username = $1;