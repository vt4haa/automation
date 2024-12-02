package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим форму
	err := r.ParseMultipartForm(10 << 20) // Ограничение на размер до 10 MB
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	log.Println("Form parsed successfully")

	// Получаем файл
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Println("Error retrieving file:", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader == nil {
		log.Println("No file header found")
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	// Сохраняем файл
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if mkdirErr := os.Mkdir(uploadDir, os.ModePerm); mkdirErr != nil {
			log.Println("Error creating upload directory:", mkdirErr)
			http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
			return
		}
	}

	filePath := filepath.Join(uploadDir, fileHeader.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error saving file:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = file.Seek(0, 0) // Возвращаемся к началу файла
	if err != nil {
		log.Println("Error seeking file:", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	if _, err := dst.ReadFrom(file); err != nil {
		log.Println("Error writing to file:", err)
		http.Error(w, "Error writing file to disk", http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded successfully: %s", filePath)

	// Ответ на запрос
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s", filePath)
}
