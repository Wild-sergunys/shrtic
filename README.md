# SHRTIC - СЕРВИС СОКРАЩЕНИЯ ССЫЛОК

[![Go Report Card](https://goreportcard.com/badge/github.com/Wild-sergunys/shrtic)](https://goreportcard.com/report/github.com/Wild-sergunys/shrtic)
[![Test](https://github.com/Wild-sergunys/shrtic/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Wild-sergunys/shrtic/actions/workflows/test.yml)

Веб-приложение для сокращения ссылок с аналитикой переходов. Генерация коротких кодов, редирект с кэшированием, сбор статистики по браузерам/устройствам/странам/источникам.

## БЫСТРЫЙ СТАРТ

```bash
git clone https://github.com/Wild-sergunys/shrtic.git
cd shrtic
cp .env.example .env
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
| Деплой | Docker, Docker Compose |

## ВОЗМОЖНОСТИ

- Сокращение ссылок: генерация коротких кодов (base62, 7 символов)
- Редирект с кэшированием в Redis (TTL 24h) - быстрый переход без запроса к БД
- Сбор статистики переходов: браузеры, устройства, страны, источники перехода
- Страна определяется по IP через ip-api.com
- Optional JWT-авторизация: незалогиненные создают ссылки, залогиненные видят статистику
- Личный кабинет со списком ссылок, поиском и удалением
- Graceful shutdown сервера
- Автоматический прогон миграций при старте
- Rate limiter на попытки входа (защита от брутфорса)

## ДОКУМЕНТАЦИЯ API

- [OpenAPI спецификация](api/openapi.yaml)

## АРХИТЕКТУРА

```
.
├── api/                    # OpenAPI спецификация
├── cmd/server/             # Точка входа
├── internal/
│   ├── config/             # Загрузка конфигурации из .env
│   ├── database/           # Подключение к PostgreSQL и Redis, миграции
│   ├── handler/            # HTTP-обработчики (auth, links, redirect)
│   ├── middleware/         # JWT-авторизация (обычная + опциональная), rate limiter
│   ├── model/              # Структуры данных
│   ├── repository/         # Доступ к БД (PostgreSQL)
│   └── service/            # Бизнес-логика (auth, ссылки, статистика)
├── migrations/             # SQL-миграции (PostgreSQL)
└── web/                    # Фронтенд (HTML/CSS/JS)
    ├── pages/              # HTML-страницы
    └── static/             # CSS, JavaScript
```

Слои: `handler → service → repository`. Handler принимает HTTP-запросы, service содержит бизнес-логику, repository работает с БД.

## СТАТИСТИКА ПЕРЕХОДОВ

При каждом переходе по короткой ссылке автоматически собирается:

- **Браузер:** Chrome, Firefox, Safari, Other (парсинг User-Agent)
- **Устройство:** Desktop, Mobile, Tablet (парсинг User-Agent)
- **Страна:** определяется по IP через ip-api.com (на localhost: `localhost`)
- **Источник:** Прямой, Twitter, Telegram, Facebook, Google и др. (заголовок Referer)

Статистика доступна в личном кабинете при раскрытии карточки ссылки.

## ТЕСТЫ

Юнит-тесты покрывают:
TODO: НАПИСАТЬ ТЕСТЫ

```bash
go test ./internal/... ./web/ -v
```

## ЛИЦЕНЗИЯ

MIT
