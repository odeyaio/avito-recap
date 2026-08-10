# Avito Recap

«Итоги года» для пользователя Авито — сервис, который на основе истории активности
(просмотры, избранное, контакты, сделки, отзывы) собирает персональный recap: ачивки,
поведенческий тип года, ключевые метрики по категориям и рекомендацию следующего действия.
Recap оформлен как история из карточек (story-формат, как в Instagram/Spotify Wrapped) и
собирается в шеринг-карточку для соцсетей.

Состоит из двух частей: Go-бэкенда, который считает метрики и ачивки по каталогу правил, и
React-фронтенда, который проигрывает готовый recap как интерактивную историю.

## Технологии

**Frontend** (`apps/frontend`)

- React 19 + TypeScript, сборка на Vite
- `react-router` (data-режим, ленивая загрузка страниц через `lazy()`)
- `@tanstack/react-query` — серверный стейт и кэш запросов к API
- MUI (`@mui/material` + `@emotion`) — UI-кит, кастомная тема в `src/theme.ts`
- `orval` — генерация типов и react-query хуков из `contracts/openapi.yaml`
- `msw` (Mock Service Worker) — мок API поверх того же контракта, для автономной разработки и тестов
- `html-to-image` — рендер share-карточки в PNG на клиенте
- `vitest` + `@testing-library/react` — тесты; `eslint` (в т.ч. `eslint-plugin-jsx-a11y`) + `prettier` — линт/форматирование

**Backend** (`apps/backend`)

- Go, HTTP-фреймворк `labstack/echo`
- `oapi-codegen` — генерация серверных типов/интерфейсов из того же `contracts/openapi.yaml`
- PostgreSQL (`jackc/pgx`) — хранилище пользователей, активности, посчитанных recap'ов
- Каталог ачивок и поведенческих типов — декларативно в YAML (`catalog/*.yaml`), заливается сервисом `catalog-import`
- Собственный генератор синтетических тестовых данных (`tools/user_data_generator`)

**Инфраструктура**

- Docker Compose — весь стек одной командой (`postgres` → `migrate` → `catalog-import` → `backend` → `frontend`, плюс одноразовый `generator`)
- Caddy — раздача собранного фронтенда и прокси `/api` на бэкенд в контейнере
- `pnpm` workspaces — монорепо на фронте, `Makefile` — общие команды для всего репозитория

## Структура проекта

```
avito-recap/
├── apps/
│   ├── frontend/                # React-приложение
│   │   ├── src/
│   │   │   ├── pages/           # экраны: Intro → Profiles → Generating → Recap → Share
│   │   │   ├── ui/
│   │   │   │   ├── organisms/   # крупные блоки экрана (RecapDashboard, StoryCardRenderer, ...)
│   │   │   │   ├── molecules/   # переиспользуемые куски (ProfileCard, AchievementBadge, ...)
│   │   │   │   └── templates/   # каркасы layout'ов и обработка ошибок роутов
│   │   │   ├── features/        # изолированная бизнес-логика (story-player, share-export, ...)
│   │   │   ├── api/
│   │   │   │   ├── generated/   # сгенерировано orval'ом из openapi.yaml (не редактируется руками)
│   │   │   │   └── mocks/       # MSW-хендлеры + ручные фикстуры (personas) поверх сгенерированных
│   │   │   ├── routes.tsx       # дерево роутов react-router
│   │   │   └── theme.ts         # тема MUI
│   │   └── package.json
│   └── backend/                 # Go-приложение
│       ├── cmd/                 # точки входа: avito-recap (HTTP-сервер), catalog-import
│       ├── internal/
│       │   ├── adapter/         # in: HTTP-хендлеры (сгенерированные + мапперы); out: репозитории
│       │   ├── engine/          # расчёт метрик, ачивок, поведенческого типа
│       │   ├── service/         # оркестрация: сборка recap из данных + движка
│       │   ├── model/, repository/, catalog/, app/
│       ├── database/migrations/ # SQL-миграции
│       └── tools/user_data_generator/  # генератор синтетических пользователей и активности
├── catalog/                     # декларативный каталог ачивок и поведенческих типов (YAML)
├── contracts/openapi.yaml       # единый контракт API — источник для orval (фронт) и oapi-codegen (бэк)
├── compose.yaml, compose.override.yaml
└── Makefile
```

## Инструкция по запуску

### Через Docker (полный стек)

Поднимает весь стек одной командой: PostgreSQL, миграции, импорт каталога ачивок/поведений
(`catalog/*.yaml`), генератор тестовых пользователей и активности, бэкенд и фронтенд поверх
Caddy. Всё описано в `compose.yaml` (+ `compose.override.yaml`, который переопределяет порт
фронтенда).

**Требования:** Docker и Docker Compose v2 (`docker compose version`); файл `.env` в корне
репозитория — за основу можно взять `.env.example`, дефолты из него уже рабочие.

```bash
docker compose up --build -d   # из корня репозитория
```

Сервисы поднимаются по цепочке зависимостей:

```
postgres → migrate → catalog-import → backend → frontend
                    ↘ generator (сеет тестовых пользователей и активность)
```

`migrate`, `catalog-import` и `generator` — одноразовые джобы: они выполняют свою работу и
завершаются. `Exited (0)` для них в `docker compose ps -a` — это ожидаемое состояние, а не
ошибка. `backend` дожидается успешного завершения `migrate` и `catalog-import`, прежде чем
стартовать.

**Адреса:**

| Сервис   | Адрес                                                       |
| -------- | ------------------------------------------------------------ |
| Frontend | http://localhost:3010 (порт задан в `compose.override.yaml`) |
| Backend  | http://localhost:8080, health-check: `GET /api/v1/health`    |
| Postgres | localhost:5439 (см. `POSTGRES_PORT` в `.env`)                |

**Полезные команды:**

```bash
docker compose ps -a                # статус всех сервисов
docker compose logs backend -f      # логи сервиса (frontend / generator / catalog-import / migrate)
docker compose down                 # остановить и удалить контейнеры (volume с БД сохранится)
```

**Чистый перезапуск (со сбросом БД):**

```bash
docker compose down
docker volume rm avito-recap_postgres-data
docker compose up --build -d
```

После этого `migrate` заново накатит схему, `catalog-import` перезальёт каталог, а `generator`
насеет новых тестовых пользователей и историю их активности.

**Проверка вручную:**

```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/profiles | jq .

curl -X POST http://localhost:8080/api/v1/profiles/<profileId>/recaps \
  -H "Content-Type: application/json" \
  -d '{"year": 2025, "locale": "ru-RU"}'
```

`<profileId>` — любой `id` из ответа `/api/v1/profiles` (все пользователи, засеянные
`generator`, помечены `is_test_profile = true` и попадают в этот список).

### Фронтенд отдельно, без бэкенда

Фронт можно поднять и проверить полностью автономно — на MSW-моках, собранных из того же
`openapi.yaml`, без Docker и без бэкенда вообще.

```bash
pnpm install                      # из корня репозитория
cd apps/frontend
echo "VITE_API_MOCK=true" > .env  # см. .env.example
pnpm dev
```

Vite поднимется на `http://localhost:3000` (если порт занят — сам уйдёт на следующий
свободный, адрес будет в выводе команды).

**Скрипты (`apps/frontend`):**

| Команда         | Что делает                                                          |
| --------------- | -------------------------------------------------------------------- |
| `pnpm dev`      | dev-сервер Vite                                                      |
| `pnpm build`    | `orval generate` → `tsc -b` → `vite build`                           |
| `pnpm generate` | только генерация API-клиента и MSW-моков из `contracts/openapi.yaml` |
| `pnpm test`     | `vitest run`                                                         |
| `pnpm lint`     | `eslint .`                                                           |
| `pnpm preview`  | локальный просмотр собранного билда                                 |

## Особенности реализации (фронтенд)

Пользовательский флоу — пять экранов на `react-router` (data-режим, каждая страница — свой
lazy-чанк): `Intro → Profiles → Generating → Recap (story-player) → Share`. `RecapPage` и
`SharePage` используют общий `recapLoader`, вынесенный в отдельный лёгкий модуль — переход по
шаренной ссылке на `/recap/:id/share` не тянет за собой код story-плеера, который на этом
экране не нужен.

- **Story-player.** `buildStorySlides` (`src/features/story-player`) собирает из ответа API
  плоский список слайдов (ачивки, поведенческий тип, ключевые метрики) для проигрывания
  историей; `useStoryPlayer` держит состояние текущего слайда, таймер автопрогресса и свайпы;
  `usePrefersReducedMotion` отключает автопрогресс и анимации для пользователей с
  `prefers-reduced-motion`.
- **Контракт и генерация.** Типы, react-query хуки (`src/api/generated`) и «сырые» MSW-хендлеры
  (`src/api/mocks/generated`) генерируются `orval`'ом из `contracts/openapi.yaml` и не
  редактируются руками — при изменении контракта всё пересобирается заново через `pnpm generate`.
- **Моки поверх контракта.** `src/api/mocks/fixtures/personas.ts` — три тестовых профиля с
  внутренне согласованными метриками, каждый подобран так, чтобы триггерить свой набор ачивок
  и свой поведенческий тип из `catalog/*.yaml`: **Алексей** (Исследователь), **Марина**
  (Целеустремлённый покупатель), **Игорь** (Эффективный продавец). `src/api/mocks/handlers.ts`
  переопределяет сгенерированные хендлеры этими фикстурами и добавляет два «скрытых» profile id
  (не показаны в списке профилей, доступны по прямому URL) — для проверки состояний 422
  (недостаточно активности) и 503 (каталог не загружен) без реального бэкенда. Переключатель —
  `VITE_API_MOCK=true` в `.env`; в проде код мока не попадает в бандл вовсе.
- **Share-карточка.** `useShareCard` (`src/features/share-export`) рендерит карточку в PNG через
  `html-to-image` и делится ей с последовательной деградацией: `navigator.canShare({files})` →
  `navigator.share()` → копирование ссылки в буфер обмена, — в зависимости от того, что
  поддерживает браузер.
- **Доступность (a11y).** Линт держит это на уровне правил (`eslint-plugin-jsx-a11y`), плюс
  ручная работа: уважение `prefers-reduced-motion` в story-player, `ExplanationDialog` и
  `RouteErrorBoundary` как доступные фолбэки, тесты на ключевые интерактивные компоненты через
  `@testing-library/react` (запросы по ролям/лейблам, а не по классам/тексту реализации).

### Ограничения MVP

- Фронт разрабатывался и гоняется тестами против MSW-моков контракта.
- `nextAction.href` в тестовых фикстурах указывает на пути вида `/search?...`,
  `/seller/statistics` — это ожидаемый формат deep link'а в основное приложение Авито, а не
  реализованные в этом MVP маршруты; сам recap-фронт их не рендерит.
- Экспорт share-карточки в PNG и шаринг через Web Share API зависят от возможностей
  браузера — см. деградацию выше.
- Клиентский AI-чат по recap исследовали отдельно и не стали включать в MVP — подробности и
  причина ниже, в «Использование ИИ».
- Только один язык интерфейса (русский), без i18n.

## Использование ИИ

### Клиентский чат поверх Chrome Prompt API (исследовано, не включено в MVP)

В ветке `feat/frontend-ai-chat` пробовали добавить чат-виджет, который отвечал бы на вопросы
пользователя про его recap полностью на клиенте — через встроенный в браузер Chrome Prompt API
(`window.LanguageModel`, on-device Gemini Nano), без единого обращения к бэкенду. Идея была в
том, чтобы контекстом для модели служили только уже показанные пользователю текстовые поля
recap'а (заголовок, тексты карточек, объяснения ачивок и поведения) — без сырых `metrics` и без
сетевого запроса, покидающего устройство.

Реализация была рабочей: feature-detection через `'LanguageModel' in window`, проверка
`availability()`, потоковый ответ через `promptStreaming()`, полное скрытие виджета при
отсутствии API. При проверке в реальном Chrome выяснилось, что Prompt API на сегодняшний день
поддерживает вход и выход только на пяти языках — `de`, `en`, `es`, `fr`, `ja`. Русского среди
них нет. Явный запрос `languages: ["ru"]` браузер отклоняет:

```
Unsupported LanguageModel API languages were specified, and the request was aborted.
API calls must only specify supported languages to ensure successful processing and
guarantee output characteristics. Please only specify supported language codes:
[de, en, es, fr, ja]
```

Без явного указания языка модель либо не отвечает, либо ведёт себя непредсказуемо — сам Chrome
предупреждает об этом в консоли ещё до отказа.

Поскольку весь продукт и recap — на русском, честно работающей фичи не получилось: показывать
пользователю нерабочий виджет или переключать чат на английский в русском продукте — хуже, чем
не показывать ничего. Фичу решили не включать в MVP.

## Команда

- **Кирилл** — фронтенд: весь UI и пользовательский флоу (`apps/frontend`), story-player,
  генерация API-клиента и MSW-моков из контракта, тестовые персоны и фикстуры, доступность
  (a11y), тесты (`vitest` + `@testing-library/react`), экспорт и шаринг карточки результата.
