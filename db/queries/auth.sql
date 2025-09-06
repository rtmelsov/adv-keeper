-- name: RegisterWithDevice :one
WITH u AS (
  INSERT INTO users (email, pwd_phc, e2ee_pub)
  VALUES ($1, $2, $3)
  RETURNING id
),
d AS (
  INSERT INTO devices (user_id, device_id)
  SELECT u.id, $4
  FROM u
  RETURNING device_id
)
SELECT
  (SELECT id FROM u)        AS user_id,
  (SELECT device_id FROM d) AS device_id;

