package models

import (
    "time"
)

type Menu struct {
    ID              string  `json:"id"`
    Name            string  `json:"name"`
    Price           float64 `json:"price"`
    Discount        float64 `json:"discount"`
    DiscountedPrice float64 `json:"discounted_price"`
    StartDate       string  `json:"start_date"`
    EndDate         string  `json:"end_date"`
    Category        string  `json:"category"`
    Status          string  `json:"status"`
    Img             string  `json:"img"`
    CafeID          int     `json:"cafe_id"`      // ✅ TAMBAHAN
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}