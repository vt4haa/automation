package handlers

import (
	"automation/db"
	"automation/models"
	"encoding/json"
	"net/http"
)

// Handler для получения истории покупок пользователя
func PurchaseHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	rows, err := dbConn.Query(`
		SELECT pr.name, p.quantity, p.purchase_date 
		FROM purchases p
		JOIN product pr ON p.product_id = pr.id
		WHERE p.user_id = ?`, userID)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		var purchaseDate string

		err := rows.Scan(&purchase.ProductName, &purchase.Quantity, &purchaseDate)
		if err != nil {
			http.Error(w, "Error scanning purchase data", http.StatusInternalServerError)
			return
		}

		purchase.PurchaseDate = purchaseDate
		purchases = append(purchases, purchase)
	}

	if len(purchases) == 0 {
		http.Error(w, "No purchases found for the user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(purchases)
	if err != nil {
		http.Error(w, "Error encoding purchase data", http.StatusInternalServerError)
		return
	}
}

// Handler для получения общего количества покупок пользователя
func TotalPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var totalPurchases int
	err = dbConn.QueryRow("SELECT COUNT(*) FROM purchases WHERE user_id = ?", userID).Scan(&totalPurchases)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"total_purchases": totalPurchases}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding purchase count", http.StatusInternalServerError)
		return
	}
}
