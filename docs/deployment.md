# Деплой

Бот использует Telegram long-polling — публичный IP/домен **не нужен**. Достаточно хоста с исходящим интернетом и Docker.

Минимум для стека: **1 vCPU, 1.5–2 GB RAM, 5 GB диска**. На 512 МБ не поместится Chromium в `browserless/chrome`.

## Рекомендуемые варианты

| Провайдер | Стоимость | Когда выбрать |
|---|---|---|
| **Oracle Cloud Always Free** (AMD 1 GB + swap, либо Ampere ARM 2×6 GB) | 0 навсегда | Нужно бесплатно, есть карта для верификации |
| **Railway.app** | $5/мес кредита | Не хочется возиться с SSH; готовы платить ~$5 после кредита |
| **Hetzner CAX11** | €4.5/мес | Платно, но максимально стабильно и просто |

Fly.io / Render / Koyeb free-уровни **не подойдут** — 256–512 МБ не хватит headless Chromium'у.

---

## Oracle Cloud Always Free (рекомендую)

### 1. Регистрация
1. [cloud.oracle.com](https://cloud.oracle.com) → **Start for free**.
2. Верифицировать карту. Списаний у Always Free нет; лимиты нельзя превысить случайно.
3. Выбрать регион с наличием Free-инстансов (обычно Frankfurt / Amsterdam / Phoenix).

### 2. Создать VM
1. Compute → Instances → **Create Instance**.
2. Image — **Canonical Ubuntu 22.04**.
3. Shape — **VM.Standard.E2.1.Micro** (AMD, 1/8 OCPU, 1 GB, **Always Free**). Альтернатива — **VM.Standard.A1.Flex** (ARM Ampere, до 4 OCPU/24 GB бесплатно), но `browserless/chrome:1.61` — amd64, ARM потребует другого образа (см. ниже).
4. Добавить свой SSH-публичный ключ.
5. Public IPv4 → **Create**.

### 3. Подготовить VM
```bash
ssh ubuntu@<public-ip>

# Swap обязателен для 1 GB хоста
sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile
sudo mkswap /swapfile && sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# Docker
sudo apt update && sudo apt install -y docker.io docker-compose-v2 git
sudo usermod -aG docker $USER && newgrp docker
```

### 4. Склонировать и заполнить конфиги
```bash
git clone <your-repo-url> content-telegram-bot
cd content-telegram-bot

cp config.example.yaml config.yaml
nano config.yaml         # telegram.api_token, bot_owner_id, telegram_chat_id, login/password Pinterest

cp deployments/.env.example deployments/.env
nano deployments/.env    # смените POSTGRES_PASSWORD на что-то надёжное
```

**Как получить значения для `config.yaml`:**
- `telegram.api_token` — [@BotFather](https://t.me/BotFather) → `/newbot`.
- `telegram.bot_owner_id` — [@userinfobot](https://t.me/userinfobot) пришлёт ваш user_id.
- `telegram_chat_id` канала — добавить бота админом, постнуть что-то в канал,  
  `curl "https://api.telegram.org/bot<TOKEN>/getUpdates"` → взять `chat.id`.

Подробности полей — в [configuration.md](configuration.md).

### 5. Запустить
```bash
cd deployments
docker compose up -d --build
docker compose logs -f bot
```

В Telegram: `/start` → `/collect nature 4k` → `/view nature 4k` → ❤️/👎 → `/publish nature 4k`.

### 6. Автозапуск
У сервисов стоит `restart: unless-stopped`, Docker включён в systemd по умолчанию — после ребута VM стек поднимется сам.

### ARM-вариант (Ampere A1)
Если взяли ARM-шейп, замените в [deployments/docker-compose.yaml](../deployments/docker-compose.yaml):
```yaml
chromeless:
  image: ghcr.io/browserless/chromium:latest    # вместо browserless/chrome:1.61
```
и проверьте `browser_ws` в `config.yaml` — у новой версии эндпойнт обычно `ws://chromeless:3000` без пути `/playwright`. Возможно нужен `TOKEN=...` в env.

---

## Railway.app (быстрый старт без SSH)

Кредита $5/мес хватает на ~1–2 недели работы всего стека 24/7. После этого нужно платить либо выключать.

1. [railway.app](https://railway.app) → Sign in with GitHub.
2. **New Project → Deploy from GitHub repo** → выбрать репозиторий. Railway найдёт [deployments/Dockerfile](../deployments/Dockerfile) и соберёт сервис `bot`.
3. **+ New → Database → PostgreSQL**. Railway выдаст `PGHOST` / `PGUSER` / `PGPASSWORD` — перенесите их в `config.yaml` и прикрепите файл к сервису `bot` через **Volumes** (mount path `/app/config.yaml`).
4. **+ New → Empty Service → Docker Image** → `browserless/chrome:1.61-puppeteer-10.4.0`. В `config.yaml` у `bot` пропишите `browser_ws: ws://<browserless-service>.railway.internal:3000/playwright`.
5. **+ New → Empty Service → Docker Image** → `ghcr.io/kukymbr/goose-docker:3.24.0`. Прокиньте `GOOSE_*` env и примонтируйте `internal/migrations` как volume. Запустите один раз, потом можно отключить.

Railway не съедает весь `docker-compose.yaml` как единое целое — сервисы настраиваются по одному через UI.

---

## Hetzner CAX11 (платно, €4.5/мес, самый простой платный вариант)

ARM Ampere, 2 vCPU / 4 GB RAM. Инструкция полностью совпадает с Oracle от шага **4. Склонировать и заполнить конфиги** (swap не нужен — RAM достаточно). Для ARM используйте `ghcr.io/browserless/chromium:latest` как в заметке выше.

---

## Обслуживание

```bash
# Логи
docker compose logs -f bot
docker compose logs -f chromeless

# Посмотреть состояние БД
docker compose exec postgres psql -U postgres -c \
  "SELECT status, count(*) FROM pin GROUP BY 1;"

# Перечитать config.yaml (бот кэширует его на старте)
docker compose restart bot

# Обновить до последней версии
git pull && docker compose up -d --build

# Полный стоп (данные в volume postgres_data сохраняются)
docker compose down

# Стоп с удалением БД
docker compose down -v
```

## Диагностика

| Симптом | Вероятная причина |
|---|---|
| Бот не отвечает на `/start` | Неверный `api_token` либо вас заблокировал собственный бот — напишите ему первым. `docker compose logs bot`. |
| `you do not have access` на свои команды | `bot_owner_id` в config.yaml не совпадает с вашим user_id. |
| `parse error: sign in: ...` | Pinterest попросил капчу или заблокировал аккаунт. Сменить пароль/аккаунт/IP. Иногда помогает подождать несколько часов. |
| `parse error: wait for images locator` | Страница не отрендерила картинки. Смотрите `docker compose logs chromeless` — возможно таймаут. |
| `/publish`: `chat not found` / `Forbidden` | Бот не админ канала или `telegram_chat_id` неверный. |
| `Out of memory` в логах контейнеров | Недостаточно RAM. Убедитесь что swap подключён (см. шаг 3). На Oracle Free обязателен. |
| Chromeless постоянно рестартует | Часто из-за `/dev/shm` (256 МБ по умолчанию). Флаг `--disable-dev-shm-usage` уже прописан в compose — если всё равно ломается, добавьте в сервис `chromeless`: `shm_size: 1gb`. |
