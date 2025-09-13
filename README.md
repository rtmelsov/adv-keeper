# Adv Keeper — мини-хранилище с TUI и gRPC

Adv Keeper — это клиент-серверное приложение для безопасного хранения пользовательских данных и файлов.
Сервер предоставляет gRPC API, а к нему подключается TUI-клиент (Bubble Tea).  
Под капотом: Go, gRPC, Postgres, JWT, потоковая передача файлов (streaming RPC), миграции БД при старте.

![screenshot](docs/screenshot.png)

## Возможности

- Регистрация и вход (JWT access token).
- Хранение файлов на сервере + запись метаданных в Postgres.
- Загрузка файлов **потоком** (без чтения всего файла в память).
- Список файлов, скачивание и удаление.
- TUI-клиент: регистрация/логин, выбор файла через file picker, просмотр списка и деталей.

---

## Архитектура (кратко)

- **Server** (`./cmd/server`):
  - gRPC на TCP (`GRPC_ADDR`).
  - Postgres (`DB_DSN`).
  - Файлы сохраняются на диск (`FILES_DIR`), метаданные — в БД.
  - JWT для авторизации (проверка токена в unary/stream interceptors).
  - Автоматический прогон миграций при старте (`MIGRATIONS_URL`, например `file:///opt/adv-keeper/migrations`).

- **Client (TUI)** (`./cmd/tui`):
  - Подключается к серверу по `AK_GRPC_ADDR` (или `GRPC_ADDR`).
  - Хранит сессию локально: `~/.config/adv-keeper/session.json`.
  - Имеет экраны: Register, Login, Vault (загрузка), FileList, FileDetails.

---

## Требования

- Go 1.22+
- Postgres 14+ (на сервере)
- Для деплоя по SSH/CI — `rsync`/`systemd` (опционально)

---

## Настройка окружения (сервер)

Создай файл окружения (пример для systemd: `/etc/adv-keeper.env`):

```
env

# gRPC адрес (0.0.0.0:8080 — доступ снаружи, 127.0.0.1:8080 — локально)
GRPC_ADDR=0.0.0.0:8080
```
# Строка подключения к Postgres
```
DB_DSN=postgres://admin:admin@127.0.0.1:5432/adv?sslmode=disable
```
# Папка для файлов
```
FILES_DIR=/opt/adv-keeper/data
```
# Секрет для подписи JWT
```
JWT_SECRET=supersecret_change_me
```
# Откуда брать миграции (локальная папка на сервере)
```
MIGRATIONS_URL=file:///opt/adv-keeper/migrations
```
Убедись, что директории существуют и доступны процессу:
```
sudo mkdir -p /opt/adv-keeper/data /opt/adv-keeper/migrations
sudo chown -R ak:ak /opt/adv-keeper
```
Сборка и запуск (сервер локально)
```
go mod tidy
go build -o bin/server ./cmd/server
# запустить
DB_DSN='postgres://admin:admin@127.0.0.1:5432/adv?sslmode=disable' \
FILES_DIR='./data' \
JWT_SECRET='devsecret' \
GRPC_ADDR='127.0.0.1:8080' \
MIGRATIONS_URL='file://./migrations' \
./bin/server
```
Проверить порт:
```
nc -vz 127.0.0.1 8080
```
Сборка и запуск (клиент)
Минимальный Makefile уже есть. Основные цели:
```
# все платформы (бинари в ./dist)
make linux
make mac
make windows

# или просто собрать текущую платформу
make build
./dist/tui
```
Перед запуском укажи адрес сервера:
```
# если сервер на этой же машине
export AK_GRPC_ADDR=127.0.0.1:8080

# если сервер в облаке
export AK_GRPC_ADDR=89.207.255.214:8080
```
