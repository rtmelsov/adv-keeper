-- 001_init.sql

-- UUID генератор (используем gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ===== USERS =====
CREATE TABLE IF NOT EXISTS users (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email      text NOT NULL UNIQUE,
  pwd_phc    text NOT NULL,
  e2ee_pub   bytea,
  created_at timestamptz NOT NULL DEFAULT now()
);

-- ===== FILES =====
CREATE TABLE IF NOT EXISTS files (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  filename   text NOT NULL,          -- имя файла
  size_bytes bigint NOT NULL,        -- размер в байтах
  path       text NOT NULL,          -- путь/ключ хранения (диск/S3/MinIO)
  created_at timestamptz NOT NULL DEFAULT now()
);

-- Индекс для выборок "мои файлы"
CREATE INDEX IF NOT EXISTS idx_files_user_created_at
  ON files(user_id, created_at DESC);

-- ===== SEED DATA (удали, если не нужно) =====

-- Пара пользователей
INSERT INTO users (email, pwd_phc, e2ee_pub) VALUES
  ('alice@example.com', 'phc$dummy', NULL),
  ('bob@example.com',   'phc$dummy', NULL)
ON CONFLICT (email) DO NOTHING;

-- Тестовые файлы, привязанные к пользователям выше
WITH u AS (
  SELECT id, email FROM users WHERE email IN ('alice@example.com','bob@example.com')
)
INSERT INTO files (user_id, filename, size_bytes, path) VALUES
  ((SELECT id FROM u WHERE email='alice@example.com'), 'report.pdf', 123456, '/data/alice/report.pdf'),
  ((SELECT id FROM u WHERE email='alice@example.com'), 'photo.png',   98765, '/data/alice/photo.png'),
  ((SELECT id FROM u WHERE email='bob@example.com'),   'notes.txt',    4321, '/data/bob/notes.txt');

