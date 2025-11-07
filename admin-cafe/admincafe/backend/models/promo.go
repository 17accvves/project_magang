package models

import (
	"time"
)

type Promo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	StartDate string    `json:"start"` // Format: "18 September 2025"
	EndDate   string    `json:"end"`   // Format: "25 September 2025"
	Revenue   string    `json:"revenue"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PromoStats struct {
	ActivePromos int      `json:"active_promos"`
	TotalRevenue string   `json:"total_revenue"`
	PromosList   []string `json:"promos_list"`
}

type PromoRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Revenue   string `json:"revenue" binding:"required"`
	Status    string `json:"status"`
}

type PromoPurchase struct {
	ID           int       `json:"id"`
	PromoID      int       `json:"promo_id"`
	MenuID       string    `json:"menu_id"`
	Quantity     int       `json:"quantity"`
	TotalAmount  float64   `json:"total_amount"`
	PurchaseDate time.Time `json:"purchase_date"`
}

type PromoPurchaseRequest struct {
	PromoID     int     `json:"promo_id" binding:"required"`
	MenuID      string  `json:"menu_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
}

type PromoWithMenus struct {
	Promo Promo  `json:"promo"`
	Menus []Menu `json:"menus"`
}

type PromoRevenue struct {
	PromoID    int     `json:"promo_id"`
	PromoTitle string  `json:"promo_title"`
	Revenue    float64 `json:"revenue"`
	TotalSales int     `json:"total_sales"`
}