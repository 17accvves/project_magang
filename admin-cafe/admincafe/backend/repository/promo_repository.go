package repository

import (
	"backend/config"
	"backend/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PromoRepository struct {
	db *sql.DB
}

func NewPromoRepository() *PromoRepository {
	return &PromoRepository{db: config.DB}
}

func (r *PromoRepository) GetAllPromos() ([]models.Promo, error) {
	query := `SELECT id, title, start_date, end_date, revenue, status, created_at, updated_at 
	          FROM promos 
	          ORDER BY start_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promos []models.Promo
	for rows.Next() {
		var promo models.Promo
		var startDate, endDate sql.NullTime

		err := rows.Scan(
			&promo.ID, &promo.Title, &startDate, &endDate, 
			&promo.Revenue, &promo.Status, &promo.CreatedAt, &promo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Format tanggal
		if startDate.Valid {
			promo.StartDate = formatDate(startDate.Time)
		}
		if endDate.Valid {
			promo.EndDate = formatDate(endDate.Time)
		}

		promos = append(promos, promo)
	}

	return promos, nil
}

func (r *PromoRepository) GetPromoStats() (*models.PromoStats, error) {
	promos, err := r.GetAllPromos()
	if err != nil {
		return nil, err
	}

	totalRevenue := 0
	var promosList []string

	for _, promo := range promos {
		// Hitung revenue (hilangkan "JT" dan konversi ke integer)
		revenueStr := strings.TrimSuffix(promo.Revenue, "JT")
		if revenueNum, err := strconv.Atoi(strings.TrimSpace(revenueStr)); err == nil {
			totalRevenue += revenueNum
		}
		promosList = append(promosList, promo.Title)
	}

	stats := &models.PromoStats{
		ActivePromos: len(promos),
		TotalRevenue: fmt.Sprintf("%dJT", totalRevenue),
		PromosList:   promosList,
	}

	return stats, nil
}

func (r *PromoRepository) GetPromoByID(id int) (*models.Promo, error) {
	query := `SELECT id, title, start_date, end_date, revenue, status, created_at, updated_at 
	          FROM promos WHERE id = $1`

	var promo models.Promo
	var startDate, endDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&promo.ID, &promo.Title, &startDate, &endDate, 
		&promo.Revenue, &promo.Status, &promo.CreatedAt, &promo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Format tanggal
	if startDate.Valid {
		promo.StartDate = formatDate(startDate.Time)
	}
	if endDate.Valid {
		promo.EndDate = formatDate(endDate.Time)
	}

	return &promo, nil
}

func (r *PromoRepository) CreatePromo(promo *models.PromoRequest) error {
	query := `INSERT INTO promos (title, start_date, end_date, revenue, status) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	
	err := r.db.QueryRow(
		query, 
		promo.Title, promo.StartDate, promo.EndDate, promo.Revenue, promo.Status,
	).Scan(&id, &createdAt, &updatedAt)

	return err
}

func (r *PromoRepository) UpdatePromo(id int, promo *models.PromoRequest) error {
	query := `UPDATE promos 
	          SET title = $1, start_date = $2, end_date = $3, revenue = $4, status = $5, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $6`

	result, err := r.db.Exec(
		query, 
		promo.Title, promo.StartDate, promo.EndDate, promo.Revenue, promo.Status, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PromoRepository) DeletePromo(id int) error {
	query := "DELETE FROM promos WHERE id = $1"
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Helper function untuk format tanggal
func formatDate(date time.Time) string {
	// Indonesian month names
	months := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	day := date.Day()
	month := months[date.Month()]
	year := date.Year()

	return fmt.Sprintf("%d %s %d", day, month, year)
}