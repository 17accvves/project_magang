package repository

import (
    "database/sql"
    "backend/models"
    "time"
    "fmt"

    "github.com/google/uuid"
)

type UlasanRepository struct {
    db *sql.DB
}

func NewUlasanRepository(db *sql.DB) *UlasanRepository {
    return &UlasanRepository{db: db}
}

func (r *UlasanRepository) GetAllUlasan() ([]models.UlasanResponse, error) {
    query := `
        SELECT id, nama, rating, teks, gambar, avatar, balasan, created_at 
        FROM ulasan 
        WHERE status = 'approved'
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ulasanList []models.UlasanResponse
    for rows.Next() {
        var ulasan models.Ulasan
        var gambar, avatar, balasan sql.NullString
        
        err := rows.Scan(
            &ulasan.ID,
            &ulasan.Nama,
            &ulasan.Rating,
            &ulasan.Teks,
            &gambar,
            &avatar,
            &balasan,
            &ulasan.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        // Format response untuk frontend
        response := models.UlasanResponse{
            ID:     ulasan.ID,
            Name:   ulasan.Nama,
            Date:   ulasan.CreatedAt.Format("02 Jan 2006"),
            Rating: ulasan.Rating,
            Text:   ulasan.Teks,
            Image:  getStringValue(gambar),
            Avatar: getStringValue(avatar),
            Reply:  getStringValue(balasan),
        }

        ulasanList = append(ulasanList, response)
    }

    return ulasanList, nil
}

func (r *UlasanRepository) GetAllUlasanForAdmin() ([]models.Ulasan, error) {
    query := `
        SELECT id, nama, email, rating, teks, gambar, avatar, balasan, status, created_at, updated_at 
        FROM ulasan 
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ulasanList []models.Ulasan
    for rows.Next() {
        var ulasan models.Ulasan
        var gambar, avatar, balasan sql.NullString
        
        err := rows.Scan(
            &ulasan.ID,
            &ulasan.Nama,
            &ulasan.Email,
            &ulasan.Rating,
            &ulasan.Teks,
            &gambar,
            &avatar,
            &balasan,
            &ulasan.Status,
            &ulasan.CreatedAt,
            &ulasan.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        ulasan.Gambar = getStringValue(gambar)
        ulasan.Avatar = getStringValue(avatar)
        ulasan.Balasan = getStringValue(balasan)

        ulasanList = append(ulasanList, ulasan)
    }

    return ulasanList, nil
}

func (r *UlasanRepository) GetUlasanByID(id string) (*models.Ulasan, error) {
    query := `
        SELECT id, nama, email, rating, teks, gambar, avatar, balasan, status, created_at, updated_at 
        FROM ulasan 
        WHERE id = $1
    `
    
    var ulasan models.Ulasan
    var gambar, avatar, balasan sql.NullString
    
    err := r.db.QueryRow(query, id).Scan(
        &ulasan.ID,
        &ulasan.Nama,
        &ulasan.Email,
        &ulasan.Rating,
        &ulasan.Teks,
        &gambar,
        &avatar,
        &balasan,
        &ulasan.Status,
        &ulasan.CreatedAt,
        &ulasan.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    ulasan.Gambar = getStringValue(gambar)
    ulasan.Avatar = getStringValue(avatar)
    ulasan.Balasan = getStringValue(balasan)

    return &ulasan, nil
}

func (r *UlasanRepository) CreateUlasan(ulasan *models.Ulasan) error {
    ulasan.ID = uuid.New().String()
    ulasan.CreatedAt = time.Now()
    ulasan.UpdatedAt = time.Now()
    
    if ulasan.Status == "" {
        ulasan.Status = "pending"
    }

    // Generate avatar jika tidak ada
    if ulasan.Avatar == "" {
        ulasan.Avatar = fmt.Sprintf("https://i.pravatar.cc/100?u=%s", ulasan.ID)
    }

    query := `
        INSERT INTO ulasan (id, nama, email, rating, teks, gambar, avatar, balasan, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

    _, err := r.db.Exec(
        query,
        ulasan.ID,
        ulasan.Nama,
        ulasan.Rating,
        ulasan.Teks,
        ulasan.Status,
        ulasan.CreatedAt,
        ulasan.UpdatedAt,
    )

    return err
}

func (r *UlasanRepository) UpdateUlasan(id string, updateData map[string]interface{}) error {
    updateData["updated_at"] = time.Now()

    query := "UPDATE ulasan SET "
    params := []interface{}{}
    paramCount := 1

    for key, value := range updateData {
        query += key + " = $" + fmt.Sprint(paramCount) + ", "
        params = append(params, value)
        paramCount++
    }

    query = query[:len(query)-2]
    query += " WHERE id = $" + fmt.Sprint(paramCount)
    params = append(params, id)

    _, err := r.db.Exec(query, params...)
    return err
}

func (r *UlasanRepository) AddReply(id string, reply string) error {
    query := "UPDATE ulasan SET balasan = $1, updated_at = $2 WHERE id = $3"
    
    _, err := r.db.Exec(query, reply, time.Now(), id)
    return err
}

func (r *UlasanRepository) DeleteReply(id string) error {
    query := "UPDATE ulasan SET balasan = NULL, updated_at = $1 WHERE id = $2"
    
    _, err := r.db.Exec(query, time.Now(), id)
    return err
}

func (r *UlasanRepository) DeleteUlasan(id string) error {
    query := "DELETE FROM ulasan WHERE id = $1"
    _, err := r.db.Exec(query, id)
    return err
}

func (r *UlasanRepository) GetUlasanStats() (map[string]interface{}, error) {
    stats := make(map[string]interface{})

    // Total ulasan
    var total int
    err := r.db.QueryRow("SELECT COUNT(*) FROM ulasan").Scan(&total)
    if err != nil {
        return nil, err
    }
    stats["total"] = total

    // Ulasan per status
    statusQuery := `
        SELECT status, COUNT(*) as count 
        FROM ulasan 
        GROUP BY status
    `
    rows, err := r.db.Query(statusQuery)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    statusCounts := make(map[string]int)
    for rows.Next() {
        var status string
        var count int
        rows.Scan(&status, &count)
        statusCounts[status] = count
    }
    stats["status"] = statusCounts

    // Rating average
    var avgRating float64
    err = r.db.QueryRow("SELECT COALESCE(AVG(rating), 0) FROM ulasan WHERE status = 'approved'").Scan(&avgRating)
    if err != nil {
        return nil, err
    }
    stats["avg_rating"] = avgRating

    // Ulasan dengan balasan
    var repliedCount int
    err = r.db.QueryRow("SELECT COUNT(*) FROM ulasan WHERE balasan IS NOT NULL AND balasan != ''").Scan(&repliedCount)
    if err != nil {
        return nil, err
    }
    stats["replied_count"] = repliedCount

    // Rating distribution
    ratingDist := make(map[int]int)
    for i := 1; i <= 5; i++ {
        var count int
        r.db.QueryRow("SELECT COUNT(*) FROM ulasan WHERE rating = $1 AND status = 'approved'", i).Scan(&count)
        ratingDist[i] = count
    }
    stats["rating_distribution"] = ratingDist

    return stats, nil
}

func getStringValue(s sql.NullString) string {
    if s.Valid {
        return s.String
    }
    return ""
}