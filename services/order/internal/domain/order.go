package domain

import "time"

// ===== ORDER =====
type Order struct {
	ID           int       `json:"id"`
	ProductsCart []int     `json:"productsCart"`
	Price        int       `json:"price"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AddProductToCartRequest struct {
	UserID    int `json:"userID"`
	ProductID int `json:"productID"`
}
