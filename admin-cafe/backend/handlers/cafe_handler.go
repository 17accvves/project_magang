package handlers

import (
    "backend/repository"
    "encoding/json"
    "fmt"
    "net/http"
)

type CafeHandler struct {
    repo *repository.UserRepository
}

func NewCafeHandler(repo *repository.UserRepository) *CafeHandler {
    return &CafeHandler{repo: repo}
}

// ==============================
// Approve cafe by admin
// ==============================
func (h *CafeHandler) ApproveCafe(w http.ResponseWriter, r *http.Request) {
    var body struct {
        CafeID int `json:"cafe_id"`
    }

    // Decode JSON body
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        fmt.Println("JSON decode error in ApproveCafe:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    fmt.Println("Approve attempt for cafe ID:", body.CafeID)

    // Panggil repository untuk verify cafe
    err = h.repo.VerifyCafe(body.CafeID)
    if err != nil {
        fmt.Println("DB error approving cafe:", err)
        http.Error(w, "Gagal approve cafe", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Cafe approved!"})
}

// ==============================
// Reject cafe by admin (opsional)
// ==============================
func (h *CafeHandler) RejectCafe(w http.ResponseWriter, r *http.Request) {
    var body struct {
        CafeID int `json:"cafe_id"`
    }

    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        fmt.Println("JSON decode error in RejectCafe:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    fmt.Println("Reject attempt for cafe ID:", body.CafeID)

    // Hapus cafe dari DB (atau bisa update status)
    _, err = h.repo.DB.Exec("DELETE FROM users WHERE id=$1 AND role='cafe' AND verified=false", body.CafeID)
    if err != nil {
        fmt.Println("DB error rejecting cafe:", err)
        http.Error(w, "Gagal menolak cafe", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Cafe ditolak!"})
}

// ==============================
// List semua cafe (pending, approved, rejected)
// ==============================
func (h *CafeHandler) ListAllCafes(w http.ResponseWriter, r *http.Request) {
	type Cafe struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		IzinUsaha string `json:"izin_usaha"`
		Verified  bool   `json:"verified"`
		Rejected  bool   `json:"rejected"` // buat field ini kalau ada
	}

	rows, err := h.repo.DB.Query("SELECT id, username, email, izin_usaha, verified, rejected FROM users WHERE role='cafe'")
	if err != nil {
		fmt.Println("DB query error in ListAllCafes:", err)
		http.Error(w, "Gagal mengambil data cafe", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cafes []Cafe
	for rows.Next() {
		var c Cafe
		err := rows.Scan(&c.ID, &c.Username, &c.Email, &c.IzinUsaha, &c.Verified, &c.Rejected)
		if err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		cafes = append(cafes, c)
	}

	json.NewEncoder(w).Encode(cafes)
}

