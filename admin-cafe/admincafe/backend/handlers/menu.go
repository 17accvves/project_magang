package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend/models"
	"backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuHandler struct {
	repo *repository.MenuRepository
}

func NewMenuHandler(repo *repository.MenuRepository) *MenuHandler {
	return &MenuHandler{repo: repo}
}

// Get All Menus
func (h *MenuHandler) GetMenus(c *gin.Context) {
	menus, err := h.repo.GetAllMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// Get Menu by ID
func (h *MenuHandler) GetMenu(c *gin.Context) {
	id := c.Param("id")

	menu, err := h.repo.GetMenuByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if menu == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// Create Menu
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi required fields
	if menu.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if menu.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be greater than 0"})
		return
	}
	if menu.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	}

	// 🔥 PERBAIKAN: Pastikan diskon tidak negatif dan handle default value
	if menu.Discount < 0 {
		menu.Discount = 0
	}

	// 🔥 PERBAIKAN: Validasi diskon dan tanggal yang lebih robust
	if menu.Discount > 0 {
		// Jika ada diskon, tanggal harus diisi
		if menu.StartDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start date is required when discount is set"})
			return
		}
		if menu.EndDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End date is required when discount is set"})
			return
		}

		// Validasi format tanggal
		if _, err := time.Parse("2006-01-02", menu.StartDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format, use YYYY-MM-DD"})
			return
		}
		if _, err := time.Parse("2006-01-02", menu.EndDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format, use YYYY-MM-DD"})
			return
		}

		// Validasi tanggal tidak boleh berlalu
		today := time.Now().Format("2006-01-02")
		if menu.EndDate < today {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be in the past"})
			return
		}

		// Validasi start date tidak boleh setelah end date
		if menu.StartDate > menu.EndDate {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be after end date"})
			return
		}
	} else {
		// 🔥 PERBAIKAN: Jika tidak ada diskon, pastikan tanggal kosong
		menu.Discount = 0
		menu.StartDate = ""
		menu.EndDate = ""
	}

	// Calculate discounted price
	if menu.Discount > 0 {
		menu.DiscountedPrice = menu.Price - (menu.Price * menu.Discount / 100)
	} else {
		menu.DiscountedPrice = menu.Price
	}

	if menu.Status == "" {
		menu.Status = "Aktif"
	}

	// 🔥 PERBAIKAN: Tambahkan debug logging untuk troubleshooting
	fmt.Printf("🔍 [DEBUG] Creating menu - Name: %s, Price: %.2f, Discount: %.2f, StartDate: '%s', EndDate: '%s'\n",
		menu.Name, menu.Price, menu.Discount, menu.StartDate, menu.EndDate)

	err := h.repo.CreateMenu(&menu)
	if err != nil {
		// 🔥 PERBAIKAN: Log error yang lebih detail
		fmt.Printf("❌ [ERROR] CreateMenu failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      menu.ID,
		"message": "Menu created successfully",
		"data":    menu,
	})
}

// Update Menu
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔥 PERBAIKAN: Tambahkan debug logging
	fmt.Printf("🔍 [DEBUG] Update request for ID: %s, Data: %+v\n", id, updateData)

	// Validasi diskon dan tanggal jika ada field diskon yang diupdate
	if discount, hasDiscount := updateData["discount"]; hasDiscount {
		discountFloat, ok := discount.(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Discount must be a number"})
			return
		}

		// 🔥 PERBAIKAN: Pastikan diskon tidak negatif
		if discountFloat < 0 {
			discountFloat = 0
			updateData["discount"] = 0.0
		}

		// Handle field name mapping untuk tanggal
		startDate, hasStartDate := updateData["startDate"]
		endDate, hasEndDate := updateData["endDate"]
		
		// Jika field menggunakan underscore, handle juga
		if !hasStartDate {
			startDate = updateData["start_date"]
			hasStartDate = startDate != nil
		}
		if !hasEndDate {
			endDate = updateData["end_date"]
			hasEndDate = endDate != nil
		}

		startDateStr, _ := startDate.(string)
		endDateStr, _ := endDate.(string)

		// Jika ada diskon > 0, validasi tanggal harus diisi
		if discountFloat > 0 {
			if !hasStartDate || startDateStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Start date is required when discount is set"})
				return
			}
			if !hasEndDate || endDateStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "End date is required when discount is set"})
				return
			}

			// Validasi format tanggal
			if _, err := time.Parse("2006-01-02", startDateStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format, use YYYY-MM-DD"})
				return
			}
			if _, err := time.Parse("2006-01-02", endDateStr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format, use YYYY-MM-DD"})
				return
			}

			// Validasi tanggal tidak boleh berlalu
			today := time.Now().Format("2006-01-02")
			if endDateStr < today {
				c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be in the past"})
				return
			}

			// Validasi start date tidak boleh setelah end date
			if startDateStr > endDateStr {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be after end date"})
				return
			}
		} else {
			// 🔥 PERBAIKAN: Jika diskon = 0, kosongkan tanggal
			updateData["startDate"] = ""
			updateData["endDate"] = ""
			updateData["start_date"] = ""
			updateData["end_date"] = ""
		}
	}

	// Handle field name mapping untuk konsistensi
	if startDate, exists := updateData["startDate"]; exists {
		updateData["start_date"] = startDate
		delete(updateData, "startDate")
	}
	if endDate, exists := updateData["endDate"]; exists {
		updateData["end_date"] = endDate
		delete(updateData, "endDate")
	}

	// Calculate discounted price if price or discount is updated
	if price, exists := updateData["price"]; exists {
		if priceFloat, ok := price.(float64); ok {
			discount := 0.0
			if disc, exists := updateData["discount"]; exists {
				if discFloat, ok := disc.(float64); ok {
					discount = discFloat
				}
			} else {
				// Get current discount from database
				currentMenu, err := h.repo.GetMenuByID(id)
				if err == nil && currentMenu != nil {
					discount = currentMenu.Discount
				}
			}

			if discount > 0 {
				updateData["discounted_price"] = priceFloat - (priceFloat * discount / 100)
			} else {
				updateData["discounted_price"] = priceFloat
			}
		}
	} else if discount, exists := updateData["discount"]; exists {
		if discountFloat, ok := discount.(float64); ok {
			// Get current price from database
			currentMenu, err := h.repo.GetMenuByID(id)
			if err == nil && currentMenu != nil {
				if discountFloat > 0 {
					updateData["discounted_price"] = currentMenu.Price - (currentMenu.Price * discountFloat / 100)
				} else {
					updateData["discounted_price"] = currentMenu.Price
				}
			}
		}
	}

	err := h.repo.UpdateMenu(id, updateData)
	if err != nil {
		// 🔥 PERBAIKAN: Log error yang lebih detail
		fmt.Printf("❌ [ERROR] UpdateMenu failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu updated successfully"})
}

// Delete Menu
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	err := h.repo.DeleteMenu(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

// Upload Image
func (h *MenuHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	ext := filepath.Ext(file.Filename)
	if !allowedTypes[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: JPG, JPEG, PNG, GIF, WEBP"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum 5MB allowed"})
		return
	}

	// Create uploads directory if not exists
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return accessible URL
	url := "http://localhost:8080/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{
		"url":      url,
		"filename": filename,
		"message":  "Image uploaded successfully",
	})
}

// Get Menus by Category
func (h *MenuHandler) GetMenusByCategory(c *gin.Context) {
	category := c.Param("category")
	
	menus, err := h.repo.GetMenusByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// Search Menus
func (h *MenuHandler) SearchMenus(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	menus, err := h.repo.SearchMenus(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}