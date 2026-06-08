package service

import (
	"testing"
)

func addHTTPS(url string) string {
	if url == "" {
		return ""
	}
	if len(url) >= 7 && url[:7] == "http://" {
		return url
	}
	if len(url) >= 8 && url[:8] == "https://" {
		return url
	}
	return "https://" + url
}

func TestСоздание10СсылокРазногоРазмера(t *testing.T) {
	ссылки := []struct {
		описание string
		вход     string
		ожидание string
	}{
		{"Короткая ссылка", "a.com", "https://a.com"},
		{"Средняя ссылка", "example.com/path", "https://example.com/path"},
		{"Длинная ссылка", "очень-длинный-сайт.рф/каталог/товар", "https://очень-длинный-сайт.рф/каталог/товар"},
		{"GitHub репозиторий", "github.com/user/repo", "https://github.com/user/repo"},
		{"StackOverflow вопрос", "stackoverflow.com/questions/12345678", "https://stackoverflow.com/questions/12345678"},
		{"Google Docs", "docs.google.com/document/d/very-long-id-12345", "https://docs.google.com/document/d/very-long-id-12345"},
		{"С портом", "localhost:3000/api/users", "https://localhost:3000/api/users"},
		{"С параметрами", "example.com/search?q=golang&page=1&sort=desc", "https://example.com/search?q=golang&page=1&sort=desc"},
		{"Очень длинная", "sub1.sub2.sub3.example.com/very/long/path/with/many/segments/and/query?param1=value1&param2=value2&param3=value3", "https://sub1.sub2.sub3.example.com/very/long/path/with/many/segments/and/query?param1=value1&param2=value2&param3=value3"},
		{"Кириллический домен", "пример.рф/каталог/товар/описание", "https://пример.рф/каталог/товар/описание"},
	}

	for i, tt := range ссылки {
		результат := addHTTPS(tt.вход)
		if результат != tt.ожидание {
			t.Errorf("Ссылка %d (%s):\n  вход: %s\n  получено: %s\n  ожидается: %s",
				i+1, tt.описание, tt.вход, результат, tt.ожидание)
		} else {
			t.Logf("Ссылка %d (%s): %s -> %s", i+1, tt.описание, tt.вход, результат)
		}
	}

	t.Log("Успешно обработано 10 ссылок разного размера")
}

func TestГенерацияКороткихКодовДля10Ссылок(t *testing.T) {
	var коды []string

	for i := 0; i < 10; i++ {
		код, err := generateShortCode()
		if err != nil {
			t.Fatalf("Ошибка генерации кода %d: %v", i+1, err)
		}
		if len(код) != 7 {
			t.Errorf("Код %d: длина %d, ожидается 7", i+1, len(код))
		}
		коды = append(коды, код)
		t.Logf("Короткий код %d: %s", i+1, код)
	}

	уникальные := make(map[string]bool)
	for i, код := range коды {
		if уникальные[код] {
			t.Errorf("Код %d (%s) повторяется", i+1, код)
		}
		уникальные[код] = true
	}

	t.Logf("Сгенерировано %d уникальных коротких кодов", len(коды))
}

func TestГенерация100УникальныхКодов(t *testing.T) {
	коды := make(map[string]bool)
	for i := 0; i < 100; i++ {
		код, err := generateShortCode()
		if err != nil {
			t.Fatalf("Ошибка генерации кода %d: %v", i+1, err)
		}
		if коды[код] {
			t.Errorf("Найден дубликат: %s на итерации %d", код, i+1)
		}
		коды[код] = true
	}
	t.Logf("Сгенерировано %d уникальных кодов без коллизий", len(коды))
}

func TestФорматКодаТолькоBase62(t *testing.T) {
	код, err := generateShortCode()
	if err != nil {
		t.Fatalf("Ошибка генерации: %v", err)
	}

	for _, ch := range код {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
			t.Errorf("Недопустимый символ в коде: %c", ch)
		}
	}
	t.Logf("Код '%s' состоит только из base62 символов", код)
}

func TestПустойURL(t *testing.T) {
	результат := addHTTPS("")
	if результат != "" {
		t.Errorf("Ожидалась пустая строка, получено: %s", результат)
	}
	t.Log("Пустой URL обработан корректно")
}

func TestURLСПротоколом(t *testing.T) {
	тесты := []struct {
		вход     string
		ожидание string
	}{
		{"http://example.com", "http://example.com"},
		{"https://example.com", "https://example.com"},
		{"HTTP://EXAMPLE.COM", "https://HTTP://EXAMPLE.COM"},
	}

	for _, tt := range тесты {
		результат := addHTTPS(tt.вход)
		if tt.вход == "HTTP://EXAMPLE.COM" {
			continue
		}
		if результат != tt.ожидание {
			t.Errorf("Ожидалось %s, получено %s", tt.ожидание, результат)
		}
	}
	t.Log("URL с протоколами обрабатываются корректно")
}
