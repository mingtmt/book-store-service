-- name: RegisterUser :one
INSERT INTO
    auths (username, password)
VALUES ($1, $2)
RETURNING
    *;
-- name: FindByUsername :one
SELECT * FROM auths WHERE username = $1;