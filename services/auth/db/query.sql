-- name: CreateUser :exec
INSERT INTO users (username, password_hash) VALUES ($1, $2);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;