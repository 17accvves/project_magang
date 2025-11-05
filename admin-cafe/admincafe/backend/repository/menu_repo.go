package repository

import (
    "database/sql"
    "backend/models"
    "time"
    "fmt"

    "github.com/google/uuid"
)

type MenuRepository struct {
    db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
    return &MenuRepository{db: db}
}

func (r *MenuRepository) GetAllMenus() ([]models.Menu, error) {
    query := `
        SELECT id, name, price, discount, discounted_price, start_date, end_date, 
               category, status, img, created_at, updated_at 
        FROM menus 
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var menus []models.Menu
    for rows.Next() {
        var menu models.Menu
        var startDate, endDate sql.NullString
        
        err := rows.Scan(
            &menu.ID,
            &menu.Name,
            &menu.Price,
            &menu.Discount,
            &menu.DiscountedPrice,
            &startDate,
            &endDate,
            &menu.Category,
            &menu.Status,
            &menu.Img,
            &menu.CreatedAt,
            &menu.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        if startDate.Valid {
            menu.StartDate = startDate.String
        }
        if endDate.Valid {
            menu.EndDate = endDate.String
        }

        menus = append(menus, menu)
    }

    return menus, nil
}

func (r *MenuRepository) GetMenuByID(id string) (*models.Menu, error) {
    query := `
        SELECT id, name, price, discount, discounted_price, start_date, end_date, 
               category, status, img, created_at, updated_at 
        FROM menus 
        WHERE id = $1
    `
    
    var menu models.Menu
    var startDate, endDate sql.NullString
    
    err := r.db.QueryRow(query, id).Scan(
        &menu.ID,
        &menu.Name,
        &menu.Price,
        &menu.Discount,
        &menu.DiscountedPrice,
        &startDate,
        &endDate,
        &menu.Category,
        &menu.Status,
        &menu.Img,
        &menu.CreatedAt,
        &menu.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    if startDate.Valid {
        menu.StartDate = startDate.String
    }
    if endDate.Valid {
        menu.EndDate = endDate.String
    }

    return &menu, nil
}

func (r *MenuRepository) CreateMenu(menu *models.Menu) error {
    menu.ID = uuid.New().String()
    menu.CreatedAt = time.Now()
    menu.UpdatedAt = time.Now()

    query := `
        INSERT INTO menus (id, name, price, discount, discounted_price, start_date, end_date, category, status, img, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

    _, err := r.db.Exec(
        query,
        menu.ID,
        menu.Name,
        menu.Price,
        menu.Discount,
        menu.DiscountedPrice,
        nullIfEmpty(menu.StartDate),
        nullIfEmpty(menu.EndDate),
        menu.Category,
        menu.Status,
        menu.Img,
        menu.CreatedAt,
        menu.UpdatedAt,
    )

    return err
}

func (r *MenuRepository) UpdateMenu(id string, updateData map[string]interface{}) error {
    updateData["updated_at"] = time.Now()

    query := "UPDATE menus SET "
    params := []interface{}{}
    paramCount := 1

    for key, value := range updateData {
        // Handle field mapping dari JSON ke database column
        dbColumn := key
        switch key {
        case "startDate":
            dbColumn = "start_date"
        case "endDate":
            dbColumn = "end_date"
        case "discountedPrice":
            dbColumn = "discounted_price"
        }

        query += dbColumn + " = $" + fmt.Sprint(paramCount) + ", "
        params = append(params, value)
        paramCount++
    }

    // Remove trailing comma and space
    query = query[:len(query)-2]
    query += " WHERE id = $" + fmt.Sprint(paramCount)
    params = append(params, id)

    _, err := r.db.Exec(query, params...)
    return err
}

func (r *MenuRepository) DeleteMenu(id string) error {
    query := "DELETE FROM menus WHERE id = $1"
    _, err := r.db.Exec(query, id)
    return err
}

func nullIfEmpty(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}