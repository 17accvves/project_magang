package handlers

import (
	"backend/models"
	"backend/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthHandler struct {
	repo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// ==============================
// LOGIN (Plain Text Password)
// ==============================
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }

    // <- letakkan di sini
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        fmt.Println("JSON decode error:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    fmt.Println("Login attempt:", body.Username, "Role:", body.Role)

    // query DB setelah ini
    var dbPassword, role, username string
    err = h.repo.DB.QueryRow(
        "SELECT username, password, role FROM users WHERE username=$1 AND role=$2",
        body.Username, body.Role,
    ).Scan(&username, &dbPassword, &role)

    if err != nil {
        fmt.Println("DB query error:", err)
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    if body.Password != dbPassword {
        fmt.Println("Password mismatch")
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    fmt.Println("Login success:", username, role)
    json.NewEncoder(w).Encode(map[string]string{
        "username": username,
        "role":     role,
    })
}


// ==============================
// REGISTER CAFE (Plain Password)
// ==============================
func (h *AuthHandler) RegisterCafe(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		fmt.Println("Parse form error:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// ===== Simpan file izin usaha =====
	file, header, err := r.FormFile("izin_usaha")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := "./uploads/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}

	// PASSWORD DISIMPAN TANPA HASH
	user := &models.User{
		Username:  username,
		Password:  password, // 🔥 plain text
		Email:     email,
		Role:      "cafe",
		IzinUsaha: filePath,
		Verified:  false,
	}

	err = h.repo.CreateCafe(user)
	if err != nil {
		fmt.Println("CreateCafe error:", err)
		http.Error(w, "Gagal registrasi", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registrasi sukses! Tunggu approval admin.",
	})
}
