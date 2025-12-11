package handlers

import (
	"backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

// Constructor
func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
	}
}

// =========================
// Login handler dengan role check
// =========================
func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"` // tangkap role dari frontend
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, password, dan role wajib diisi"})
		return
	}

	user, err := h.userRepo.Login(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah!"})
		return
	}

	// ✅ Validasi role
	if user.Role != body.Role {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak sesuai dengan akun"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
