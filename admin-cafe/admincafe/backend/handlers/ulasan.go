package handlers

import (
    "net/http"
    "os"
    "path/filepath"

    "backend/models"
    "backend/repository"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type UlasanHandler struct {
    repo *repository.UlasanRepository
}

func NewUlasanHandler(repo *repository.UlasanRepository) *UlasanHandler {
    return &UlasanHandler{repo: repo}
}

// Get All Ulasan (for public)
func (h *UlasanHandler) GetUlasan(c *gin.Context) {
    ulasan, err := h.repo.GetAllUlasan()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ulasan)
}

// Get All Ulasan for Admin
func (h *UlasanHandler) GetUlasanAdmin(c *gin.Context) {
    ulasan, err := h.repo.GetAllUlasanForAdmin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ulasan)
}

// Get Ulasan Stats
func (h *UlasanHandler) GetUlasanStats(c *gin.Context) {
    stats, err := h.repo.GetUlasanStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, stats)
}

// Create Ulasan
func (h *UlasanHandler) CreateUlasan(c *gin.Context) {
    var ulasan models.Ulasan
    if err := c.ShouldBindJSON(&ulasan); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.repo.CreateUlasan(&ulasan)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":      ulasan.ID,
        "message": "Ulasan berhasil dikirim dan menunggu persetujuan",
    })
}

// Add Reply to Ulasan
func (h *UlasanHandler) AddReply(c *gin.Context) {
    id := c.Param("id")
    
    var replyReq models.ReplyRequest
    if err := c.ShouldBindJSON(&replyReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.repo.AddReply(id, replyReq.Reply)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Balasan berhasil ditambahkan"})
}

// Delete Reply
func (h *UlasanHandler) DeleteReply(c *gin.Context) {
    id := c.Param("id")
    
    err := h.repo.DeleteReply(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Balasan berhasil dihapus"})
}

// Update Ulasan Status (approve/reject)
func (h *UlasanHandler) UpdateUlasanStatus(c *gin.Context) {
    id := c.Param("id")
    
    var updateData struct {
        Status string `json:"status" binding:"required"`
    }
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.repo.UpdateUlasan(id, map[string]interface{}{
        "status": updateData.Status,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Status ulasan berhasil diupdate"})
}

// Delete Ulasan
func (h *UlasanHandler) DeleteUlasan(c *gin.Context) {
    id := c.Param("id")
    
    err := h.repo.DeleteUlasan(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ulasan berhasil dihapus"})
}

// Upload Gambar Ulasan
func (h *UlasanHandler) UploadGambarUlasan(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Create uploads directory if not exists
    if err := os.MkdirAll("uploads/ulasan", os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
        return
    }

    // Generate unique filename
    filename := uuid.New().String() + filepath.Ext(file.Filename)
    filepath := filepath.Join("uploads/ulasan", filename)

    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Return accessible URL
    url := "http://localhost:8080/uploads/ulasan/" + filename
    c.JSON(http.StatusOK, gin.H{"url": url})
}

// Upload Avatar Ulasan
func (h *UlasanHandler) UploadAvatarUlasan(c *gin.Context) {
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Create uploads directory if not exists
    if err := os.MkdirAll("uploads/avatar", os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
        return
    }

    // Generate unique filename
    filename := uuid.New().String() + filepath.Ext(file.Filename)
    filepath := filepath.Join("uploads/avatar", filename)

    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Return accessible URL
    url := "http://localhost:8080/uploads/avatar/" + filename
    c.JSON(http.StatusOK, gin.H{"url": url})
}