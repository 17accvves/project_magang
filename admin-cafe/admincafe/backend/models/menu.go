package models

import (
    "time"
)

type Menu struct {
    ID              string    `json:"id" db:"id"`
    Name            string    `json:"name" db:"name" binding:"required"`
    Price           float64   `json:"price" db:"price" binding:"required"`
    Discount        float64   `json:"discount" db:"discount"`
    DiscountedPrice float64   `json:"discountedPrice" db:"discounted_price"`
    StartDate       string    `json:"startDate" db:"start_date"`
    EndDate         string    `json:"endDate" db:"end_date"`
    Category        string    `json:"category" db:"category" binding:"required"`
    Status          string    `json:"status" db:"status"`
    Img             string    `json:"img" db:"img"`
    CreatedAt       time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}