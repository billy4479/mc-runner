-- name: GetUserFromToken :one
SELECT sqlc.embed(users), token_store.expires FROM users
INNER JOIN token_store ON token_store.user_id = users.id
WHERE token_store.token = ? AND token_store.type = ?;

-- name: GetAllTokensForUser :many
SELECT * FROM token_store
WHERE user_id = ?;

-- name: SetToken :exec
INSERT INTO token_store (token, expires, type, user_id)
VALUES (?, ?, ?, ?);

-- name: RemoveTokenById :exec
DELETE FROM token_store
WHERE user_id = ? and type = ?;

-- name: RemoveTokenExact :exec
DELETE FROM token_store
WHERE token = ? and user_id = ?;

-- name: CreateUser :one
INSERT INTO users (name)
VALUES (NULL)
RETURNING id;

-- name: SetUserName :exec
UPDATE users
SET name = ?
WHERE id = ?;
