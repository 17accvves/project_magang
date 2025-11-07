package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
	"backend/models"
	"backend/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CafeProfileHandler struct {
	repo *repository.CafeProfileRepository
}

func NewCafeProfileHandler(repo *repository.CafeProfileRepository) *CafeProfileHandler {
	return &CafeProfileHandler{repo: repo}
}

// GetCafeProfile - Get main cafe profile
func (h *CafeProfileHandler) GetCafeProfile(c *gin.Context) {
	profile, err := h.repo.GetCafeProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Get related data
	socialMedia, _ := h.repo.GetSocialMedia()
	operationalHours, _ := h.repo.GetOperationalHours()
	facilities, _ := h.repo.GetFacilities()
	
	response := gin.H{
		"id": profile.ID,
		"nama": profile.Nama,
		"alamat": profile.Alamat,
		"telepon": profile.Telepon,
		"deskripsi": profile.Deskripsi,
		"main_image": profile.MainImage,
		"verified": profile.Verified,
		"social_media": socialMedia,
		"operational_hours": operationalHours,
		"facilities": facilities,
	}
	
	c.JSON(http.StatusOK, response)
}

// UpdateCafeProfile - Update main cafe profile
func (h *CafeProfileHandler) UpdateCafeProfile(c *gin.Context) {
	var profile models.CafeProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.UpdateCafeProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// UploadProfileImage - Upload profile image (FIXED - RETURN FULL URL)
func (h *CafeProfileHandler) UploadProfileImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided: " + err.Error()})
		return
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory: " + err.Error()})
		return
	}

	// Generate unique filename to avoid conflicts
	fileExt := filepath.Ext(file.Filename)
	fileName := "profile_" + time.Now().Format("20060102150405") + fileExt
	filePath := "uploads/" + fileName

	// Save file to uploads directory
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
		return
	}

	// Update profile with new image path (gunakan FULL URL untuk frontend)
	imageURL := "http://localhost:8080/uploads/" + fileName
	err = h.repo.UpdateProfileImage(imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile image uploaded successfully",
		"image_url": imageURL,
	})
}

// Social Media methods
func (h *CafeProfileHandler) GetSocialMedia(c *gin.Context) {
	socialMedia, err := h.repo.GetSocialMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, socialMedia)
}

func (h *CafeProfileHandler) AddSocialMedia(c *gin.Context) {
	var socialMedia models.CafeSocialMedia
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.AddSocialMedia(&socialMedia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Social media added successfully"})
}

func (h *CafeProfileHandler) DeleteAllSocialMedia(c *gin.Context) {
	err := h.repo.DeleteAllSocialMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All social media deleted successfully"})
}

// Operational Hours methods
func (h *CafeProfileHandler) GetOperationalHours(c *gin.Context) {
	hours, err := h.repo.GetOperationalHours()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hours)
}

func (h *CafeProfileHandler) UpdateOperationalHours(c *gin.Context) {
	var hours []models.CafeOperationalHours
	if err := c.ShouldBindJSON(&hours); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.UpdateOperationalHours(hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Operational hours updated successfully"})
}

// ✅ METHOD BARU: GetTodayOperationalHours - Get jam operasional hari ini (format 24 jam)
func (h *CafeProfileHandler) GetTodayOperationalHours(c *gin.Context) {
	openingHours, err := h.repo.GetTodayOperationalHours()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, openingHours)
}

// ✅ METHOD BARU: UpdateSingleOperationalHours - Update jam operasional untuk hari tertentu
func (h *CafeProfileHandler) UpdateSingleOperationalHours(c *gin.Context) {
	var request struct {
		Hari  string `json:"hari"`  // "Senin", "Selasa", etc.
		Buka  string `json:"buka"`  // Format: "07:00"
		Tutup string `json:"tutup"` // Format: "22:00"
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.repo.UpdateSingleOperationalHours(request.Hari, request.Buka, request.Tutup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Return success response
	response := gin.H{
		"message": "Jam operasional berhasil diupdate",
		"hari":    request.Hari,
		"waktu":   request.Buka + " - " + request.Tutup,
	}
	
	c.JSON(http.StatusOK, response)
}

// ✅ METHOD BARU: GetCurrentStatus - Hanya ambil status buka/tutup
func (h *CafeProfileHandler) GetCurrentStatus(c *gin.Context) {
	status, err := h.repo.GetCurrentStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// Facilities methods
func (h *CafeProfileHandler) GetFacilities(c *gin.Context) {
	facilities, err := h.repo.GetFacilities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, facilities)
}

func (h *CafeProfileHandler) UpdateFacilities(c *gin.Context) {
	var facilities []models.CafeFacility
	if err := c.ShouldBindJSON(&facilities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.UpdateFacilities(facilities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Facilities updated successfully"})
}

// Gallery methods
func (h *CafeProfileHandler) GetGallery(c *gin.Context) {
	gallery, err := h.repo.GetGallery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gallery)
}

func (h *CafeProfileHandler) AddGalleryImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Create gallery directory if not exists
	if err := os.MkdirAll("uploads/gallery", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gallery directory: " + err.Error()})
		return
	}

	// Generate unique filename for gallery
	fileExt := filepath.Ext(file.Filename)
	fileName := "gallery_" + time.Now().Format("20060102150405") + fileExt
	filePath := "uploads/gallery/" + fileName

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
		return
	}

	urutan, _ := strconv.Atoi(c.PostForm("urutan"))
	
	gallery := &models.CafeGallery{
		ImageURL: "http://localhost:8080/uploads/gallery/" + fileName, // FULL URL
		Urutan:   urutan,
	}

	err = h.repo.AddGalleryImage(gallery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery image added successfully",
		"id": gallery.ID,
		"image_url": gallery.ImageURL,
		"urutan": gallery.Urutan,
	})
}

func (h *CafeProfileHandler) DeleteGalleryImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.repo.DeleteGalleryImage(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Gallery image deleted successfully"})
}