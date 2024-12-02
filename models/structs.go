package models

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	Id    int    `json:"Id"`
	Login string `json:"login"`
	Fio   string `json:"fio"`
	Post  string `json:"post"`
	Pass  string `json:"-"`
}

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
type Client struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"` // Заменено на поле contact
}
