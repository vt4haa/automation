package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// Purchase структура для хранения данных о покупке
type Purchase struct {
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	PurchaseDate string `json:"purchase_date"`
}
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Photo string  `json:"photo"`
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

// Обработчик для получения истории покупок
func purchaseHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
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

	// Получаем историю покупок пользователя
	rows, err := db.Query(`
		SELECT pr.name, p.quantity, p.purchase_date 
		FROM purchases p
		JOIN product pr ON p.product_id = pr.id
		WHERE p.user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var purchase Purchase
		var purchaseDate string

		err := rows.Scan(&purchase.ProductName, &purchase.Quantity, &purchaseDate)
		if err != nil {
			http.Error(w, "Error scanning purchase data", http.StatusInternalServerError)
			return
		}

		// Преобразуем строку даты в нужный формат
		purchase.PurchaseDate = purchaseDate
		purchases = append(purchases, purchase)
	}

	// Проверяем наличие покупок
	if len(purchases) == 0 {
		http.Error(w, "No purchases found for the user", http.StatusNotFound)
		return
	}

	// Отправляем историю покупок в JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(purchases)
	if err != nil {
		http.Error(w, "Error encoding purchase data", http.StatusInternalServerError)
		return
	}
}

// Обработчик для получения общего количества покупок
func totalPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
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

	// Получаем количество покупок пользователя
	var totalPurchases int
	err = db.QueryRow("SELECT COUNT(*) FROM purchases WHERE user_id = ?", userID).Scan(&totalPurchases)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	// Отправляем количество покупок в JSON
	response := map[string]int{"total_purchases": totalPurchases}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding purchase count", http.StatusInternalServerError)
		return
	}
}

// Обработчик для добавления товара с изображением
// Обработчик для добавления товара с изображением
func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму для получения файла
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	log.Println("Form parsed successfully")

	// Получаем информацию о товаре
	name := r.FormValue("name")
	price := r.FormValue("price")

	// Получаем файл изображения
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println("Error retrieving image:", err)
		http.Error(w, "Error retrieving image", http.StatusBadRequest)
		return
	}
	log.Println("File retrieved successfully")
	defer file.Close()

	// Создаем папку для изображений, если она не существует
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating uploads directory", http.StatusInternalServerError)
		return
	}

	// Сохраняем изображение на сервере
	imagePath := fmt.Sprintf("./uploads/%s.jpg", name) // Путь для сохранения файла
	out, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Error saving image", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Копируем содержимое файла в новое место
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error copying image", http.StatusInternalServerError)
		return
	}

	// Добавляем товар в базу данных
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO product (name, price, photo) VALUES (?, ?, ?)", name, price, imagePath)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product added successfully"))
}

// Обработчик для получения информации о товаре
func getProductHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID товара из URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
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

	// Получаем информацию о товаре
	var product Product
	err = db.QueryRow("SELECT id, name, price, photo FROM product WHERE id = ?", id).
		Scan(&product.ID, &product.Name, &product.Price, &product.Photo)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Отправляем данные товара в JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, "Error encoding product data", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/purchase-history", purchaseHistoryHandler)
	http.HandleFunc("/total-purchases", totalPurchasesHandler)
	http.HandleFunc("/add-product", addProductHandler)
	http.HandleFunc("/get-product", getProductHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
