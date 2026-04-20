# Конфигурация

Бот читает два файла:
- `config.yaml` в корне проекта — секреты и настройки приложения.
- `deployments/.env` — креденшелы PostgreSQL для docker-compose (их же подхватывает сервис `migrations`).

Оба файла **не должны попадать в git** — пользуйтесь `config.example.yaml` и `deployments/.env.example`.

## config.yaml

```yaml
telegram:
  api_token: 123456:ABC-your-tg-bot-token
  bot_owner_id: 111111111

accounts:
  - channel: Wall Paper
    telegram_chat_id: -1001234567890
    login: pinterest_account_login
    password: pinterest_account_password

database:
  user: postgres
  db: postgres
  password: postgres
  host: postgres        # "localhost" при локальном go run, "postgres" при docker-compose
  port: 5432

browser_ws: ws://chromeless:3000/playwright   # "ws://localhost:3000/playwright" локально
```

### telegram.api_token
Токен бота. Получают у [@BotFather](https://t.me/BotFather):
```
/newbot → имя бота → username (должен заканчиваться на _bot) → токен вида 123456:ABC...
```
После этого желательно `/setprivacy → Disable` (иначе бот в группах не будет видеть сообщения — для нашего сценария это не критично, но не мешает).

### telegram.bot_owner_id
Ваш Telegram user ID. Узнать можно у [@userinfobot](https://t.me/userinfobot) — он ответит числом. Нужен потому что все команды доступны только владельцу.

### accounts[*]
Список «связок» Pinterest-аккаунт ↔ Telegram-канал. Сейчас бот использует только первый элемент списка (MVP), но структура готова к мульти-каналу.

| Поле | Назначение |
|---|---|
| `channel` | Человекочитаемое имя. Пишется в БД в колонку `channel` — по нему фильтруется выдача. |
| `telegram_chat_id` | `chat_id` вашего канала в Telegram, куда публиковать (см. ниже). |
| `login` | Логин Pinterest — email или username. |
| `password` | Пароль Pinterest. |

**Как узнать `telegram_chat_id` канала:**
1. Добавьте бота администратором в канал (права: `Post Messages`).
2. Запостите любое сообщение в канал.
3. `curl "https://api.telegram.org/bot<TOKEN>/getUpdates"` — в ответе найдите `"chat": {"id": -1001234567890, ...}`. Это и есть `telegram_chat_id`.

Для публичных каналов можно использовать `@username` вместо числа, но лучше числовой — он работает для приватных каналов тоже.

### database
Параметры подключения к PostgreSQL. `host: postgres` работает внутри docker-compose (имя сервиса). Для локального запуска `go run ./cmd` вне Docker поменяйте на `localhost`.

### browser_ws
WebSocket-эндпойнт к browserless-контейнеру. Playwright подключается к нему вместо запуска браузера локально. Аналогично database.host: в compose — `ws://chromeless:3000/playwright`, локально — `ws://localhost:3000/playwright`.

## deployments/.env

```dotenv
POSTGRES_USER=postgres
POSTGRES_DB=postgres
POSTGRES_PASSWORD=change-me
```

Читается docker-compose для сервисов `postgres` и `migrations`. Значения должны совпадать с блоком `database` в `config.yaml` — ничего не синхронизируется автоматически.

**Прод-рекомендации:**
- `POSTGRES_PASSWORD` — сгенерируйте через `openssl rand -base64 24`.
- Не открывайте порт 5432 наружу. В [deployments/docker-compose.yaml](../deployments/docker-compose.yaml) его нет в `ports` специально.

## Где правильно держать секреты

- `config.yaml` и `deployments/.env` — **вне git**. Оба пути уже в `.gitignore` (если его нет — добавьте).
- На VM храните их с правами `600` для текущего пользователя.
- Для «серьёзной» прода переносите `telegram.api_token` и Pinterest-пароль в секрет-менеджер (Docker secret, Vault и т. п.) — у нас пока такого не сделано.

## Где НЕ нужно ничего менять

- Номера статусов пинов, имена колонок БД, таймауты Playwright — они в коде, не в конфиге. Намеренно: «магические» настройки у конфига превращают его в свалку.
- Константы `pinterestLoginURL`, `pinterestSearchURL` — хардкод в `internal/service/parser/pinterest/parser.go`.

## Пример для локального запуска без Docker

`config.yaml`:
```yaml
database:
  user: postgres
  db: postgres
  password: postgres
  host: localhost
  port: 5432

browser_ws: ws://localhost:3000/playwright
```

Потом поднимаете только инфраструктуру (postgres + chrome) в compose, а бот запускаете `go run ./cmd`. Детали — в [development.md](development.md).
