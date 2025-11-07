package repository

import (
	"backend/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MenuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// ===============================
// FUNGSI UNTUK PROMO (MENU DENGAN DISKON)
// ===============================

// GetMenusWithDiscount - Ambil menu yang memiliki diskon (untuk promo)
func (r *MenuRepository) GetMenusWithDiscount() ([]models.Menu, error) {
	query := `
		SELECT id, name, price, discount, discounted_price, start_date, end_date, 
			   category, status, img, created_at, updated_at
		FROM menus 
		WHERE discount > 0 AND status = 'Aktif'
		ORDER BY discount DESC, created_at DESC
	`

	fmt.Printf("🔍 [DEBUG PROMO] Executing query for menus with discount...\n")

	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Printf("❌ [DEBUG PROMO] Query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	count := 0
	today := time.Now().Format("2006-01-02")
	
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
			fmt.Printf("❌ [DEBUG PROMO] Scan failed: %v\n", err)
			return nil, err
		}

		// Handle null values
		if startDate.Valid {
			menu.StartDate = startDate.String
		} else {
			menu.StartDate = ""
		}
		if endDate.Valid {
			menu.EndDate = endDate.String
		} else {
			menu.EndDate = ""
		}

		// ✅ FITUR BARU: Auto-check tanggal diskon expired
		if menu.Discount > 0 && menu.EndDate != "" {
			if menu.EndDate < today {
				// Diskon sudah berakhir, kembalikan ke harga normal
				fmt.Printf("🔍 [AUTO-RESET] Discount expired for menu %s, resetting to normal price\n", menu.Name)
				menu.DiscountedPrice = menu.Price
				menu.Discount = 0
				menu.StartDate = ""
				menu.EndDate = ""
				
				// Update di database juga (background)
				go r.resetExpiredDiscount(menu.ID)
			}
		}

		count++
		fmt.Printf("✅ [DEBUG PROMO] Menu %d: %s - Discount: %.0f%%, Price: Rp %.0f\n", 
			count, menu.Name, menu.Discount, menu.Price)
		
		menus = append(menus, menu)
	}

	fmt.Printf("🎯 [DEBUG PROMO] TOTAL FOUND: %d menus with discount\n", count)
	
	if count == 0 {
		fmt.Printf("⚠️  [DEBUG PROMO] No menus found with discount. Checking database...\n")
		
		// Debug: Check what's actually in the database
		checkQuery := `SELECT COUNT(*) as total, 
							  COUNT(CASE WHEN discount > 0 THEN 1 END) as with_discount,
							  COUNT(CASE WHEN status = 'Aktif' THEN 1 END) as active
					   FROM menus`
		var total, withDiscount, active int
		r.db.QueryRow(checkQuery).Scan(&total, &withDiscount, &active)
		fmt.Printf("📊 [DEBUG PROMO] Database stats - Total: %d, With Discount: %d, Active: %d\n", 
			total, withDiscount, active)
	}

	return menus, nil
}

// GetPromoStats - Get statistik promo dari menu dengan diskon
func (r *MenuRepository) GetPromoStats() (*models.PromoStats, error) {
	menus, err := r.GetMenusWithDiscount()
	if err != nil {
		return nil, err
	}

	totalRevenue := 0
	var promosList []string

	for _, menu := range menus {
		// Simulasi revenue berdasarkan harga dan diskon
		simulatedRevenue := (menu.Price * menu.Discount) / 100
		if simulatedRevenue > 0 {
			totalRevenue += int(simulatedRevenue / 10000) // Convert to JT
		}
		promosList = append(promosList, menu.Name)
	}

	stats := &models.PromoStats{
		ActivePromos: len(menus),
		TotalRevenue: fmt.Sprintf("%dJT", totalRevenue), // ✅ PERBAIKAN: format revenue
		PromosList:   promosList,
	}

	return stats, nil
}

// ===============================
// FUNGSI UTAMA MENU (YANG SUDAH ADA)
// ===============================

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
	today := time.Now().Format("2006-01-02")
	
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

		// Handle null values
		if startDate.Valid {
			menu.StartDate = startDate.String
		} else {
			menu.StartDate = ""
		}
		if endDate.Valid {
			menu.EndDate = endDate.String
		} else {
			menu.EndDate = ""
		}

		// ✅ FITUR BARU: Auto-check tanggal diskon expired
		if menu.Discount > 0 && menu.EndDate != "" {
			if menu.EndDate < today {
				// Diskon sudah berakhir, kembalikan ke harga normal
				fmt.Printf("🔍 [AUTO-RESET] Discount expired for menu %s, resetting to normal price\n", menu.Name)
				menu.DiscountedPrice = menu.Price
				menu.Discount = 0
				menu.StartDate = ""
				menu.EndDate = ""
				
				// Update di database juga (background)
				go r.resetExpiredDiscount(menu.ID)
			}
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

	// Handle null values
	if startDate.Valid {
		menu.StartDate = startDate.String
	} else {
		menu.StartDate = ""
	}
	if endDate.Valid {
		menu.EndDate = endDate.String
	} else {
		menu.EndDate = ""
	}

	// ✅ FITUR BARU: Auto-check tanggal diskon expired untuk single menu
	today := time.Now().Format("2006-01-02")
	if menu.Discount > 0 && menu.EndDate != "" {
		if menu.EndDate < today {
			// Diskon sudah berakhir, kembalikan ke harga normal
			fmt.Printf("🔍 [AUTO-RESET] Discount expired for menu %s, resetting to normal price\n", menu.Name)
			menu.DiscountedPrice = menu.Price
			menu.Discount = 0
			menu.StartDate = ""
			menu.EndDate = ""
			
			// Update di database juga (background)
			go r.resetExpiredDiscount(menu.ID)
		}
	}

	return &menu, nil
}

func (r *MenuRepository) CreateMenu(menu *models.Menu) error {
	// Pastikan diskon tidak negatif
	if menu.Discount < 0 {
		menu.Discount = 0
	}

	// Validasi diskon dan tanggal
	if menu.Discount > 0 {
		// Jika ada diskon, tanggal harus diisi
		if menu.StartDate == "" {
			return fmt.Errorf("start date is required when discount is set")
		}
		if menu.EndDate == "" {
			return fmt.Errorf("end date is required when discount is set")
		}

		// Validasi format tanggal
		if _, err := time.Parse("2006-01-02", menu.StartDate); err != nil {
			return fmt.Errorf("invalid start date format, use YYYY-MM-DD")
		}
		if _, err := time.Parse("2006-01-02", menu.EndDate); err != nil {
			return fmt.Errorf("invalid end date format, use YYYY-MM-DD")
		}

		// Validasi tanggal tidak boleh berlalu
		today := time.Now().Format("2006-01-02")
		if menu.EndDate < today {
			return fmt.Errorf("end date cannot be in the past")
		}

		// Validasi start date tidak boleh setelah end date
		if menu.StartDate > menu.EndDate {
			return fmt.Errorf("start date cannot be after end date")
		}
	} else {
		// Jika tidak ada diskon, pastikan tanggal kosong
		menu.Discount = 0
		menu.StartDate = ""
		menu.EndDate = ""
	}

	// Hitung discounted_price
	if menu.Discount > 0 {
		menu.DiscountedPrice = menu.Price - (menu.Price * menu.Discount / 100)
	} else {
		menu.DiscountedPrice = menu.Price
	}

	menu.ID = uuid.New().String()
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()

	query := `
        INSERT INTO menus (id, name, price, discount, discounted_price, start_date, end_date, category, status, img, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `

	// ✅ SEKARANG SIMPLE: Langsung kirim string (bisa kosong karena kolom VARCHAR)
	fmt.Printf("🔍 [DEBUG] Creating menu - Discount: %.2f, StartDate: '%s', EndDate: '%s'\n", 
		menu.Discount, menu.StartDate, menu.EndDate)

	_, err := r.db.Exec(
		query,
		menu.ID,              // $1
		menu.Name,            // $2
		menu.Price,           // $3
		menu.Discount,        // $4
		menu.DiscountedPrice, // $5
		menu.StartDate,       // $6 - bisa string kosong "" (VARCHAR)
		menu.EndDate,         // $7 - bisa string kosong "" (VARCHAR)
		menu.Category,        // $8
		menu.Status,          // $9
		menu.Img,             // $10
		menu.CreatedAt,       // $11
		menu.UpdatedAt,       // $12
	)

	if err != nil {
		fmt.Printf("❌ [ERROR] CreateMenu failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ [SUCCESS] Menu created: %s\n", menu.ID)
	return nil
}

func (r *MenuRepository) UpdateMenu(id string, updateData map[string]interface{}) error {
	fmt.Printf("🔍 [DEBUG] UpdateMenu - ID: %s, Data: %+v\n", id, updateData)

	// ✅ PERBAIKAN: Debug semua field yang ada SEBELUM validasi
	fmt.Printf("🔍 [DEBUG] === ALL FIELDS IN updateData ===\n")
	for key, value := range updateData {
		fmt.Printf("  %s: %v (type: %T)\n", key, value, value)
	}
	fmt.Printf("🔍 [DEBUG] ================================\n")

	// ✅ PERBAIKAN: Handle discounted_price calculation FIRST untuk menghindari duplicate
	if discount, hasDiscount := updateData["discount"]; hasDiscount {
		discountFloat, ok := discount.(float64)
		if !ok {
			return fmt.Errorf("discount must be a number")
		}

		// Pastikan diskon tidak negatif
		if discountFloat < 0 {
			discountFloat = 0
			updateData["discount"] = 0.0
		}

		// Hitung discounted_price
		if price, hasPrice := updateData["price"]; hasPrice {
			priceFloat, ok := price.(float64)
			if ok {
				if discountFloat > 0 {
					updateData["discounted_price"] = priceFloat - (priceFloat * discountFloat / 100)
				} else {
					updateData["discounted_price"] = priceFloat
				}
			}
		} else if discountFloat > 0 {
			// Jika hanya discount yang berubah, hitung berdasarkan harga current
			currentMenu, err := r.GetMenuByID(id)
			if err == nil && currentMenu != nil {
				updateData["discounted_price"] = currentMenu.Price - (currentMenu.Price * discountFloat / 100)
			}
		}

		// ✅ PERBAIKAN: Hapus discountedPrice untuk menghindari duplicate
		delete(updateData, "discountedPrice")
	}

	// Validasi diskon dan tanggal jika ada field diskon yang diupdate
	if discount, hasDiscount := updateData["discount"]; hasDiscount {
		discountFloat, ok := discount.(float64)
		if !ok {
			return fmt.Errorf("discount must be a number")
		}

		// ✅ PERBAIKAN: Handle kedua format field dengan cara yang lebih robust
		var startDate, endDate string
		var hasStart, hasEnd bool

		// Cari field start_date (prioritas) atau startDate
		if startDateVal, exists := updateData["start_date"]; exists {
			if startDateStr, ok := startDateVal.(string); ok && startDateStr != "" {
				startDate = startDateStr
				hasStart = true
				fmt.Printf("🔍 [DEBUG] ✅ Found start_date: '%s'\n", startDate)
			} else {
				fmt.Printf("🔍 [DEBUG] ❌ start_date exists but empty/invalid: %v\n", startDateVal)
			}
		} else {
			fmt.Printf("🔍 [DEBUG] ❌ start_date field not found\n")
		}
		
		// Jika start_date tidak ada atau kosong, coba startDate
		if !hasStart {
			if startDateVal, exists := updateData["startDate"]; exists {
				if startDateStr, ok := startDateVal.(string); ok && startDateStr != "" {
					startDate = startDateStr
					hasStart = true
					fmt.Printf("🔍 [DEBUG] ✅ Found startDate: '%s', mapping to start_date\n", startDate)
					// Mapping ke format database
					updateData["start_date"] = startDate
					delete(updateData, "startDate")
				} else {
					fmt.Printf("🔍 [DEBUG] ❌ startDate exists but empty/invalid: %v\n", startDateVal)
				}
			} else {
				fmt.Printf("🔍 [DEBUG] ❌ startDate field not found\n")
			}
		}

		// Cari field end_date (prioritas) atau endDate
		if endDateVal, exists := updateData["end_date"]; exists {
			if endDateStr, ok := endDateVal.(string); ok && endDateStr != "" {
				endDate = endDateStr
				hasEnd = true
				fmt.Printf("🔍 [DEBUG] ✅ Found end_date: '%s'\n", endDate)
			} else {
				fmt.Printf("🔍 [DEBUG] ❌ end_date exists but empty/invalid: %v\n", endDateVal)
			}
		} else {
			fmt.Printf("🔍 [DEBUG] ❌ end_date field not found\n")
		}
		
		// Jika end_date tidak ada atau kosong, coba endDate
		if !hasEnd {
			if endDateVal, exists := updateData["endDate"]; exists {
				if endDateStr, ok := endDateVal.(string); ok && endDateStr != "" {
					endDate = endDateStr
					hasEnd = true
					fmt.Printf("🔍 [DEBUG] ✅ Found endDate: '%s', mapping to end_date\n", endDate)
					// Mapping ke format database
					updateData["end_date"] = endDate
					delete(updateData, "endDate")
				} else {
					fmt.Printf("🔍 [DEBUG] ❌ endDate exists but empty/invalid: %v\n", endDateVal)
				}
			} else {
				fmt.Printf("🔍 [DEBUG] ❌ endDate field not found\n")
			}
		}

		fmt.Printf("🔍 [DEBUG] === FINAL RESULT ===\n")
		fmt.Printf("🔍 [DEBUG] Discount: %.2f\n", discountFloat)
		fmt.Printf("🔍 [DEBUG] StartDate: '%s' (found: %t)\n", startDate, hasStart)
		fmt.Printf("🔍 [DEBUG] EndDate: '%s' (found: %t)\n", endDate, hasEnd)
		fmt.Printf("🔍 [DEBUG] ====================\n")

		// Jika ada diskon > 0, validasi tanggal harus diisi
		if discountFloat > 0 {
			if !hasStart || startDate == "" {
				return fmt.Errorf("start date is required when discount is set")
			}
			if !hasEnd || endDate == "" {
				return fmt.Errorf("end date is required when discount is set")
			}

			// Validasi format tanggal
			if _, err := time.Parse("2006-01-02", startDate); err != nil {
				return fmt.Errorf("invalid start date format, use YYYY-MM-DD")
			}
			if _, err := time.Parse("2006-01-02", endDate); err != nil {
				return fmt.Errorf("invalid end date format, use YYYY-MM-DD")
			}

			// Validasi tanggal tidak boleh berlalu
			today := time.Now().Format("2006-01-02")
			if endDate < today {
				return fmt.Errorf("end date cannot be in the past")
			}

			// Validasi start date tidak boleh setelah end date
			if startDate > endDate {
				return fmt.Errorf("start date cannot be after end date")
			}
		} else {
			// Jika diskon = 0, kosongkan tanggal
			updateData["start_date"] = ""
			updateData["end_date"] = ""
			// Hapus field lama jika ada
			delete(updateData, "startDate")
			delete(updateData, "endDate")
		}
	}

	updateData["updated_at"] = time.Now()

	// Build query - hanya update field yang perlu
	query := "UPDATE menus SET "
	params := []interface{}{}
	paramCount := 1

	// ✅ PERBAIKAN: Skip field yang tidak perlu diupdate ke database
	excludedFields := map[string]bool{
		"id":         true,
		"createdAt":  true,
		"created_at": true,
		"updatedAt":  true,
		"discountedPrice": true, // ✅ PERBAIKAN: Exclude untuk hindari duplicate
	}

	for key, value := range updateData {
		if excludedFields[key] {
			continue
		}

		// Handle field name mapping untuk konsistensi database
		dbColumn := key
		switch key {
		case "discountedPrice":
			dbColumn = "discounted_price"
		}

		query += dbColumn + " = $" + fmt.Sprint(paramCount) + ", "
		params = append(params, value)
		paramCount++
	}

	// Remove trailing comma and add WHERE clause
	if len(params) > 0 {
		query = query[:len(query)-2]
	} else {
		return fmt.Errorf("no fields to update")
	}

	query += " WHERE id = $" + fmt.Sprint(paramCount)
	params = append(params, id)

	fmt.Printf("🔍 [DEBUG] Executing: %s\n", query)
	fmt.Printf("🔍 [DEBUG] Params: %+v\n", params)

	result, err := r.db.Exec(query, params...)
	if err != nil {
		fmt.Printf("❌ [ERROR] UpdateMenu failed: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("menu not found")
	}

	fmt.Printf("✅ [SUCCESS] Menu updated: %s\n", id)
	return nil
}

func (r *MenuRepository) DeleteMenu(id string) error {
	query := "DELETE FROM menus WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Printf("❌ [ERROR] DeleteMenu failed: %v\n", err)
		return err
	}
	fmt.Printf("✅ [SUCCESS] Menu deleted: %s\n", id)
	return nil
}

// ===============================
// FUNGSI TAMBAHAN MENU
// ===============================

// GetMenusByCategory untuk filter berdasarkan kategori
func (r *MenuRepository) GetMenusByCategory(category string) ([]models.Menu, error) {
	query := `
        SELECT id, name, price, discount, discounted_price, start_date, end_date, 
               category, status, img, created_at, updated_at 
        FROM menus 
        WHERE category = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	today := time.Now().Format("2006-01-02")
	
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
		} else {
			menu.StartDate = ""
		}
		if endDate.Valid {
			menu.EndDate = endDate.String
		} else {
			menu.EndDate = ""
		}

		// ✅ FITUR BARU: Auto-check tanggal diskon expired
		if menu.Discount > 0 && menu.EndDate != "" {
			if menu.EndDate < today {
				// Diskon sudah berakhir, kembalikan ke harga normal
				fmt.Printf("🔍 [AUTO-RESET] Discount expired for menu %s, resetting to normal price\n", menu.Name)
				menu.DiscountedPrice = menu.Price
				menu.Discount = 0
				menu.StartDate = ""
				menu.EndDate = ""
				
				// Update di database juga (background)
				go r.resetExpiredDiscount(menu.ID)
			}
		}

		menus = append(menus, menu)
	}

	return menus, nil
}

// SearchMenus untuk pencarian menu
func (r *MenuRepository) SearchMenus(query string) ([]models.Menu, error) {
	sqlQuery := `
        SELECT id, name, price, discount, discounted_price, start_date, end_date, 
               category, status, img, created_at, updated_at 
        FROM menus 
        WHERE name ILIKE $1 OR category ILIKE $1
        ORDER BY created_at DESC
    `

	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	today := time.Now().Format("2006-01-02")
	
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
		} else {
			menu.StartDate = ""
		}
		if endDate.Valid {
			menu.EndDate = endDate.String
		} else {
			menu.EndDate = ""
		}

		// ✅ FITUR BARU: Auto-check tanggal diskon expired
		if menu.Discount > 0 && menu.EndDate != "" {
			if menu.EndDate < today {
				// Diskon sudah berakhir, kembalikan ke harga normal
				fmt.Printf("🔍 [AUTO-RESET] Discount expired for menu %s, resetting to normal price\n", menu.Name)
				menu.DiscountedPrice = menu.Price
				menu.Discount = 0
				menu.StartDate = ""
				menu.EndDate = ""
				
				// Update di database juga (background)
				go r.resetExpiredDiscount(menu.ID)
			}
		}

		menus = append(menus, menu)
	}

	return menus, nil
}

// CalculateDiscountedPrice menghitung harga diskon
func (r *MenuRepository) CalculateDiscountedPrice(price, discount float64) float64 {
	if discount > 0 {
		return price - (price * discount / 100)
	}
	return price
}

// ===============================
// FUNGSI UTILITAS DISKON
// ===============================

// ✅ FITUR BARU: resetExpiredDiscount reset diskon yang sudah expired di database
func (r *MenuRepository) resetExpiredDiscount(menuID string) error {
	query := `
        UPDATE menus 
        SET discount = 0, discounted_price = price, start_date = NULL, end_date = NULL, updated_at = $1
        WHERE id = $2 AND discount > 0 AND end_date IS NOT NULL AND end_date < $3
    `
    
    today := time.Now().Format("2006-01-02")
    result, err := r.db.Exec(query, time.Now(), menuID, today)
    if err != nil {
        fmt.Printf("❌ [ERROR] resetExpiredDiscount failed: %v\n", err)
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        fmt.Printf("✅ [AUTO-RESET] Reset expired discount for menu: %s\n", menuID)
    }
    
    return nil
}

// ✅ FITUR BARU: StartDiscountChecker menjalankan checker diskon expired setiap hari
func (r *MenuRepository) StartDiscountChecker() {
    // Jalankan sekali saat startup
    fmt.Println("🔍 [SCHEDULER] Initial check for expired discounts...")
    r.checkAndResetAllExpiredDiscounts()
    
    // Schedule harian
    ticker := time.NewTicker(24 * time.Hour) // Check setiap 24 jam
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            fmt.Println("🔍 [SCHEDULER] Daily check for expired discounts...")
            r.checkAndResetAllExpiredDiscounts()
        }
    }
}

// ✅ FITUR BARU: checkAndResetAllExpiredDiscounts reset semua diskon yang sudah expired
func (r *MenuRepository) checkAndResetAllExpiredDiscounts() error {
    query := `
        UPDATE menus 
        SET discount = 0, discounted_price = price, start_date = NULL, end_date = NULL, updated_at = $1
        WHERE discount > 0 AND end_date IS NOT NULL AND end_date < $2
    `
    
    today := time.Now().Format("2006-01-02")
    result, err := r.db.Exec(query, time.Now(), today)
    if err != nil {
        fmt.Printf("❌ [ERROR] checkAndResetAllExpiredDiscounts failed: %v\n", err)
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        fmt.Printf("✅ [AUTO-RESET] Reset %d expired discounts\n", rowsAffected)
    }
    
    return nil
}