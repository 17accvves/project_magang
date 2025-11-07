package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

// ==================== NULL HANDLING FUNCTIONS ====================

// ToNullString converts string to sql.NullString
func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// FromNullString converts sql.NullString to string
func FromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// ToNullInt64 converts int64 to sql.NullInt64
func ToNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

// FromNullInt64 converts sql.NullInt64 to int64
func FromNullInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// ToNullFloat64 converts float64 to sql.NullFloat64
func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// FromNullFloat64 converts sql.NullFloat64 to float64
func FromNullFloat64(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}

// ToNullBool converts bool to sql.NullBool
func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// FromNullBool converts sql.NullBool to bool
func FromNullBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}

// ==================== STRING UTILITY FUNCTIONS ====================

// NullIfEmpty converts empty string to nil for database optional fields
func NullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// IsEmpty checks if string is empty
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// GetStringOrDefault returns string value or default value if empty
func GetStringOrDefault(s, defaultValue string) string {
	if IsEmpty(s) {
		return defaultValue
	}
	return s
}

// GetStringOrEmpty returns string value or empty string if nil
func GetStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ==================== QUERY BUILDING FUNCTIONS ====================

// BuildQuery builds SQL query with conditions
func BuildQuery(baseQuery string, conditions []string) string {
	if len(conditions) == 0 {
		return baseQuery
	}
	
	query := baseQuery + " WHERE "
	for i, condition := range conditions {
		if i > 0 {
			query += " AND "
		}
		query += condition
	}
	return query
}

// BuildQueryWithOrder builds SQL query with conditions and order
func BuildQueryWithOrder(baseQuery string, conditions []string, orderBy string) string {
	query := BuildQuery(baseQuery, conditions)
	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}
	return query
}

// BuildPaginationQuery builds SQL query with pagination
func BuildPaginationQuery(baseQuery string, conditions []string, orderBy string, limit, offset int) string {
	query := BuildQueryWithOrder(baseQuery, conditions, orderBy)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}
	return query
}

// AddLikeCondition adds LIKE condition for search
func AddLikeCondition(field, searchTerm string) string {
	if IsEmpty(searchTerm) {
		return ""
	}
	return fmt.Sprintf("%s ILIKE '%%%s%%'", field, searchTerm)
}

// AddInCondition adds IN condition for multiple values
func AddInCondition(field string, values []string) string {
	if len(values) == 0 {
		return ""
	}
	
	quotedValues := make([]string, len(values))
	for i, v := range values {
		quotedValues[i] = fmt.Sprintf("'%s'", v)
	}
	
	return fmt.Sprintf("%s IN (%s)", field, strings.Join(quotedValues, ", "))
}

// ==================== VALIDATION FUNCTIONS ====================

// IsValidRating validates rating is between 1-5
func IsValidRating(rating int) bool {
	return rating >= 1 && rating <= 5
}

// IsValidStatus validates ulasan status
func IsValidStatus(status string) bool {
	validStatuses := []string{"pending", "approved", "rejected"}
	for _, s := range validStatuses {
		if strings.ToLower(status) == s {
			return true
		}
	}
	return false
}

// IsValidEmail basic email validation
func IsValidEmail(email string) bool {
	if IsEmpty(email) {
		return true // Email is optional
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ==================== FORMATTING FUNCTIONS ====================

// FormatRating formats rating for display
func FormatRating(rating int) string {
	if rating < 1 || rating > 5 {
		return "0/5"
	}
	return fmt.Sprintf("%d/5", rating)
}

// FormatStatus formats status for display
func FormatStatus(status string) string {
	switch strings.ToLower(status) {
	case "pending":
		return "Menunggu"
	case "approved":
		return "Disetujui"
	case "rejected":
		return "Ditolak"
	default:
		return status
	}
}

// TruncateText truncates text to specified length
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// ==================== ERROR HANDLING FUNCTIONS ====================

// HandleDBError handles common database errors
func HandleDBError(err error, operation string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("data tidak ditemukan")
	}
	return fmt.Errorf("gagal %s: %v", operation, err)
}

// IsDuplicateError checks if error is due to duplicate entry
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	// PostgreSQL duplicate error
	return strings.Contains(err.Error(), "duplicate key value") ||
		strings.Contains(err.Error(), "violates unique constraint")
}

// IsForeignKeyError checks if error is due to foreign key constraint
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "foreign key constraint") ||
		strings.Contains(err.Error(), "violates foreign key constraint")
}

// ==================== SLICE UTILITY FUNCTIONS ====================

// StringSliceContains checks if slice contains string
func StringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicates from string slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// FilterSlice filters string slice based on condition
func FilterSlice(slice []string, condition func(string) bool) []string {
	result := []string{}
	for _, item := range slice {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}