-- name: Register :one
INSERT INTO users (email, pwd_phc, e2ee_pub)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, email, pwd_phc, e2ee_pub, created_at
FROM users
WHERE email = $1;

-- name: AddFile :one
INSERT INTO files (user_id, filename, path, size_bytes)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, filename, path, size_bytes, created_at;

-- name: ListFilesByUser :many
SELECT id, filename, path, size_bytes, created_at
FROM files
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteFile :one
DELETE FROM files
WHERE id = $1 AND user_id = $2
RETURNING id;-- name: GetUserByID :one

-- name: GetUserByID :one
SELECT id, email, created_at
FROM users
WHERE id = $1;

-- name: GetFileForUser :one
SELECT id, user_id, filename, path, size_bytes, created_at
FROM files
WHERE id = $1 AND user_id = $2;
