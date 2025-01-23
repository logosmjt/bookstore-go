-- name: CreateBook :one
INSERT INTO books (
  title,
  author,
  price,
  description,
  cover_image_url,
  published_date,
  user_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetBook :one
SELECT * FROM books
WHERE id = $1 LIMIT 1;

-- name: ListBooks :many
SELECT * FROM books
WHERE user_id = $1
ORDER BY published_date
LIMIT $2
OFFSET $3;

-- name: UpdateBook :one
UPDATE books
SET
  title = COALESCE(sqlc.narg(title), title),
  author = COALESCE(sqlc.narg(author), author),
  price = COALESCE(sqlc.narg(price), price),
  description = COALESCE(sqlc.narg(description), description),
  cover_image_url = COALESCE(sqlc.narg(cover_image_url), cover_image_url),
  published_date = COALESCE(sqlc.narg(published_date), published_date)
WHERE id = sqlc.arg(id)
RETURNING *;