package models

import "time"

type CafeProfile struct {
	ID          int       `json:"id"`
	Nama        string    `json:"nama"`
	Alamat      string    `json:"alamat"`
	Telepon     string    `json:"telepon"`
	Deskripsi   string    `json:"deskripsi"`
	MainImage   string    `json:"main_image"`
	Verified    bool      `json:"verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CafeSocialMedia struct {
	ID            int       `json:"id"`
	CafeProfileID int       `json:"cafe_profile_id"`
	Platform      string    `json:"platform"`
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
}

type CafeOperationalHours struct {
	ID            int       `json:"id"`
	CafeProfileID int       `json:"cafe_profile_id"`
	Hari          string    `json:"hari"`
	Buka          string    `json:"buka"`   // Format: "07:00" (24 jam)
	Tutup         string    `json:"tutup"`  // Format: "22:00" (24 jam)
	CreatedAt     time.Time `json:"created_at"`
}

type CafeFacility struct {
	ID            int       `json:"id"`
	CafeProfileID int       `json:"cafe_profile_id"`
	NamaFasilitas string    `json:"nama_fasilitas"`
	Tersedia      bool      `json:"tersedia"`
	CreatedAt     time.Time `json:"created_at"`
}

type CafeGallery struct {
	ID            int       `json:"id"`
	CafeProfileID int       `json:"cafe_profile_id"`
	ImageURL      string    `json:"image_url"`
	Urutan        int       `json:"urutan"`
	CreatedAt     time.Time `json:"created_at"`
}

// ✅ MODEL BARU: OpeningHoursResponse untuk response API format 24 jam
type OpeningHoursResponse struct {
	Waktu  string `json:"waktu"`  // Format: "07:00 - 22:00"
	Status string `json:"status"` // "Buka Sekarang" atau "Tutup Sekarang"
}

// ✅ MODEL BARU: UpdateOpeningHoursRequest untuk update jam operasional
type UpdateOpeningHoursRequest struct {
	OpenTime  string `json:"open_time"`  // Format: "07:00"
	CloseTime string `json:"close_time"` // Format: "22:00"
}

// ✅ MODEL BARU: CurrentStatusResponse untuk response status saja
type CurrentStatusResponse struct {
	Status string `json:"status"` // "Buka Sekarang" atau "Tutup Sekarang"
}