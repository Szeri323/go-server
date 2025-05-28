-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, hashed_password) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name=$1;

-- name: UpdateUser :one
UPDATE users SET name=$1, hashed_password=$2, updated_at=$3 WHERE id=$4 RETURNING *;

-- name: UpdateUsersMembership :exec
UPDATE users SET is_chirpy_red=$1 WHERE id=$2;

-- name: TruncateUsersTable :exec
TRUNCATE TABLE users CASCADE;