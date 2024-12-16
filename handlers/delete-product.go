package handlers

import (
	"automation/db" // Подключение к базе данных
	"encoding/json"
	"log"
	"net/http"
)

// Структура для получения данных из запроса
type DeleteProductRequest struct {
	ID    int    `json:"id"`    // ID товара
	Name  string `json:"name"`  // Название товара
	Stock int    `json:"stock"` // Остаток на складе
}

// Удаление товара
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем подключение к базе данных
	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	// Декодируем JSON-запрос
	var req DeleteProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		log.Println("Error decoding request:", err)
		return
	}

	// Удаляем товар из базы данных
	query := `DELETE FROM product WHERE id = ? AND name = ? AND stock = ?`
	result, err := database.Exec(query, req.ID, req.Name, req.Stock)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		log.Println("Error deleting product:", err)
		return
	}

	// Проверяем, были ли удалены строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve delete result", http.StatusInternalServerError)
		log.Println("Error retrieving rows affected:", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No product found with the given criteria", http.StatusNotFound)
		log.Println("No matching product found")
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}
