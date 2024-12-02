package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"automation/db"
	"automation/models"
)

// Функция для получения всех работников
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Подключение к базе данных
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Запрос на получение всех пользователей
	rows, err := dbConn.Query(`
		SELECT w.login, w.fio, p.name AS post, w.pass
		FROM workers w
		LEFT JOIN positions p ON w.post = p.id
	`)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Массив для хранения пользователей
	var users []models.User

	// Обработка результатов запроса
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Login, &user.Fio, &user.Post, &user.Pass)
		if err != nil {
			log.Println("Error scanning user:", err)
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			return
		}

		// Исключаем передачу пароля
		user.Pass = "" // Очищаем пароль перед отправкой клиенту
		users = append(users, user)
	}

	// Проверка на наличие записей
	if len(users) == 0 {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	// Отправка данных пользователей в JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println("Error encoding users:", err)
		http.Error(w, "Error encoding users data", http.StatusInternalServerError)
		return
	}
}
func GetClientsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Получаем список всех клиентов
	rows, err := dbConn.Query("SELECT id, name, contact FROM clients")
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.ID, &client.Name, &client.Contact) // Используем contact вместо email
		if err != nil {
			http.Error(w, "Error scanning client data", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	// Проверяем наличие клиентов
	if len(clients) == 0 {
		http.Error(w, "No clients found", http.StatusNotFound)
		return
	}

	// Отправляем список клиентов в JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		http.Error(w, "Error encoding clients data", http.StatusInternalServerError)
		return
	}
}
