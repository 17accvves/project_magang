package repository

import (
	"backend/config"
	"database/sql"
	"fmt"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type UserRepository struct {
	db *sql.DB
}

// =========================
// Constructor
// =========================
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: config.DB,
	}
}

// =========================
// Login (tanpa hash)
// =========================
func (r *UserRepository) Login(username, password string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, password, role FROM users WHERE username=$1"
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Cocokkan password langsung (tanpa hash)
	if user.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// =========================
// Buat user baru
// =========================
func (r *UserRepository) CreateUser(username, password, role string) (*User, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id",
		username, password, role,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Username: username,
		Password: password,
		Role:     role,
	}, nil
}

// =========================
// Cek apakah username sudah ada
// =========================
func (r *UserRepository) Exists(username string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// =========================
// Hapus user
// =========================
func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
