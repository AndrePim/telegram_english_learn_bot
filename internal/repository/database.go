package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database представляет собой структуру для работы с базой данных
type Database struct {
	db *sql.DB
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{db: db}

	// Создаем таблицы при инициализации
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Successfully connected to database")
	return database, nil
}

// createTables создает необходимые таблицы
func (d *Database) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(255),
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			state VARCHAR(50) DEFAULT 'idle',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS words (
			id SERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id),
			word VARCHAR(255) NOT NULL,
			translation VARCHAR(255) NOT NULL,
			context TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_review TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			next_review TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			interval INTEGER DEFAULT 1,
			difficulty INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS quizzes (
			id SERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id),
			word_id INTEGER REFERENCES words(id),
			correct BOOLEAN NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := d.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
	return d.db.Close()
}
