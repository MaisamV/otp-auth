package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// Exists checks if a user exists by phone number
func (r *UserRepository) Exists(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, phoneNumber.String()).Scan(&exists)
	if err != nil {
		return false, errors.NewInternalError("Failed to check user existence", err)
	}
	
	return exists, nil
}

// List retrieves users with pagination and optional search
func (r *UserRepository) List(ctx context.Context, offset, limit int, searchPhone string, searchDateFrom, searchDateTo string) ([]*entities.User, int64, error) {
	// Build query with optional filters
	var query string
	var countQuery string
	var args []interface{}
	argIndex := 1
	
	// Base query
	baseQuery := `
		SELECT id, phone_number, scope, created_at, updated_at
		FROM users
	`
	baseCountQuery := `SELECT COUNT(*) FROM users`
	
	// Build WHERE conditions
	var conditions []string
	
	if searchPhone != "" {
		conditions = append(conditions, "phone_number ILIKE $"+fmt.Sprintf("%d", argIndex))
		args = append(args, "%"+searchPhone+"%")
		argIndex++
	}
	
	if searchDateFrom != "" {
		conditions = append(conditions, "created_at >= $"+fmt.Sprintf("%d", argIndex))
		args = append(args, searchDateFrom)
		argIndex++
	}
	
	if searchDateTo != "" {
		conditions = append(conditions, "created_at <= $"+fmt.Sprintf("%d", argIndex))
		args = append(args, searchDateTo)
		argIndex++
	}
	
	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query = baseQuery + whereClause
		countQuery = baseCountQuery + whereClause
	} else {
		query = baseQuery
		countQuery = baseCountQuery
	}
	
	// Add ORDER BY and pagination
	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, limit, offset)
	
	// Get total count
	var total int64
	countArgs := args[:len(args)-2] // Remove limit and offset for count query
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.NewInternalError("Failed to count users", err)
	}
	
	// Get users
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.NewInternalError("Failed to query users", err)
	}
	defer rows.Close()
	
	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		var phoneStr string
		
		err := rows.Scan(
			&user.ID,
			&phoneStr,
			&user.Scope,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.NewInternalError("Failed to scan user", err)
		}
		
		// Parse phone number
		phoneNumber, err := valueobjects.NewPhoneNumber(phoneStr)
		if err != nil {
			return nil, 0, errors.NewValidationError("Failed to parse phone number", err)
		}
		user.PhoneNumber = phoneNumber
		
		users = append(users, user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, 0, errors.NewInternalError("Failed to iterate users", err)
	}
	
	return users, total, nil
}