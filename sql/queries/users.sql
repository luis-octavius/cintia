-- name: CreateUser :exec
INSERT INTO users (id, name, created_at, updated_at, password_hash)
VALUES (
  gen_random_uuid(),
  $1, 
  $2, 
  $3,
  $4
);

-- name: GetUserByID :one 
SELECT id, name, created_at, updated_at 
FROM users WHERE id = $1; 




