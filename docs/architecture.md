# Архитектура

## Слои

```
cmd/main.go                          — точка входа, DI
└── internal/
    ├── config/                      — загрузка config.yaml
    ├── infrastructure/
    │   ├── storage/                 — PostgreSQL pool (sqlx)
    │   └── clients/
    │       ├── browser/             — Playwright-клиент (WebSocket к browserless)
    │       └── telegram/            — tgbotapi.BotAPI обёртка
    ├── models/                      — Pin, Account, Channel, PinStatus
    ├── repository/                  — CRUD пинов (squirrel + sqlx)
    ├── service/
    │   ├── parser/pinterest/        — парсер Pinterest через Playwright
    │   ├── pin/                     — бизнес-логика пинов
    │   └── telegram/                — handlers команд и callback'ов
    └── utils/                       — log, slices, save
```

Зависимости идут только «внутрь»: handler'ы Telegram зависят от `pin.Service` через интерфейс `pinService`, `pin.Service` — от `parser` и `repository` через свои интерфейсы. Инфраструктурные клиенты (Playwright, Postgres, tgbotapi) инстанцируются в `main.go` и передаются конструкторами.

## Граф сборки в main.go

```
config.NewConfig()
  → storage.New(cfg.Database)              — *sqlx.DB
  → repository.New(db)                     — Repository
  → browser.NewBrowser(cfg.BrowserWS)      — playwright.Browser
  → pinterest.New(browser)                 — Parser
  → pin.NewService(parser, repo)           — *pin.Service
  → tgClient.New(cfg.Telegram.Token)       — *tgbotapi.BotAPI
  → telegram.New(bot, pinService, accounts, cfg.Telegram)
  → tgWrapper.RegisterHandlers(ctx)        — блокирующий цикл long-polling
```

## Модель данных

Одна таблица `pin`:

| колонка | тип | назначение |
|---|---|---|
| `id` | BIGINT | внешний ID пина в Pinterest |
| `url` | TEXT | прямой URL оригинала картинки (i.pinimg.com/originals/…) |
| `type` | TEXT | `pin` или `video` |
| `status` | INTEGER | 1=New, 2=Viewed, 3=Posted, 4=Selected |
| `channel` | TEXT | имя Telegram-канала из конфига |
| `query` | TEXT | поисковая тема, по которой пин собран |
| `created_at` | TIMESTAMPTZ | для возможной ретроспективы |

`PRIMARY KEY (id, channel)` — защита от дубликатов при повторном `/collect <тема>`. `CreatePins` вставляет с `ON CONFLICT DO NOTHING`.

## Машина состояний пина

```
          /collect                     ❤️ like
Pinterest ────────▶  New  ──/view──▶  Viewed  ──────▶  Selected
                             ▲                            │
                             │ 👎/⏭ остаётся Viewed       │
                                                         /publish
                                                          ▼
                                                        Posted
```

- **Переход `New → Viewed` атомарен внутри `GetPinsForView`**: пин помечается `Viewed` в той же функции, что его вернула — защита от повторного показа при спаме команды.
- **`Viewed → Selected`** делает callback `like` (`pinService.Select`).
- **`Selected → Posted`** делает `/publish` после успешной отправки media group.
- Обратных переходов нет — если пользователь нажал 👎, пин навсегда остаётся `Viewed` для этой `(channel, query)`.

## Потоки данных

### /collect \<query\>
```
Update → CollectHandler
       → pinService.Parse(account, query)
         → parser.Parse(account, query)                 — Playwright
           1. getNewPage()                              — новая страница через WS
           2. signIn(page, account)                     — логин по форме
           3. gotoSearch(page, query) если query != ""  — search URL Pinterest
           4. evaluate(getPinInfoRowFunc)               — JS извлечение img[srcset]
           5. transformImageURL(...)                    — 236x/564x → originals
         → repository.CreatePins(pins)                  — INSERT ... ON CONFLICT DO NOTHING
       → sendMessage("Готово...")
```

### /view \<query\>
```
Update → ViewHandler
       → showNextPin(chatID, channel, query)
         → pinService.GetPinsForView(filter{New, channel, query, limit:1})
           ├─ errgroup: CountPins + GetPins                  — параллельно
           └─ UpdateStatuses([id], Viewed)                   — атомарно перед выдачей
         → если пусто → pinService.Parse(...)                — автодозапуск
         → sendPinWithCheckboxes(chatID, pin, count-1)
           → генерирует callback_data "like:<id>"
           → sendContent(photo/video + inline keyboard)
```

### callback (like / dislike / skip)
```
Update.CallbackQuery → CallbackHandler
  → parseCallback(data) → action, pinID
  → pinService.GetByID(pinID)                         — узнать query пина
  → action=like   → pinService.Select(pinID)          — статус Selected
    action=skip   → ничего (уже Viewed)
    action=dislike→ ничего (уже Viewed)
  → bot.Request(NewCallback(...))                     — убрать спиннер
  → bot.Request(EditMessageCaption ✅/❌/⏭)            — снять клавиатуру, добавить отметку
  → showNextPin(chatID, channel, pin.Query)           — той же функцией, что и /view
```

### /publish \<query\>
```
Update → PublishHandler
       → pinService.GetSelected(channel, query)              — WHERE status=4 AND channel=? AND query=?
       → chunk pins по mediaGroupLimit (=10)
       → для каждого чанка:
         → buildMediaGroup(chunk)                            — []InputMediaPhoto / InputMediaVideo
         → bot.SendMediaGroup(NewMediaGroup(chat_id, media)) — Telegram сам скачивает по URL
       → pinService.MarkPosted(posted_ids)                   — UPDATE status=3
       → sendMessage("Опубликовано N пинов")
```

## Формат callback_data

Формат `<action>:<pinID>`, например `like:123456789`. Тема (`query`) в callback не передаётся — её берут из БД по ID через `pinService.GetByID`. Это обход 64-байтного лимита Telegram на `callback_data` (русскоязычные темы легко его превышают).

## Безопасность

Все handler'ы начинаются с `validateUser(update.Message.From.ID)` — сверка с `bot_owner_id` из конфига. Посторонний юзер получает `ErrAccessDenied`. Callback-кнопки тоже проверяются: даже если кто-то перешлёт сообщение с кнопками, нажатие отвергнется.

Playwright подключается к **существующему** browserless-контейнеру через WebSocket (`ws://chromeless:3000/playwright`) — бот не тянет браузерные бинарники в свой образ. Его итоговый размер ~17 МБ.

## Что сознательно не сделано

- **FSM / сессии пользователя** — тема передаётся аргументом каждой команды, глобального состояния нет.
- **Мультиканальность в UI** — в конфиге может быть несколько `accounts`, но бот берёт первый. Добавить `/channel <name>` — задача на будущее.
- **Автолайки в Pinterest** — `parser.LikePins` реализован, но не вызывается. Цель — прогреть ленту рекомендаций после успешной публикации; можно добавить отдельной командой `/train`.
- **Планировщик публикаций** — закомментированный `gocron` в `main.go` намекал на крон, но это не вошло в MVP.
