// db/db.go
package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"log"
)

var db *sql.DB

// Инициализация подключения к базе данных MySQL
func InitDB() (*sql.DB, error) {
	// Строка подключения к базе данных
	connStr := "root@tcp(127.0.0.1:3306)/autocast"
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}
	log.Println("Successfully connected to the MySQL database")
	return db, nil
}

// Получение пользователя по логину
func GetUserByLogin(login string) (*User, error) {
	var user User
	query := `SELECT id, login, pass FROM workers WHERE login = ?`
	row := db.QueryRow(query, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}
	return &user, nil
}

// Структура для пользователя
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"pass"`
}
