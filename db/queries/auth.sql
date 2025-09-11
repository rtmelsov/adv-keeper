-- name: RegisterWithDevice :one
INSERT INTO users (email, pwd_phc, e2ee_pub)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, email, pwd_phc, e2ee_pub, created_at
FROM users
WHERE email = $1;

-- name: AddFile :one
INSERT INTO files (user_id, name, path)
VALUES ($1, $2, $3)
RETURNING id, user_id, name, path, created_at;

-- name: ListFilesByUser :many
SELECT id, name, path, created_at
FROM files
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1 AND user_id = $2;
