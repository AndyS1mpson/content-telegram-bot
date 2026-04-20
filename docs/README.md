# Документация

Telegram-бот для сбора обоев с Pinterest и публикации их пачкой в канал.

## Как работает (в двух словах)

1. Владелец шлёт боту `/collect <тема>` — бот логинится в Pinterest через headless Chromium (Playwright) и сохраняет ~30 картинок в PostgreSQL со статусом `New`.
2. Владелец шлёт `/view <тема>` — бот присылает по одному пину с кнопками ❤️ / 👎 / ⏭. Лайк переводит пин в `Selected`, остальное — в `Viewed`. Если `New` закончились, бот сам запускает новый парсинг.
3. Владелец шлёт `/publish <тема>` — бот собирает все `Selected` по теме и отправляет их в канал как media group (пачками по 10), затем проставляет `Posted`.

## Разделы

- **[architecture.md](architecture.md)** — слои приложения, зависимости компонентов, машина состояний пина.
- **[commands.md](commands.md)** — все команды бота, формат аргументов, примеры.
- **[configuration.md](configuration.md)** — `config.yaml`, `.env`, как получить токен, как узнать `chat_id` канала.
- **[development.md](development.md)** — локальный запуск для разработки, структура кода, миграции, как добавить команду.
- **[deployment.md](deployment.md)** — деплой в облако (Oracle Cloud Always Free, Railway, Hetzner), диагностика.

## Стек

- **Go 1.22** — основной язык.
- **PostgreSQL 16** — хранение пинов.
- **Playwright + browserless/chrome** — парсинг Pinterest через headless-браузер.
- **go-telegram-bot-api/v5** — Telegram Bot API (long-polling, публичный URL не нужен).
- **goose** — миграции БД.
- **squirrel** — SQL query builder.
- **logrus** — структурное логирование.
