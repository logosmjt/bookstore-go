-- name: CreateUser :one
INSERT INTO users (
  name,
  hashed_password,
  email,
  role
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  name = COALESCE(sqlc.narg(name), name),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at),
  role = COALESCE(sqlc.narg(role), role),
  email = COALESCE(sqlc.narg(email), email)
WHERE id = sqlc.arg(id)
RETURNING *;
