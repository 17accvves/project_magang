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

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("JSON decode error in ApproveCafe:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Approve attempt for cafe ID:", body.CafeID)

	// Update cafe menjadi verified=true, rejected=false
	err := h.repo.VerifyCafe(body.CafeID)
	if err != nil {
		fmt.Println("DB error approving cafe:", err)
		http.Error(w, "Gagal approve cafe", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Cafe approved!"})
}

// ==============================
// Reject cafe by admin
// ==============================
func (h *CafeHandler) RejectCafe(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CafeID int `json:"cafe_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("JSON decode error in RejectCafe:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Reject attempt for cafe ID:", body.CafeID)

	// Update cafe menjadi rejected=true, verified=false
	_, err := h.repo.DB.Exec("UPDATE users SET rejected=true, verified=false WHERE id=$1 AND role='cafe'", body.CafeID)
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
		Rejected  bool   `json:"rejected"`
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
		if err := rows.Scan(&c.ID, &c.Username, &c.Email, &c.IzinUsaha, &c.Verified, &c.Rejected); err != nil {
			fmt.Println("Row scan error:", err)
			continue
		}
		cafes = append(cafes, c)
	}

	json.NewEncoder(w).Encode(cafes)
}
