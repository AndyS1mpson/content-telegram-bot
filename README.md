# wallpaper-bot

Telegram-бот для сбора обоев с Pinterest и публикации их пачкой в канал.

**TL;DR:** `/collect <тема>` — парсит Pinterest; `/view <тема>` — листает картинки с кнопками ❤️/👎/⏭; `/publish <тема>` — публикует отобранные в канал через `sendMediaGroup`.

## Быстрый старт (docker-compose)

```bash
cp config.example.yaml config.yaml          # заполнить токены и данные канала
cp deployments/.env.example deployments/.env # пароль БД
cd deployments && docker compose up -d --build
```

## Документация

- **[docs/](docs/README.md)** — индекс.
- [docs/architecture.md](docs/architecture.md) — слои, компоненты, машина состояний пина.
- [docs/commands.md](docs/commands.md) — команды бота.
- [docs/configuration.md](docs/configuration.md) — `config.yaml`, `.env`, где взять токен и chat_id.
- [docs/development.md](docs/development.md) — локальная разработка, миграции, структура кода.
- [docs/deployment.md](docs/deployment.md) — деплой в Oracle Cloud Always Free / Railway / Hetzner.

## Стек

Go 1.22 · PostgreSQL 16 · Playwright + browserless/chrome · go-telegram-bot-api/v5 · goose · squirrel.
