package models

import (
    "time"
)

type Menu struct {
    ID             string    `json:"id"`
    Name           string    `json:"name"`
    Price          float64   `json:"price"`
    Discount       float64   `json:"discount"`
    DiscountedPrice float64   `json:"discountedPrice"`
    StartDate      string    `json:"startDate"`
    EndDate        string    `json:"endDate"`
    Category       string    `json:"category"`
    Status         string    `json:"status"`
    Img            string    `json:"img"`
    CreatedAt      time.Time `json:"createdAt"`
    UpdatedAt      time.Time `json:"updatedAt"`
}