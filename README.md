# Avito Recap

## Использование ИИ

### Клиентский чат поверх Chrome Prompt API (исследовано, не включено в MVP)

В ветке `feat/frontend-ai-chat` пробовали добавить чат-виджет, который отвечал бы на вопросы
пользователя про его recap полностью на клиенте - через встроенный в браузер Chrome Prompt API
(`window.LanguageModel`, on-device Gemini Nano), без единого обращения к бэкенду. Идея была в
том, чтобы контекстом для модели служили только уже показанные пользователю текстовые поля
recap'а (заголовок, тексты карточек, объяснения ачивок и поведения) - без сырых `metrics` и без
сетевого запроса, покидающего устройство.

Реализация была рабочей: feature-detection через `'LanguageModel' in window`, проверка
`availability()`, потоковый ответ через `promptStreaming()`, полное скрытие виджета при
отсутствии API. При проверке в реальном Chrome выяснилось, что Prompt API на сегодняшний день
поддерживает вход и выход только на пяти языках - `de`, `en`, `es`, `fr`, `ja`. Русского среди
них нет. Явный запрос `languages: ["ru"]` браузер отклоняет:

```
Unsupported LanguageModel API languages were specified, and the request was aborted.
API calls must only specify supported languages to ensure successful processing and
guarantee output characteristics. Please only specify supported language codes:
[de, en, es, fr, ja]
```

Без явного указания языка модель либо не отвечает, либо ведёт себя непредсказуемо - сам Chrome
предупреждает об этом в консоли ещё до отказа.

Поскольку весь продукт и recap - на русском, честно работающей фичи не получилось: показывать
пользователю нерабочий виджет или переключать чат на английский в русском продукте - хуже, чем
не показывать ничего. Фичу просто выпилить решили.
