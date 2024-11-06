// main.go
package main

import (
	"automation/db"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go" // Импортируем пакет jwt-go
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var secretKey = []byte("your-secret-key")

// Структура для получения логина и пароля от клиента
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

// Генерация JWT токена
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // токен будет действовать 24 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Проверка JWT токена
func validateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// Проверка пароля с хешем
func checkPasswordHash(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Получаем пользователя из базы данных
	user, err := db.GetUserByLogin(creds.Login)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Println("User found in DB:", user.Login)

	// Проверка пароля
	if !checkPasswordHash(creds.Password, user.Password) {
		fmt.Println("Password check failed")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Println("Password check passed")

	// Генерация токена
	token, err := generateJWT(user.ID)
	if err != nil {
		fmt.Println("Error generating token:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	fmt.Println("Generated token:", token)

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Инициализация подключения к базе данных
	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Обработчики маршрутов
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
