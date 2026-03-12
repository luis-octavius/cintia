-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, password_hash, role, created_at, updated_at;

-- name: GetUserByID :one 
SELECT id, name, email, password_hash, role, created_at, updated_at 
FROM users 
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, created_at, updated_at 
FROM users 
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users 
SET 
  name = COALESCE(sqlc.narg('name'), name),
  email = COALESCE(sqlc.narg('email'), email),
  password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
  updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, password_hash, role, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1; 




