package main

import (
    "backend/config"
    "backend/handlers"
    "backend/repository"
    "backend/routes"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/rs/cors"
)

func main() {
	// =========================
	// 1️⃣ Connect PostgreSQL
	// =========================
	config.ConnectDB() // di sini juga sudah bikin tabel users & admin default

	// =========================
	// 2️⃣ Initialize Repositories
	// =========================
	userRepo := repository.NewUserRepository()
	cafeRepo := repository.NewCafeProfileRepository(config.DB)
	menuRepo := repository.NewMenuRepository(config.DB)
	ulasanRepo := repository.NewUlasanRepository(config.DB)

	// =========================
	// 3️⃣ Initialize Handlers
	// =========================
	authHandler := handlers.NewAuthHandler(userRepo) // pakai userRepo
	cafeHandler := handlers.NewCafeProfileHandler(cafeRepo)
	menuHandler := handlers.NewMenuHandler(menuRepo)
	ulasanHandler := handlers.NewUlasanHandler(ulasanRepo)
	promoHandler := handlers.NewPromoHandler(menuRepo)

	// =========================
	// 4️⃣ Start discount checker (background)
	// =========================
	go menuRepo.StartDiscountChecker()

	// =========================
	// 5️⃣ Setup router
	// =========================
	router := gin.Default()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"*"}
	router.Use(cors.New(corsConfig))

	// Logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	})

	// Serve static files
	router.Static("/uploads", "./uploads")

	// =========================
	// 6️⃣ Auth Routes
	// =========================
	router.POST("/login", authHandler.Login)

	// =========================
	// 7️⃣ Cafe Profile Routes
	// =========================
	cafeRoutes := router.Group("/cafe")
	{
		cafeRoutes.GET("/profile", cafeHandler.GetCafeProfile)
		cafeRoutes.PUT("/profile", cafeHandler.UpdateCafeProfile)
		cafeRoutes.POST("/profile/image", cafeHandler.UploadProfileImage)

		cafeRoutes.GET("/social-media", cafeHandler.GetSocialMedia)
		cafeRoutes.POST("/social-media", cafeHandler.AddSocialMedia)
		cafeRoutes.DELETE("/social-media/all", cafeHandler.DeleteAllSocialMedia)

		cafeRoutes.GET("/operational-hours", cafeHandler.GetOperationalHours)
		cafeRoutes.PUT("/operational-hours", cafeHandler.UpdateOperationalHours)

		cafeRoutes.GET("/facilities", cafeHandler.GetFacilities)
		cafeRoutes.PUT("/facilities", cafeHandler.UpdateFacilities)

		cafeRoutes.GET("/gallery", cafeHandler.GetGallery)
		cafeRoutes.POST("/gallery", cafeHandler.AddGalleryImage)
		cafeRoutes.DELETE("/gallery/:id", cafeHandler.DeleteGalleryImage)
	}

	// Operational Hours public API
	operational := router.Group("/api/operational-hours")
	{
		operational.GET("/today", cafeHandler.GetTodayOperationalHours)
		operational.PUT("/single", cafeHandler.UpdateSingleOperationalHours)
		operational.GET("/status", cafeHandler.GetCurrentStatus)
	}

	// =========================
	// 8️⃣ Menu Routes
	// =========================
	menuApi := router.Group("/menus")
	{
		menuApi.GET("", menuHandler.GetMenus)
		menuApi.GET("/:id", menuHandler.GetMenu)
		menuApi.POST("", menuHandler.CreateMenu)
		menuApi.PUT("/:id", menuHandler.UpdateMenu)
		menuApi.DELETE("/:id", menuHandler.DeleteMenu)
	}
	router.POST("/upload", menuHandler.UploadImage)

	// =========================
	// 9️⃣ Ulasan Routes
	// =========================
	ulasanApi := router.Group("/ulasan")
	{
		ulasanApi.GET("", ulasanHandler.GetUlasan)
		ulasanApi.POST("", ulasanHandler.CreateUlasan)
		ulasanApi.POST("/upload", ulasanHandler.UploadGambarUlasan)
		ulasanApi.POST("/upload-avatar", ulasanHandler.UploadAvatarUlasan)
	}

	adminUlasan := router.Group("/admin/ulasan")
	{
		adminUlasan.GET("", ulasanHandler.GetUlasanAdmin)
		adminUlasan.GET("/stats", ulasanHandler.GetUlasanStats)
		adminUlasan.PUT("/:id/reply", ulasanHandler.AddReply)
		adminUlasan.DELETE("/:id/reply", ulasanHandler.DeleteReply)
		adminUlasan.PUT("/:id/status", ulasanHandler.UpdateUlasanStatus)
		adminUlasan.DELETE("/:id", ulasanHandler.DeleteUlasan)
	}

	// =========================
	// 10️⃣ Promo Routes
	// =========================
	promoApi := router.Group("/api/v1/promos")
	{
		promoApi.GET("", promoHandler.GetAllPromos)
		promoApi.GET("/stats", promoHandler.GetPromoStats)
		promoApi.GET("/:id", promoHandler.GetPromoByID)
		promoApi.POST("", promoHandler.CreatePromo)
		promoApi.PUT("/:id", promoHandler.UpdatePromo)
		promoApi.DELETE("/:id", promoHandler.DeletePromo)
	}

	// =========================
	// 11️⃣ Health Check
	// =========================
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "Cafe Management API",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// =========================
	// 12️⃣ Debug Routes
	// =========================
	logRegisteredRoutes(router)

	log.Println("✅ Server started successfully on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// =========================
// Helper: Print semua routes
// =========================
func logRegisteredRoutes(router *gin.Engine) {
	log.Println("📋 REGISTERED ROUTES:")
	for _, r := range router.Routes() {
		log.Printf("%s %s", r.Method, r.Path)
	}
	log.Println("📋 END OF REGISTERED ROUTES")
}
