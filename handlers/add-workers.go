package handlers

import (
	"automation/db"     // Импортируем пакет db для подключения к базе данных
	"automation/models" // Импортируем модели для работы с данными
	"encoding/json"
	"log"
	"net/http"
)

// Добавление работника
func AddWorkerHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем подключение к базе данных
	database, err := db.ConnectDB() // Вызов функции ConnectDB из пакета db
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование данных работника из тела запроса
	var worker models.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Вставка данных работника в базу данных
	query := `INSERT INTO workers (fio, post, login, pass) VALUES (?, ?, ?, ?)`
	_, err = database.Exec(query, worker.Name, worker.Position, worker.Login, worker.Password)
	if err != nil {
		http.Error(w, "Failed to add worker to database", http.StatusInternalServerError)
		log.Println("Error adding worker:", err)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Worker added successfully"))
}
