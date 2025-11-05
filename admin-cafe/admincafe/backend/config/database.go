package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "111122" // Ganti dengan password PostgreSQL Anda
	dbname   = "menu_db"
)

func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Create tables if they don't exist
	createTables()
}

func createTables() {
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
		log.Fatal("Failed to create tables:", err)
	}

	fmt.Println("Tables created successfully!")
}
