# Code Style Guide — SHRTIK

## Go

### Именование
| Тип        | Правило               | Пример                    |
| ---------- | --------------------- | ------------------------- |
| Пакет      | lowercase, одно слово | handler, service          |
| Файл       | snake_case            | link.go, link_test.go     |
| Структура  | PascalCase            | CreateLinkInput, LinkRepo |
| Интерфейс  | PascalCase + er       | LinkRepository            |
| Функция    | PascalCase (публ)     | CreateLink(), GetStats()  |
| Переменная | camelCase             | shortCode, clickCount     |
| Константа  | PascalCase            | MaxShortLen, CacheTTL     |

### Обработка ошибок

```go
// Всегда проверяем ошибки
link, err := repo.FindByCode(ctx, code)
if err != nil {
    return nil, fmt.Errorf("failed to find link: %w", err)
}

// Не игнорируем
link, _ := repo.FindByCode(ctx, code) // так нельзя
```

### Комментарии

```go
// GenerateCode создаёт случайный короткий код из 7 символов base62.
// Использует криптографически безопасный генератор случайных чисел.
func GenerateCode() (string, error) {
    // ...
}
```

### Форматирование
- Табы (go fmt стандарт)
- Максимальная длина строки — 120 символов
- Импорты: стандартная библиотека -> внешние -> внутренние
```go
import (
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"

    "shrtik/internal/model"
)
```

---

## SQL (Миграции)

### Именование
| Тип           | Правило                 | Пример                   |
|---------------|-------------------------|--------------------------|
| Файл миграции | NNN_description.sql     | 001_create_users.sql     |
| Таблица       | snake_case, множ. число | links, link_stats        |
| Поле          | snake_case, ед. число   | created_at, short_code   |
| Первичный ключ| id                      | id SERIAL PRIMARY KEY    |
| Внешний ключ  | {table}_id              | link_id                  |

### Пример миграции
```sql
-- 001_create_links.sql
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    long_url TEXT NOT NULL,
    clicks INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Git

### Ветки
- main — только рабочий, проверенный код
- feature/имя-фичи — разработка новой функциональности
- fix/что-чиним — исправление багов

### Коммиты
<префикс>: <описание на русском>

```
feat: реализована генерация short_code
fix: исправлен редирект при отсутствии кода в Redis
chore: обновлён compose.yml (Redis 7 -> 7-alpine)
docs: добавлен CODE_STYLE.md
refactor: вынесен LinkRepository в интерфейс
test: добавлены тесты для LinkService
```

---

## Соглашения по проекту

| Договорённость         | Значение               |
|------------------------|------------------------|
| Язык комментариев      | Русский                |
| Язык логов             | Русский                |
| Формат даты в БД       | TIMESTAMP              |
| Порт по умолчанию      | 8080                   |
| Точка входа            | cmd/server/main.go     |
| База данных            | PostgreSQL 16          |
| Кэш                    | Redis 7                |
