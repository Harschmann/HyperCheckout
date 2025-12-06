package models

import (
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
