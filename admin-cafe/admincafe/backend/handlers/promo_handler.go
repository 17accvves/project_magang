package handlers

import (
	"backend/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PromoHandler struct {
	menuRepo *repository.MenuRepository
}

func NewPromoHandler(menuRepo *repository.MenuRepository) *PromoHandler {
	return &PromoHandler{
		menuRepo: menuRepo,
	}
}

// GetAllPromos menampilkan semua promo DARI MENU YANG ADA DISKON
func (h *PromoHandler) GetAllPromos(c *gin.Context) {
    fmt.Printf("🎯 [HANDLER] GET /api/v1/promos called\n")
    
    // Ambil menu yang ada diskon
    menus, err := h.menuRepo.GetMenusWithDiscount()
    if err != nil {
        fmt.Printf("❌ [HANDLER] Error fetching promo menus: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   true,
            "message": "Failed to fetch promo menus",
            "data":    nil,
        })
        return
    }

    fmt.Printf("✅ [HANDLER] Returning %d promo menus\n", len(menus))
    
    c.JSON(http.StatusOK, gin.H{
        "error":   false,
        "message": "Promos retrieved successfully",
        "data":    menus,
    })
}

// GetPromoStats menampilkan statistik promo DARI MENU YANG ADA DISKON
func (h *PromoHandler) GetPromoStats(c *gin.Context) {
	fmt.Printf("📊 [HANDLER] GET /api/v1/promos/stats called\n")
	
	stats, err := h.menuRepo.GetPromoStats()
	if err != nil {
		fmt.Printf("❌ [HANDLER] Error fetching promo stats: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to fetch promo stats",
			"data":    nil,
		})
		return
	}

	fmt.Printf("✅ [HANDLER] Returning promo stats - Active: %d, Revenue: %s\n", 
		stats.ActivePromos, stats.TotalRevenue)
	
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Promo stats retrieved successfully",
		"data":    stats,
	})
}

// GetPromoByID - LEGACY: Tidak digunakan karena promo sekarang dari menu
func (h *PromoHandler) GetPromoByID(c *gin.Context) {
	promoID := c.Param("id")
	fmt.Printf("⚠️  [HANDLER] LEGACY CALL: GET /api/v1/promos/%s (promos are now from menus)\n", promoID)
	
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Promos are now retrieved from menus with discounts",
		"data":    nil,
	})
}

// CreatePromo - LEGACY: Tidak digunakan karena promo sekarang dari menu
func (h *PromoHandler) CreatePromo(c *gin.Context) {
	fmt.Printf("⚠️  [HANDLER] LEGACY CALL: POST /api/v1/promos (promos are now created via menus)\n")
	
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Promos are now created by adding discount to menus in the Menu page",
		"data":    nil,
	})
}

// UpdatePromo - LEGACY: Tidak digunakan karena promo sekarang dari menu
func (h *PromoHandler) UpdatePromo(c *gin.Context) {
	promoID := c.Param("id")
	fmt.Printf("⚠️  [HANDLER] LEGACY CALL: PUT /api/v1/promos/%s (promos are now updated via menus)\n", promoID)
	
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Promos are now updated by modifying menu discounts in the Menu page",
		"data":    nil,
	})
}

// DeletePromo - LEGACY: Tidak digunakan karena promo sekarang dari menu
func (h *PromoHandler) DeletePromo(c *gin.Context) {
	promoID := c.Param("id")
	fmt.Printf("⚠️  [HANDLER] LEGACY CALL: DELETE /api/v1/promos/%s (promos are now removed via menus)\n", promoID)
	
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Promos are now removed by setting menu discount to 0 in the Menu page",
		"data":    nil,
	})
}