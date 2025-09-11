package helpers

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(dsn string) {
	m, err := migrate.New(
		"file://migrations", // путь до миграций
		dsn,
	)
	if err != nil {
		log.Fatal("migration init error:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed:", err)
	}
	log.Println("migrations applied")
}
