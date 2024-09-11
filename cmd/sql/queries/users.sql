-- name: CreateUser :exec
INSERT INTO users (id, name, email, hashedPassword, createdAt, updatedAt)
VALUES (?, ?, ?, ?, DEFAULT, DEFAULT);

-- name: GetUserById :one
SELECT * FROM users  WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users  WHERE email = ?;

