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

    // Calculate discounted price
    if menu.Discount > 0 {
        menu.DiscountedPrice = menu.Price - (menu.Price * menu.Discount / 100)
    } else {
        menu.DiscountedPrice = menu.Price
    }

    if menu.Status == "" {
        menu.Status = "Aktif"
    }

    err := h.repo.CreateMenu(&menu)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":      menu.ID,
        "message": "Menu created successfully",
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

    // Handle field name mapping
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
    c.JSON(http.StatusOK, gin.H{"url": url})
}