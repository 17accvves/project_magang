// main.go - PERBAIKI ROUTING
package main

import (
	"log"
	"time"

	"backend/config"
	"backend/handlers"
	"backend/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to PostgreSQL
	config.ConnectDB()

	// Initialize repositories
	cafeProfileRepo := repository.NewCafeProfileRepository(config.DB)
	menuRepo := repository.NewMenuRepository(config.DB)
	ulasanRepo := repository.NewUlasanRepository(config.DB)

	// Initialize handlers
	cafeProfileHandler := handlers.NewCafeProfileHandler(cafeProfileRepo)
	menuHandler := handlers.NewMenuHandler(menuRepo)
	ulasanHandler := handlers.NewUlasanHandler(ulasanRepo)
	promoHandler := handlers.NewPromoHandler(menuRepo) // ✅ PERBAIKI: Pass menuRepo

	// Start discount checker
	go menuRepo.StartDiscountChecker()

	// Setup router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	// Middleware logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	})

	// Serve static files
	router.Static("/uploads", "./uploads")

	// === CAFE PROFILE ROUTES ===
	cafeRoutes := router.Group("/cafe")
	{
		cafeRoutes.GET("/profile", cafeProfileHandler.GetCafeProfile)
		cafeRoutes.PUT("/profile", cafeProfileHandler.UpdateCafeProfile)
		cafeRoutes.POST("/profile/image", cafeProfileHandler.UploadProfileImage)
		
		cafeRoutes.GET("/social-media", cafeProfileHandler.GetSocialMedia)
		cafeRoutes.POST("/social-media", cafeProfileHandler.AddSocialMedia)
		cafeRoutes.DELETE("/social-media/all", cafeProfileHandler.DeleteAllSocialMedia)
		
		cafeRoutes.GET("/operational-hours", cafeProfileHandler.GetOperationalHours)
		cafeRoutes.PUT("/operational-hours", cafeProfileHandler.UpdateOperationalHours)
		
		cafeRoutes.GET("/facilities", cafeProfileHandler.GetFacilities)
		cafeRoutes.PUT("/facilities", cafeProfileHandler.UpdateFacilities)
		
		cafeRoutes.GET("/gallery", cafeProfileHandler.GetGallery)
		cafeRoutes.POST("/gallery", cafeProfileHandler.AddGalleryImage)
		cafeRoutes.DELETE("/gallery/:id", cafeProfileHandler.DeleteGalleryImage)
	}

	// ✅ ROUTES BARU: Operational Hours dengan format 24 jam - PERBAIKI PATH
	operationalRoutes := router.Group("/api")
	{
		operationalRoutes.GET("/operational-hours/today", cafeProfileHandler.GetTodayOperationalHours)
		operationalRoutes.PUT("/operational-hours/single", cafeProfileHandler.UpdateSingleOperationalHours)
		operationalRoutes.GET("/operational-hours/status", cafeProfileHandler.GetCurrentStatus)
	}

	// Menu Routes (existing)
	menuApi := router.Group("/menus")
	{
		menuApi.GET("", menuHandler.GetMenus)
		menuApi.GET("/:id", menuHandler.GetMenu)
		menuApi.POST("", menuHandler.CreateMenu)
		menuApi.PUT("/:id", menuHandler.UpdateMenu)
		menuApi.DELETE("/:id", menuHandler.DeleteMenu)
	}

	// Ulasan Routes (existing)
	ulasanApi := router.Group("/ulasan")
	{
		ulasanApi.GET("", ulasanHandler.GetUlasan)
		ulasanApi.POST("", ulasanHandler.CreateUlasan)
		ulasanApi.POST("/upload", ulasanHandler.UploadGambarUlasan)
		ulasanApi.POST("/upload-avatar", ulasanHandler.UploadAvatarUlasan)
	}

	// Menu Upload route
	router.POST("/upload", menuHandler.UploadImage)

	// Admin Ulasan Routes (existing)
	adminUlasan := router.Group("/admin/ulasan")
	{
		adminUlasan.GET("", ulasanHandler.GetUlasanAdmin)
		adminUlasan.GET("/stats", ulasanHandler.GetUlasanStats)
		adminUlasan.PUT("/:id/reply", ulasanHandler.AddReply)
		adminUlasan.DELETE("/:id/reply", ulasanHandler.DeleteReply)
		adminUlasan.PUT("/:id/status", ulasanHandler.UpdateUlasanStatus)
		adminUlasan.DELETE("/:id", ulasanHandler.DeleteUlasan)
	}

	// ✅ PROMO ROUTES - SEKARANG AMBIL DATA DARI MENU YANG ADA DISKON
	promoApi := router.Group("/api/v1/promos")
	{
		promoApi.GET("", promoHandler.GetAllPromos)           // Get semua promo dari menu
		promoApi.GET("/stats", promoHandler.GetPromoStats)    // Get statistik promo
		// Routes legacy bisa dihapus atau dipertahankan untuk compatibility
		promoApi.GET("/:id", promoHandler.GetPromoByID)
		promoApi.POST("", promoHandler.CreatePromo)
		promoApi.PUT("/:id", promoHandler.UpdatePromo)
		promoApi.DELETE("/:id", promoHandler.DeletePromo)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "Cafe Management API",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// ✅ DEBUG: Print semua routes yang terdaftar
	logRegisteredRoutes(router)

	log.Println("✅ Server started successfully on :8080")
	log.Println("🏪 Cafe Profile API ready")
	log.Println("⏰ Operational Hours API ready (24-hour format)")
	log.Println("🎯 Promo API ready (Now using menu data with discounts)")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// ✅ FUNCTION BARU: Debug untuk melihat semua routes yang terdaftar
func logRegisteredRoutes(router *gin.Engine) {
	log.Println("📋 REGISTERED ROUTES:")
	routes := router.Routes()
	for _, route := range routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}
	log.Println("📋 END OF REGISTERED ROUTES")
}