package repository

import (
    "backend/config"
    "backend/models"
		"database/sql"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        DB: config.DB,
    }
}

// Cari user berdasarkan username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
    user := &models.User{}
    row := config.DB.QueryRow("SELECT id, username, password, role, verified FROM users WHERE username=$1", username)
    err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Verified)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Buat user cafe baru
func (r *UserRepository) CreateCafe(user *models.User) error {
    _, err := config.DB.Exec(
        "INSERT INTO users (username, password, email, role, izin_usaha, verified) VALUES ($1,$2,$3,$4,$5,false)",
        user.Username, user.Password, user.Email, user.Role, user.IzinUsaha,
    )
    return err
}

// Update verified cafe
func (r *UserRepository) VerifyCafe(userID int) error {
    _, err := config.DB.Exec("UPDATE users SET verified=true WHERE id=$1", userID)
    return err
}
