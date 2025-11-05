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

	// Initialize repository and handlers
	menuRepo := repository.NewMenuRepository(config.DB)
	menuHandler := handlers.NewMenuHandler(menuRepo)

	// Setup router
	router := gin.Default()

	// SUPER SIMPLE CORS - untuk development
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	router.Use(cors.New(config))

	// Middleware logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	})

	// Serve static files
	router.Static("/uploads", "./uploads")

	// Routes
	api := router.Group("/menus")
	{
		api.GET("", menuHandler.GetMenus)
		api.GET("/:id", menuHandler.GetMenu)
		api.POST("", menuHandler.CreateMenu)
		api.PUT("/:id", menuHandler.UpdateMenu)
		api.DELETE("/:id", menuHandler.DeleteMenu)
	}

	// Upload route
	router.POST("/upload", menuHandler.UploadImage)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "Menu API",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Menu API Server is running",
			"version": "1.0",
		})
	})

	log.Println("✅ Server started successfully on :8080")
	log.Println("🌐 CORS: Enabled for all origins (development)")
	log.Println("📊 Test: http://localhost:8080/health")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
