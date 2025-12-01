package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
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
		title VARCHAR(255) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		revenue VARCHAR(50) NOT NULL,
		status VARCHAR(20) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_promo_status ON promos(status);
	CREATE INDEX IF NOT EXISTS idx_promo_dates ON promos(start_date, end_date);
	`
	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create promos table:", err)
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
		log.Printf("Error checking existing promos: %v", err)
		return
	}

	if count == 0 {
		samplePromos := []struct {
			title     string
			startDate string
			endDate   string
			revenue   string
		}{
			{"Diskon Promo 25% Matcha", "2025-09-18", "2025-09-25", "17JT"},
			{"Diskon Promo 15% Kopi Latte", "2025-09-10", "2025-09-20", "12JT"},
			{"Diskon Promo 10% Snack Pukis", "2025-09-05", "2025-09-12", "9JT"},
			{"Diskon Promo 20% Extra Jose", "2025-09-01", "2025-09-10", "14JT"},
			{"Diskon Promo 30% Pentolaan James", "2025-08-22", "2025-08-28", "20JT"},
		}

		for _, promo := range samplePromos {
			_, err := DB.Exec(
				"INSERT INTO promos (title, start_date, end_date, revenue) VALUES ($1, $2, $3, $4)",
				promo.title, promo.startDate, promo.endDate, promo.revenue,
			)
			if err != nil {
				log.Printf("Failed to insert sample promo: %v", err)
			}
		}
		fmt.Println("Sample promos inserted successfully!")
	}
}
