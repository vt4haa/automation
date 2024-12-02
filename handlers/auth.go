package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"automation/db"
	"automation/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Получение данных пользователя
	var user models.User
	err = dbConn.QueryRow("SELECT workers.id, workers.login, workers.fio, positions.name as post, workers.pass FROM workers JOIN positions on workers.post = positions. id WHERE login = ?", creds.Login).Scan(&user.Id, &user.Login, &user.Fio, &user.Post, &user.Pass)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Получение названия должности

	// Проверка, является ли пароль хэшированным
	if !db.IsPasswordHashed(user.Pass) {
		hashedPassword, err := db.HashPassword(user.Pass)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		_, err = dbConn.Exec("UPDATE workers SET pass = ? WHERE login = ?", hashedPassword, user.Login)
		if err != nil {
			http.Error(w, "Database update error", http.StatusInternalServerError)
			return
		}
		user.Pass = hashedPassword
	}

	// Проверка пароля
	if !db.CheckPasswordHash(creds.Password, user.Pass) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Отправка данных пользователя в ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
