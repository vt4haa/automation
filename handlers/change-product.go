package handlers

import (
	"automation/db"
	"encoding/json"
	"log"
	"net/http"
)

// Структура для изменения продукта
type ChangeProductRequest struct {
	ID          int    `json:"id"`           // ID продукта (обязательный)
	NewName     string `json:"new_name"`     // Новое название продукта
	NewPhoto    string `json:"new_photo"`    // Новое фото (URL)
	NewCategory int    `json:"new_category"` // Новый ID категории
	NewBrand    int    `json:"new_brand"`    // Новый ID бренда
	NewStock    int    `json:"new_stock"`    // Новый остаток
	NewPrice    int    `json:"new_price"`    // Новая цена
}

// Обработчик для изменения данных продукта
func ChangeProductHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Подключаемся к базе данных
	database, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer database.Close()

	// Декодируем тело запроса
	var req ChangeProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		log.Println("Error decoding request:", err)
		return
	}

	// Проверяем обязательное поле ID
	if req.ID == 0 {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		log.Println("Product ID is missing")
		return
	}

	// Обновляем данные продукта в базе данных
	query := `UPDATE product 
		SET name = ?, photo = ?, idCategories = ?, idBrands = ?, stock = ?, price = ? 
		WHERE id = ?`
	result, err := database.Exec(query, req.NewName, req.NewPhoto, req.NewCategory, req.NewBrand, req.NewStock, req.NewPrice, req.ID)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		log.Println("Error updating product:", err)
		return
	}

	// Проверяем, были ли обновлены строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve update result", http.StatusInternalServerError)
		log.Println("Error retrieving rows affected:", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No product found with the given ID", http.StatusNotFound)
		log.Println("No product found with ID:", req.ID)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}
