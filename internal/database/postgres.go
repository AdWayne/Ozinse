package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось достучаться до БД: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("ошибка миграции: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_name = 'users'
	)`).Scan(&exists)
	if err != nil {
		return err
	}

	// Таблицы уже есть — только наполняем справочники
	if exists {
		migration, err := os.ReadFile("migrations/002_insert.sql")
		if err != nil {
			return fmt.Errorf("ошибка чтения 002_insert.sql: %w", err)
		}
		_, err = db.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("ошибка выполнения 002_insert.sql: %w", err)
		}
		return nil
	}

	// Первый запуск — создаём таблицы и наполняем данными
	files := []string{"migrations/001_init.sql", "migrations/002_insert.sql"}
	for _, file := range files {
		migration, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("ошибка чтения %s: %w", file, err)
		}
		_, err = db.Exec(string(migration))
		if err != nil {
			return fmt.Errorf("ошибка выполнения %s: %w", file, err)
		}
	}

	return nil
}