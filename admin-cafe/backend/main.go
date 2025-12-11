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
    // 1️⃣ Koneksi ke database
    config.ConnectDB()

    // 2️⃣ Buat repository & handler
    userRepo := repository.NewUserRepository()
    authHandler := handlers.NewAuthHandler(userRepo)
    cafeHandler := handlers.NewCafeHandler(userRepo)

    // 3️⃣ Pastikan folder uploads ada
    ensureUploadsFolder()

    // 4️⃣ Setup router & routes
    router := routes.SetupRoutes(authHandler, cafeHandler)

    // 5️⃣ Serve file uploads
    router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

    // 6️⃣ Konfigurasi CORS untuk React
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"}, // React dev server
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })
    handler := c.Handler(router)

    // 7️⃣ Jalankan server
    fmt.Println("Server running at :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}

// ensureUploadsFolder memastikan folder uploads ada
func ensureUploadsFolder() {
    if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
        err := os.Mkdir("./uploads", os.ModePerm)
        if err != nil {
            log.Fatal("Gagal membuat folder uploads:", err)
        }
        fmt.Println("Folder uploads dibuat")
    }
}
