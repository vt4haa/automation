package main

import (
	"automation/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/add-product", handlers.AddProductHandler)
	http.HandleFunc("/purchase-history", handlers.PurchaseHistoryHandler)
	http.HandleFunc("/total-purchases", handlers.TotalPurchasesHandler)
	http.HandleFunc("/get-users", handlers.GetUsersHandler)
	http.HandleFunc("/get-all-products", handlers.GetAllProductsHandler)    // Обработчик для получения товаров
	http.HandleFunc("/get-product-image/", handlers.GetProductImageHandler) // Обработчик для получения изображения
	http.HandleFunc("/get-clients", handlers.GetClientsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
