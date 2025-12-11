package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ganteng707"
	dbname   = "carispot_db"
)

func ConnectDB() {
	// DSN menggunakan database carispot_db
	dsn := "host=localhost port=5432 user=postgres password=111122 dbname=carispot_db sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	// Set connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Buat semua tabel
	createTables()

	// Buat tabel users & admin default
	createUsersTable()
	createAdminDefault()
}

// =========================
// CREATE ALL TABLES
// =========================
func createTables() {
	createCafeTables()
	createMenuTable()
	createUlasanTable()
	createPromoTable()
}

// =========================
// USERS TABLE & ADMIN DEFAULT
// =========================
func createUsersTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}

func createAdminDefault() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username='admin'").Scan(&count)
	if err != nil {
		log.Printf("Error checking admin user: %v", err)
		return
	}

	if count == 0 {
		_, err := DB.Exec(
			"INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
			"admin", "1234", "admin",
		)
		if err != nil {
			log.Printf("Failed to insert admin default: %v", err)
			return
		}
		fmt.Println("✅ Admin default dibuat -> username: admin | password: 1234")
	} else {
		fmt.Println("Admin default sudah ada, tidak dibuat ulang.")
	}
}

// =========================
// CAFE TABLES
// =========================
func createCafeTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS cafe_profiles (
			id SERIAL PRIMARY KEY,
			nama VARCHAR(255) NOT NULL,
			alamat TEXT NOT NULL,
			telepon VARCHAR(20),
			deskripsi TEXT,
			main_image VARCHAR(500),
			verified BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cafe_social_media (
			id SERIAL PRIMARY KEY,
			cafe_profile_id INTEGER REFERENCES cafe_profiles(id) ON DELETE CASCADE,
			platform VARCHAR(50) NOT NULL,
			url VARCHAR(500) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cafe_operational_hours (
			id SERIAL PRIMARY KEY,
			cafe_profile_id INTEGER REFERENCES cafe_profiles(id) ON DELETE CASCADE,
			hari VARCHAR(10) NOT NULL,
			buka TIME NOT NULL,
			tutup TIME NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cafe_facilities (
			id SERIAL PRIMARY KEY,
			cafe_profile_id INTEGER REFERENCES cafe_profiles(id) ON DELETE CASCADE,
			nama_fasilitas VARCHAR(100) NOT NULL,
			tersedia BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cafe_gallery (
			id SERIAL PRIMARY KEY,
			cafe_profile_id INTEGER REFERENCES cafe_profiles(id) ON DELETE CASCADE,
			image_url VARCHAR(500) NOT NULL,
			urutan INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Printf("Failed to create table: %v", err)
		}
	}
	fmt.Println("Cafe tables created successfully!")
}

// =========================
// MENUS TABLE
// =========================
func createMenuTable() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS menus (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		discount DECIMAL(5,2) DEFAULT 0,
		discounted_price DECIMAL(10,2) DEFAULT 0,
		start_date DATE,
		end_date DATE,
		category VARCHAR(100) NOT NULL,
		status VARCHAR(50) DEFAULT 'Aktif',
		img TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_menu_name ON menus(name);
	CREATE INDEX IF NOT EXISTS idx_menu_category ON menus(category);
	CREATE INDEX IF NOT EXISTS idx_menu_status ON menus(status);
	`
	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create menus table:", err)
	}
	fmt.Println("Menus table created successfully!")
}

// =========================
// ULASAN TABLE
// =========================
func createUlasanTable() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS ulasan (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nama VARCHAR(255) NOT NULL,
		email VARCHAR(255),
		rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
		teks TEXT NOT NULL,
		gambar TEXT,
		avatar TEXT,
		balasan TEXT,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_ulasan_status ON ulasan(status);
	CREATE INDEX IF NOT EXISTS idx_ulasan_rating ON ulasan(rating);
	CREATE INDEX IF NOT EXISTS idx_ulasan_created_at ON ulasan(created_at);
	`
	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create ulasan table:", err)
	}
	fmt.Println("Ulasan table created successfully!")
}

// =========================
// PROMOS TABLE
// =========================
func createPromoTable() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS promos (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(100),
		role VARCHAR(20) NOT NULL,
		izin_usaha TEXT,
		verified BOOLEAN DEFAULT false,
		rejected BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Promos table created successfully!")
	insertSamplePromos()
}

// =========================
// SAMPLE PROMOS
// =========================
func insertSamplePromos() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM promos").Scan(&count)
	if err != nil {
		log.Fatal("Check admin error:", err)
	}

	if !exists {
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
