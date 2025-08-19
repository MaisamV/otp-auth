package postgres

import (
	"context"
	"database/sql"
	"github.com/otp-auth/internal/application/ports/repositories"
	"time"

	"github.com/lib/pq"
	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// UserRepository implements the user repository using PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, phone_number, scope, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Scope,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("User not found", nil)
		}
		return nil, errors.NewInternalError("Failed to get user by ID", err)
	}

	return &user, nil
}

// GetByPhoneNumber retrieves a user by phone number
func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber valueobjects.PhoneNumber) (*entities.User, error) {
	query := `
		SELECT id, phone_number, scope, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, string(phoneNumber)).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Scope,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("User not found", nil)
		}
		return nil, errors.NewInternalError("Failed to get user by phone number", err)
	}

	return &user, nil
}

// GetAll retrieves all users with pagination
func (r *UserRepository) GetAll(ctx context.Context, limit, offset int, scope string) ([]*entities.User, int, error) {
	// Build query with optional scope filter
	var query string
	var countQuery string
	var args []interface{}

	if scope != "" {
		query = `
			SELECT id, phone_number, scope, created_at, updated_at
			FROM users
			WHERE scope = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		countQuery = `SELECT COUNT(*) FROM users WHERE scope = $1`
		args = []interface{}{scope, limit, offset}
	} else {
		query = `
			SELECT id, phone_number, scope, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		countQuery = `SELECT COUNT(*) FROM users`
		args = []interface{}{limit, offset}
	}

	// Get total count
	var total int
	var countArgs []interface{}
	if scope != "" {
		countArgs = []interface{}{scope}
	}
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, errors.NewInternalError("Failed to count users", err)
	}

	// Get users
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.NewInternalError("Failed to get users", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID,
			&user.PhoneNumber,
			&user.Scope,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, errors.NewInternalError("Failed to scan user", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, errors.NewInternalError("Error iterating users", err)
	}

	return users, total, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, phone_number, scope, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.PhoneNumber,
		user.Scope,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return errors.NewValidationError("User with this phone number already exists", err)
			}
		}
		return errors.NewInternalError("Failed to create user", err)
	}

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET phone_number = $2, scope = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.PhoneNumber,
		user.Scope,
		user.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return errors.NewValidationError("User with this phone number already exists", err)
			}
		}
		return errors.NewInternalError("Failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("User not found", nil)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewInternalError("Failed to delete user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("User not found", nil)
	}

	return nil
}

// UpdateScope updates a user's scope
func (r *UserRepository) UpdateScope(ctx context.Context, id, scope string) error {
	query := `
		UPDATE users
		SET scope = $2, updated_at = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, scope, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to update user scope", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("User not found", nil)
	}

	return nil
}
