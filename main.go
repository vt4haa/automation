package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Credentials структура для приема данных из запроса
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// User структура для хранения данных пользователя из БД
type User struct {
	Login string `json:"login"`
	Fio   string `json:"fio"`
	Post  string `json:"post"` // Теперь здесь будет должность, извлеченная из positions
	Pass  string `json:"-"`    // Не возвращаем пароль в ответ
}

// Подключение к базе данных
func connectDB() (*sql.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/autocast" // Замените yourdbname на имя вашей базы данных
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Функция для проверки хеша пароля
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Функция для создания bcrypt-хэша пароля
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Обработчик авторизации
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Декодируем тело запроса в структуру Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Ищем пользователя в базе данных по логину
	var user User
	var postID int
	err = db.QueryRow("SELECT login, fio, post, pass FROM workers WHERE login = ?", creds.Login).Scan(&user.Login, &user.Fio, &postID, &user.Pass)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Получаем название должности из таблицы positions
	var postName string
	err = db.QueryRow("SELECT name FROM positions WHERE id = ?", postID).Scan(&postName)
	if err != nil {
		http.Error(w, "Database error while retrieving position", http.StatusInternalServerError)
		return
	}

	// Устанавливаем полученную должность в поле post
	user.Post = postName

	// Проверяем, является ли пароль хэшированным
	if !strings.HasPrefix(user.Pass, "$2a$") {
		fmt.Println("Password is not hashed; hashing and updating the database.")
		hashedPassword, err := hashPassword(user.Pass)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Обновляем пароль в базе данных на хэшированный
		_, err = db.Exec("UPDATE workers SET pass = ? WHERE login = ?", hashedPassword, user.Login)
		if err != nil {
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
		user.Pass = hashedPassword
	}

	// Проверяем пароль
	fmt.Println("Login from request:", creds.Login)
	fmt.Println("Password from request:", creds.Password)
	fmt.Println("Hashed password from DB:", user.Pass)

	if !checkPasswordHash(creds.Password, user.Pass) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Успешная авторизация: отправляем данные пользователя в JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
