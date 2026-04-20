# Разработка

## Требования

- Go 1.22+
- Docker + Docker Compose (для postgres и browserless)
- [goose](https://github.com/pressly/goose) для миграций: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Быстрый старт локально

```bash
git clone <repo> && cd content-telegram-bot

# 1. Подготовить конфиги
cp config.example.yaml config.yaml
cp deployments/.env.example deployments/.env
# в config.yaml: host → localhost, browser_ws → ws://localhost:3000/playwright

# 2. Поднять инфраструктуру (без бота)
cd deployments && docker compose up -d postgres chromeless
cd ..

# 3. Накатить миграции
goose -dir internal/migrations postgres \
  "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable" up

# 4. Запустить бота
GO111MODULE=on go run ./cmd
```

При `go run` проект подхватывает `config.yaml` из текущей директории — запускайте из корня репозитория.

## Сборка

```bash
GO111MODULE=on go build -o bin/bot ./cmd    # локальный бинарник
cd deployments && docker compose build bot  # docker-образ (~17 МБ на выходе)
```

## Проверки

```bash
GO111MODULE=on go vet ./...
GO111MODULE=on go build ./...
```

Unit-тестов в проекте пока нет. Если добавляете — кладите рядом с кодом (`_test.go`).

## Миграции

Формат [goose](https://github.com/pressly/goose/blob/master/docs/goose.md). Лежат в `internal/migrations/`, имя файла `YYYYMMDDhhmmss_<name>.sql`.

Шаблон:
```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE pin ADD COLUMN something TEXT;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE pin DROP COLUMN something;
-- +goose StatementEnd
```

Применить:
```bash
goose -dir internal/migrations postgres "$DSN" up
goose -dir internal/migrations postgres "$DSN" status
goose -dir internal/migrations postgres "$DSN" down     # откатить последнюю
```

В docker-compose миграции накатывает сервис `migrations` (образ `kukymbr/goose-docker`) один раз перед стартом бота — см. [deployment.md](deployment.md).

## Структура кода

См. [architecture.md](architecture.md) для слоёв. Ниже — практические подсказки где что менять:

| Хочу… | Файл |
|---|---|
| Добавить новую команду `/foo` | `internal/service/telegram/const.go` (константа), `internal/service/telegram/foo.go` (handler), `internal/service/telegram/client.go` (`handleMessage` switch) |
| Добавить callback-действие | `internal/service/telegram/const.go` (префикс), `internal/service/telegram/callback.go` (ветка switch) |
| Новое поле у пина | миграция + `internal/models/pin.go` + `internal/repository/entities.go` + `internal/repository/get.go` / `create.go` |
| Новый фильтр при выборке пинов | `internal/service/pin/filter.go` (поле) + `internal/repository/get.go` (`applyFilter`) |
| Поменять селекторы на Pinterest | `internal/service/parser/pinterest/parse.go` (`getPinInfoRowFunc`, `gotoSearch`) |
| Добавить таймаут / retry на парсинг | `internal/service/parser/pinterest/parser.go` (`signIn`, `gotoSearch`) |

## Паттерны

- **Интерфейсы зависимостей пакета — в `deps.go`.** Пример: `internal/service/pin/deps.go`, `internal/service/telegram/deps.go`. Так пакет не знает конкретные реализации — они передаются из `main.go`.
- **Репозиторий возвращает модели, не entity.** Entity (`internal/repository/entities.go`) — только для маппинга БД. Конверсия — в `get.go`.
- **Бизнес-логика в `service/*`**, handler'ы Telegram — тонкая прослойка: валидация юзера → вызов `pinService.*` → ответ в чат.
- **Ошибки оборачиваются `errors.Wrap(err, "context")`** через `github.com/pkg/errors`. В логах используется `log.Error(err, log.Data{...})` из `internal/utils/log`.

## Отладка

- Логи бота — stdout. При `go run` идут в терминал, в Docker — `docker compose logs -f bot`.
- Запросы Telegram API бот не логирует. Если надо — включите `bot.Debug = true` на `*tgbotapi.BotAPI` (в `internal/infrastructure/clients/telegram/client.go`).
- Парсинг Pinterest падает чаще всего из-за капчи. Чтобы увидеть, что происходит в браузере, запустите browserless с `CONNECTION_TIMEOUT=0 PREBOOT_CHROME=false` и откройте `http://localhost:3000` — там UI со screenshot'ами активных сессий.

## Типовые задачи

### Добавить команду /stats \<query\>

1. Константа в `const.go`:
   ```go
   CommandStats Command = "/stats"
   ```
2. Метод в `pinService` (deps.go) и реализация в `internal/service/pin/stats.go`:
   ```go
   Stats(ctx, channel, query) (map[models.PinStatus]int64, error)
   ```
3. Handler `internal/service/telegram/stats.go` — валидация, вызов, форматирование ответа.
4. Ветка в `client.go` `handleMessage`:
   ```go
   case CommandStats: c.StatsHandler(ctx, update)
   ```

### Перевести `channel` в отдельный аргумент команды

1. В handler'ах парсить `update.Message.CommandArguments()` так, чтобы первое слово — ключ канала, остальное — query.
2. Сделать ресолв `c.accounts[models.Channel(key)]` с падением в `ErrIncorrectAction`, если ключ не найден.
3. `defaultAccount` можно оставить как fallback для случая одного аккаунта.

Это примерно 30 строк изменений, ломать ничего не нужно.
