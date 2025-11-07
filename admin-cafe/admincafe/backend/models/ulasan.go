package models

import (
    "time"
)

type Ulasan struct {
    ID        string    `json:"id" db:"id"`
    Nama      string    `json:"nama" db:"nama" binding:"required"`
    Email     string    `json:"email" db:"email"`
    Rating    int       `json:"rating" db:"rating" binding:"required,min=1,max=5"`
    Teks      string    `json:"teks" db:"teks" binding:"required"`
    Gambar    string    `json:"gambar" db:"gambar"`
    Avatar    string    `json:"avatar" db:"avatar"`
    Balasan   string    `json:"balasan" db:"balasan"`
    Status    string    `json:"status" db:"status"` // "pending", "approved", "rejected"
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UlasanResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Date      string    `json:"date"`
    Rating    int       `json:"rating"`
    Text      string    `json:"text"`
    Image     string    `json:"image"`
    Avatar    string    `json:"avatar"`
    Reply     string    `json:"reply"`
}

type ReplyRequest struct {
    Reply string `json:"reply" binding:"required"`
}