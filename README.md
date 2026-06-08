# SHRTIC - СЕРВИС СОКРАЩЕНИЯ ССЫЛОК

[![Go Report Card](https://goreportcard.com/badge/github.com/Wild-sergunys/shrtic)](https://goreportcard.com/report/github.com/Wild-sergunys/shrtic)
[![Test](https://github.com/Wild-sergunys/shrtic/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Wild-sergunys/shrtic/actions/workflows/test.yml)

Веб-приложение для сокращения ссылок с аналитикой переходов. Генерация коротких кодов, редирект с кэшированием, сбор статистики по браузерам/устройствам/странам/источникам.

## БЫСТРЫЙ СТАРТ

```bash
git clone https://github.com/Wild-sergunys/shrtic.git
cd shrtic
cp .env.example .env
# Обязательно измените JWT_SECRET в .env!
docker compose up -d
```

Сервер запустится на `http://localhost:8080`.

## СТЕК

| Компонент | Технология |
|-----------|------------|
| Backend | Go, стандартный `net/http`, JWT-аутентификация |
| База данных | PostgreSQL 16 |
| Кэш | Redis 7 |
| Фронтенд | Чистый HTML/CSS/JS |
| Мониторинг | Prometheus + Grafana |
| Деплой | Docker, Docker Compose |

## ВОЗМОЖНОСТИ

- Сокращение ссылок: генерация коротких кодов (base62, 7 символов)
- Редирект с кэшированием в Redis (TTL 24h) - быстрый переход без запроса к БД
- Сбор статистики переходов: браузеры, устройства, страны, источники перехода
- Страна определяется по IP через ip-api.com (на localhost: "Локальный")
- JWT-авторизация через cookie или Bearer токен
- Личный кабинет со списком ссылок, поиском и статистикой
- Rate limiter на попытки входа (5 попыток за 15 минут, блокировка 15 минут)
- Graceful shutdown сервера
- Автоматический прогон миграций при старте
- Prometheus метрики (RPS, latency, активные ссылки, пользователи)

## АУТЕНТИФИКАЦИЯ

После успешного входа JWT токен устанавливается в cookie `shrtic_token` (HttpOnly, 24 часа).
Также токен принимается в заголовке `Authorization: Bearer <token>`.

Ответ при входе:
```json
{
  "role": "user"
}
```

## ДОКУМЕНТАЦИЯ API

- [OpenAPI спецификация](api/openapi.yaml)
- [Форматы запросов и ошибок](api/FORMATS.md)
- [Сбор статистики](api/STATS.md)
- [Визуализация данных](api/VISUALIZATION.md)

## АРХИТЕКТУРА

```
.
├── api/                    # Документация API
├── cmd/server/             # Точка входа
├── internal/
│   ├── config/             # Загрузка конфигурации из .env
│   ├── database/           # PostgreSQL, Redis, миграции
│   ├── handler/            # HTTP-обработчики
│   ├── middleware/         # JWT, rate limiter, метрики
│   ├── model/              # Структуры данных
│   ├── repository/         # Доступ к БД
│   └── service/            # Бизнес-логика
├── migrations/             # SQL-миграции
├── web/                    # Фронтенд
│   ├── pages/              # HTML страницы
│   └── static/             # CSS, JS
├── prometheus.yml
├── shrtic-dashboard.json
└── docker-compose.yml
```

Слои: `handler → service → repository`.

## СТАТИСТИКА ПЕРЕХОДОВ

При каждом переходе по короткой ссылке **синхронно** собирается:

- **Браузер:** Chrome, Firefox, Safari, Edge, Other (парсинг User-Agent)
- **Устройство:** Desktop, Mobile, Tablet (парсинг User-Agent)
- **Страна:** определяется по IP через ip-api.com (на localhost: "Локальный")
- **Источник:** Прямой, Twitter, Telegram, Facebook, Google, Yandex, Other (заголовок Referer)

Статистика доступна в личном кабинете при раскрытии карточки ссылки.

## МОНИТОРИНГ

- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

Импорт дашборда:
```bash
curl -X POST http://localhost:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d @shrtic-dashboard.json
```

## ТЕСТЫ

```bash
go test ./internal/... -v
```

## ЛИЦЕНЗИЯ

MIT

