package repository

import (
	"backend/models"
	"database/sql"
	"fmt"
	"time"
)

type CafeProfileRepository struct {
	db *sql.DB
}

func NewCafeProfileRepository(db *sql.DB) *CafeProfileRepository {
	return &CafeProfileRepository{db: db}
}

func (r *CafeProfileRepository) GetCafeProfile() (*models.CafeProfile, error) {
	query := `SELECT id, nama, alamat, telepon, deskripsi, main_image, verified, created_at, updated_at 
	          FROM cafe_profiles WHERE id = 1`
	
	profile := &models.CafeProfile{}
	err := r.db.QueryRow(query).Scan(
		&profile.ID, &profile.Nama, &profile.Alamat, &profile.Telepon,
		&profile.Deskripsi, &profile.MainImage, &profile.Verified,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *CafeProfileRepository) UpdateCafeProfile(profile *models.CafeProfile) error {
	query := `UPDATE cafe_profiles 
	          SET nama = $1, alamat = $2, telepon = $3, deskripsi = $4, 
	              main_image = $5, verified = $6, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = 1`
	
	_, err := r.db.Exec(query, profile.Nama, profile.Alamat, profile.Telepon,
		profile.Deskripsi, profile.MainImage, profile.Verified)
	return err
}

// UpdateProfileImage - Update hanya gambar profile
func (r *CafeProfileRepository) UpdateProfileImage(imageURL string) error {
	query := `UPDATE cafe_profiles SET main_image = $1, updated_at = CURRENT_TIMESTAMP WHERE id = 1`
	_, err := r.db.Exec(query, imageURL)
	return err
}

func (r *CafeProfileRepository) GetSocialMedia() ([]models.CafeSocialMedia, error) {
	query := `SELECT id, cafe_profile_id, platform, url, created_at 
	          FROM cafe_social_media WHERE cafe_profile_id = 1`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var socialMedia []models.CafeSocialMedia
	for rows.Next() {
		var sm models.CafeSocialMedia
		err := rows.Scan(&sm.ID, &sm.CafeProfileID, &sm.Platform, &sm.URL, &sm.CreatedAt)
		if err != nil {
			return nil, err
		}
		socialMedia = append(socialMedia, sm)
	}
	return socialMedia, nil
}

func (r *CafeProfileRepository) AddSocialMedia(socialMedia *models.CafeSocialMedia) error {
	query := `INSERT INTO cafe_social_media (cafe_profile_id, platform, url) 
	          VALUES (1, $1, $2)`
	
	_, err := r.db.Exec(query, socialMedia.Platform, socialMedia.URL)
	return err
}

// DeleteAllSocialMedia - Hapus semua social media
func (r *CafeProfileRepository) DeleteAllSocialMedia() error {
	_, err := r.db.Exec("DELETE FROM cafe_social_media WHERE cafe_profile_id = 1")
	return err
}

func (r *CafeProfileRepository) GetOperationalHours() ([]models.CafeOperationalHours, error) {
	query := `SELECT id, cafe_profile_id, hari, buka, tutup, created_at 
	          FROM cafe_operational_hours WHERE cafe_profile_id = 1`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hours []models.CafeOperationalHours
	for rows.Next() {
		var hour models.CafeOperationalHours
		err := rows.Scan(&hour.ID, &hour.CafeProfileID, &hour.Hari, &hour.Buka, &hour.Tutup, &hour.CreatedAt)
		if err != nil {
			return nil, err
		}
		hours = append(hours, hour)
	}
	return hours, nil
}

func (r *CafeProfileRepository) UpdateOperationalHours(hours []models.CafeOperationalHours) error {
	// Delete existing hours
	_, err := r.db.Exec("DELETE FROM cafe_operational_hours WHERE cafe_profile_id = 1")
	if err != nil {
		return err
	}

	// Insert new hours
	for _, hour := range hours {
		query := `INSERT INTO cafe_operational_hours (cafe_profile_id, hari, buka, tutup) 
		          VALUES (1, $1, $2, $3)`
		_, err := r.db.Exec(query, hour.Hari, hour.Buka, hour.Tutup)
		if err != nil {
			return err
		}
	}
	return nil
}

// ✅ METHOD BARU: GetTodayOperationalHours - Ambil jam operasional hari ini (format 24 jam)
func (r *CafeProfileRepository) GetTodayOperationalHours() (*models.OpeningHoursResponse, error) {
	// Dapatkan nama hari dalam Bahasa Indonesia
	today := getTodayIndonesian()
	
	query := `
		SELECT buka, tutup 
		FROM cafe_operational_hours 
		WHERE hari = $1 AND cafe_profile_id = 1
	`
	
	var buka, tutup string
	err := r.db.QueryRow(query, today).Scan(&buka, &tutup)
	if err != nil {
		if err == sql.ErrNoRows {
			// Jika tidak ada data untuk hari ini, return default
			return &models.OpeningHoursResponse{
				Waktu:  "07:00 - 22:00",
				Status: "Buka Sekarang",
			}, nil
		}
		return nil, err
	}
	
	result := &models.OpeningHoursResponse{
		Waktu:  fmt.Sprintf("%s - %s", buka, tutup),
		Status: r.getSimpleStatus(buka, tutup),
	}
	
	fmt.Printf("🔍 [DEBUG] Hari: %s, Buka: %s, Tutup: %s, Status: %s\n", 
		today, buka, tutup, result.Status)
	
	return result, nil
}

// ✅ METHOD BARU: UpdateSingleOperationalHours - Update jam operasional untuk hari tertentu
func (r *CafeProfileRepository) UpdateSingleOperationalHours(hari, buka, tutup string) error {
	// Validasi format HH:MM
	if !isValidTimeFormat(buka) || !isValidTimeFormat(tutup) {
		return fmt.Errorf("invalid time format, use HH:MM (24-hour format)")
	}
	
	// Validasi buka tidak boleh setelah tutup
	if buka >= tutup {
		return fmt.Errorf("buka time cannot be after or equal to tutup time")
	}
	
	query := `
		INSERT INTO cafe_operational_hours (cafe_profile_id, hari, buka, tutup, created_at) 
		VALUES (1, $1, $2, $3, $4)
		ON CONFLICT (cafe_profile_id, hari) 
		DO UPDATE SET buka = $2, tutup = $3, created_at = $4
	`
	
	_, err := r.db.Exec(query, hari, buka, tutup, time.Now())
	if err != nil {
		fmt.Printf("❌ [ERROR] UpdateSingleOperationalHours failed: %v\n", err)
		return err
	}
	
	fmt.Printf("✅ [SUCCESS] Operational hours updated for %s: %s - %s\n", hari, buka, tutup)
	return nil
}

// ✅ METHOD BARU: GetCurrentStatus - Hanya ambil status buka/tutup
func (r *CafeProfileRepository) GetCurrentStatus() (*models.CurrentStatusResponse, error) {
	hours, err := r.GetTodayOperationalHours()
	if err != nil {
		return &models.CurrentStatusResponse{Status: "Tutup Sekarang"}, err
	}
	return &models.CurrentStatusResponse{Status: hours.Status}, nil
}

// ✅ HELPER FUNCTION: getSimpleStatus - Cek status buka/tutup berdasarkan waktu sekarang
func (r *CafeProfileRepository) getSimpleStatus(buka, tutup string) string {
	now := time.Now().Format("15:04") // Format 24 jam: "HH:MM"
	
	if now >= buka && now <= tutup {
		return "Buka Sekarang"
	}
	return "Tutup Sekarang"
}

func (r *CafeProfileRepository) GetFacilities() ([]models.CafeFacility, error) {
	query := `SELECT id, cafe_profile_id, nama_fasilitas, tersedia, created_at 
	          FROM cafe_facilities WHERE cafe_profile_id = 1`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facilities []models.CafeFacility
	for rows.Next() {
		var facility models.CafeFacility
		err := rows.Scan(&facility.ID, &facility.CafeProfileID, &facility.NamaFasilitas, 
			&facility.Tersedia, &facility.CreatedAt)
		if err != nil {
			return nil, err
		}
		facilities = append(facilities, facility)
	}
	return facilities, nil
}

func (r *CafeProfileRepository) UpdateFacilities(facilities []models.CafeFacility) error {
	// Delete existing facilities
	_, err := r.db.Exec("DELETE FROM cafe_facilities WHERE cafe_profile_id = 1")
	if err != nil {
		return err
	}

	// Insert new facilities
	for _, facility := range facilities {
		query := `INSERT INTO cafe_facilities (cafe_profile_id, nama_fasilitas, tersedia) 
		          VALUES (1, $1, $2)`
		_, err := r.db.Exec(query, facility.NamaFasilitas, facility.Tersedia)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CafeProfileRepository) GetGallery() ([]models.CafeGallery, error) {
	query := `SELECT id, cafe_profile_id, image_url, urutan, created_at 
	          FROM cafe_gallery WHERE cafe_profile_id = 1 ORDER BY urutan`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gallery []models.CafeGallery
	for rows.Next() {
		var item models.CafeGallery
		err := rows.Scan(&item.ID, &item.CafeProfileID, &item.ImageURL, &item.Urutan, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		gallery = append(gallery, item)
	}
	return gallery, nil
}

func (r *CafeProfileRepository) AddGalleryImage(gallery *models.CafeGallery) error {
	query := `INSERT INTO cafe_gallery (cafe_profile_id, image_url, urutan) 
	          VALUES (1, $1, $2) RETURNING id`
	
	err := r.db.QueryRow(query, gallery.ImageURL, gallery.Urutan).Scan(&gallery.ID)
	return err
}

// DeleteGalleryImage - Hapus gambar gallery
func (r *CafeProfileRepository) DeleteGalleryImage(id int) error {
	_, err := r.db.Exec("DELETE FROM cafe_gallery WHERE id = $1 AND cafe_profile_id = 1", id)
	return err
}

// ✅ HELPER FUNCTION: getTodayIndonesian - Mendapatkan hari dalam Bahasa Indonesia
func getTodayIndonesian() string {
	days := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}
	return days[time.Now().Weekday()]
}

// ✅ HELPER FUNCTION: isValidTimeFormat - Validasi format HH:MM
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	
	// Cek format HH:MM
	if timeStr[2] != ':' {
		return false
	}
	
	// Cek angka
	hour := timeStr[0:2]
	minute := timeStr[3:5]
	
	if hour < "00" || hour > "23" || minute < "00" || minute > "59" {
		return false
	}
	
	return true
}