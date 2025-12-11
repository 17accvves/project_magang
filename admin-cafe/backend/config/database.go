package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
    // DSN menggunakan database carispot_db
    dsn := "host=localhost port=5432 user=postgres password=111122 dbname=carispot_db sslmode=disable"
    var err error
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("DB connection error:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("DB ping error:", err)
    }

    fmt.Println("Database connected")

    // Pastikan folder uploads ada
    ensureUploadsFolder()

    // Buat table users jika belum ada
    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        email VARCHAR(100),
        role VARCHAR(20) NOT NULL,
        izin_usaha TEXT,
        verified BOOLEAN DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err = DB.Exec(createTable)
    if err != nil {
        log.Fatal("Failed to create table:", err)
    }
    fmt.Println("Table users ready")

    // Buat default admin jika belum ada
    createDefaultAdmin()
}

// Buat default admin
func createDefaultAdmin() {
    // Cek apakah admin sudah ada
    var exists bool
    row := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username='admin')")
    err := row.Scan(&exists)
    if err != nil {
        log.Fatal("Check admin error:", err)
    }

    if !exists {
        // Simpan password plain text
        _, err := DB.Exec(
            "INSERT INTO users (username, password, role, verified) VALUES ($1,$2,$3,true)",
            "admin", "admin123", "admin",
        )
        if err != nil {
            log.Fatal("Failed to create default admin:", err)
        }
        fmt.Println("Default admin created (username: admin, password: admin123)")
    } else {
        fmt.Println("Default admin already exists")
    }
}

// Pastikan folder uploads ada
func ensureUploadsFolder() {
    if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
        os.Mkdir("./uploads", os.ModePerm)
        fmt.Println("Folder uploads dibuat")
    }
}
