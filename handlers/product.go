package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// AddProductHandler обрабатывает загрузку файла
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем только метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Ограничиваем размер данных формы
	err := r.ParseMultipartForm(32 << 20) // Ограничение на 32 MB
	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
		log.Printf("Ошибка разбора формы: %v\n", err)
		return
	}

	// Получаем файл из формы
	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Ошибка получения файла. Проверьте ключ 'photo'", http.StatusBadRequest)
		log.Printf("Ошибка получения файла: %v\n", err)
		return
	}
	defer file.Close()

	// Проверяем наличие имени файла
	if fileHeader.Filename == "" {
		http.Error(w, "Имя файла не указано", http.StatusBadRequest)
		log.Println("Имя файла отсутствует в запросе")
		return
	}

	// Создаем папку для загрузки, если её нет
	uploadDir := "./uploads"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Ошибка создания директории для загрузки", http.StatusInternalServerError)
		log.Printf("Ошибка создания директории: %v\n", err)
		return
	}

	// Путь для сохранения файла
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	// Проверяем, не существует ли файл уже
	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "Файл с таким именем уже существует", http.StatusConflict)
		log.Printf("Файл уже существует: %s\n", filePath)
		return
	}

	// Создаем файл для записи
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		log.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer outFile.Close()

	// Копируем содержимое загруженного файла в файл на сервере
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Ошибка копирования содержимого файла", http.StatusInternalServerError)
		log.Printf("Ошибка копирования файла: %v\n", err)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	log.Printf("Файл успешно загружен: %s\n", filePath)
	fmt.Fprintf(w, "Файл успешно загружен: %s", filePath)
}
