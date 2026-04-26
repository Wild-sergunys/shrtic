package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
	if os.Getenv("SKIP_MIGRATIONS") == "true" {
		log.Println("Миграции пропущены (SKIP_MIGRATIONS=true)")
		return nil
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("ошибка инициализации миграций: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	log.Println("Миграции применены успешно")
	return nil
}
