Adv Keeper — «еле-еле, но работает»

Маленький демо-проект: gRPC-сервер + TUI-клиент (Go) для регистрации/логина и загрузки/скачивания файлов.
План был «сделать за неделю», получилось «сделали — и это уже праздник»

Что умеет

Регистрация и логин (пароли в БД, клиент хранит только сессию).

JWT-авторизация на сервере.

Загрузка/скачивание файлов с прогрессом.

Кроссплатформенный TUI (Linux / macOS / Windows).

Автомиграции БД при старте сервера.

systemd-юнит для деплоя (чтобы «само вставало и работало… ну почти»).

Быстрый старт (клиент)

Клиенту нужна одна переменная окружения: AK_GRPC_ADDR — адрес gRPC сервера, например 89.207.255.214:8080.

Windows (PowerShell)
$env:AK_GRPC_ADDR = "89.207.255.214:8080"
.\tui_windows_amd64.exe


Сделать навсегда:

[Environment]::SetEnvironmentVariable(
  "AK_GRPC_ADDR","89.207.255.214:8080",[EnvironmentVariableTarget]::User
)

macOS / Linux
export AK_GRPC_ADDR=89.207.255.214:8080
./tui                   # или ./tui_darwin_arm64, ./tui_linux_amd64


Опционально: куда сохранять скачанные файлы

export AK_FILES_DIR="$HOME/Downloads/adv-keeper"


Навигация в TUI: ↑/↓ и Enter. После логина сессия лежит в:

Linux/macOS: ~/.config/adv-keeper/session.json

Windows: %AppData%\adv-keeper\session.json

Сборка клиента

Требуется Go 1.22+.

make build    # текущая платформа -> ./dist/tui
make linux    # linux_amd64
make mac      # darwin_arm64 (Apple Silicon)
make windows  # windows_amd64.exe

Сервер (шпаргалка владельцу)
.env
# /etc/adv-keeper.env
AK_GRPC_ADDR=0.0.0.0:8080
DB_DSN=postgres://admin:admin@127.0.0.1:5432/adv?sslmode=disable
JWT_SECRET=super-secret-not-for-git
FILES_DIR=/opt/adv-keeper/data
MIGRATIONS_PATH=file:///opt/adv-keeper/migrations

systemd unit
# /etc/systemd/system/adv-keeper.service
[Unit]
Description=Adv Keeper Server
After=network.target
Requires=postgresql.service

[Service]
User=ak
Group=ak
EnvironmentFile=/etc/adv-keeper.env
WorkingDirectory=/opt/adv-keeper
ExecStart=/opt/adv-keeper/current
Restart=always
RestartSec=2
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target


Развёртывание (у нас через GitHub Actions) кладёт бинарь в /opt/adv-keeper/releases/<server-TS>, двигает current и рестартит сервис.
Миграции берутся из MIGRATIONS_PATH.

Проверки:

sudo systemctl status adv-keeper --no-pager
journalctl -u adv-keeper -n 100 --no-pager
ss -ltnp 'sport = :8080'         # слушает ли порт


Открыть порт наружу (UFW):

sudo ufw allow 8080/tcp
sudo ufw status

FAQ

connection refused
На сервере AK_GRPC_ADDR=0.0.0.0:8080, порт открыт, сервис запущен.

Дисклеймер безопасности

Это демо. Порт 8080 может быть открыт, JWT секрет — в переменных окружения, логика — нарочно упрощена.
Не храните чувствительные данные, только тестовые!!!

Статус

«Работает у меня на облаке»

⚠️ «Улучшить можно всё»

🥳 «Но оно живёт, и это главное»

Если что-то не взлетело — заводите issue. Если взлетело — тоже заводите, нам нужна моральная поддержка. 🙌
