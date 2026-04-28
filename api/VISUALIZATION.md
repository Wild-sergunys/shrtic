# Данные для визуализации SHRTIC

## Обзор

Бэкенд возвращает агрегированную статистику переходов по каждой ссылке. Фронтенд отображает её через шахматные полоски с процентами.

## Структура ответа /api/links/{id}/stats

{
  "total_clicks": 1247,
  "browsers": [
    {"name": "Chrome", "percentage": 62.5, "count": 780},
    {"name": "Firefox", "percentage": 18.0, "count": 225},
    {"name": "Safari", "percentage": 12.0, "count": 150},
    {"name": "Other", "percentage": 7.5, "count": 92}
  ],
  "devices": [
    {"name": "Desktop", "percentage": 71.0, "count": 885},
    {"name": "Mobile", "percentage": 24.0, "count": 299},
    {"name": "Tablet", "percentage": 5.0, "count": 63}
  ],
  "countries": [
    {"name": "Россия", "percentage": 45.0, "count": 561},
    {"name": "США", "percentage": 28.0, "count": 349},
    {"name": "Германия", "percentage": 15.0, "count": 187},
    {"name": "Other", "percentage": 12.0, "count": 150}
  ],
  "referrers": [
    {"name": "Прямой", "percentage": 45.0, "count": 561},
    {"name": "Twitter", "percentage": 28.0, "count": 349},
    {"name": "Telegram", "percentage": 22.0, "count": 274},
    {"name": "Other", "percentage": 5.0, "count": 63}
  ]
}

## Описание полей

### Корневые поля

| Поле | Тип | Описание |
|------|-----|----------|
| total_clicks | integer | Общее количество переходов по ссылке |

### StatItem (элементы массивов)

| Поле | Тип | Описание |
|------|-----|----------|
| name | string | Название (Chrome, Desktop, Россия, Twitter) |
| percentage | number | Процент от общего числа переходов |
| count | integer | Абсолютное количество переходов |

## Визуализация

### Шахматные полоски

В кабинете статистика отображается через шахматные полоски (chess-bar):
- Каждый блок - StatItem с шириной равной percentage
- Тёмные и светлые блоки чередуются
- Под полоской - названия с процентами

### Табличный вывод

На странице кабинета ссылки отображаются в виде раскрывающихся карточек:
| short_url | long_url | clicks | created_at |
|-----------|----------|--------|------------|
| /r/x7kQ2 | github.com/user/repo | 1247 | 26.04.2026 |

При раскрытии - статистика по браузерам, устройствам, странам, источникам.

## Источник данных

- Браузер/устройство: парсинг User-Agent
- Страна: IP-адрес через ip-api.com (на localhost: "localhost")
- Источник: заголовок Referer