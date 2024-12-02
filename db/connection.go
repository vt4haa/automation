package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Импортируем драйвер MySQL
)

// ConnectDB создает подключение к базе данных и возвращает его
func ConnectDB() (*sql.DB, error) {
	// Замените параметры подключения на свои
	dsn := "root:@tcp(127.0.0.1:3306)/autocast"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	return db, nil
}
