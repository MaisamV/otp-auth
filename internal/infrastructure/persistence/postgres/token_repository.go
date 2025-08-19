package postgres

import (
	"context"
	"database/sql"
	"github.com/otp-auth/internal/application/ports/repositories"
	"time"

	"github.com/otp-auth/internal/domain/entities"
	"github.com/otp-auth/internal/domain/valueobjects"
	"github.com/otp-auth/pkg/errors"
)

// TokenRepository implements the token repository using PostgreSQL
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository creates a new PostgreSQL token repository
func NewTokenRepository(db *sql.DB) repositories.TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

// GetByID retrieves a refresh token by ID
func (r *TokenRepository) GetByID(ctx context.Context, id string) (*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE id = $1
	`

	var token entities.RefreshToken
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.SessionID,
		&token.TokenHash,
		&token.CreatedAt,
		&token.ExpiresAt,
		&revokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Refresh token not found", nil)
		}
		return nil, errors.NewInternalError("Failed to get refresh token by ID", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
		token.Revoked = true
	}

	return &token, nil
}

// GetByTokenHash retrieves a refresh token by token hash
func (r *TokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token entities.RefreshToken
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.SessionID,
		&token.TokenHash,
		&token.CreatedAt,
		&token.ExpiresAt,
		&revokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Refresh token not found", nil)
		}
		return nil, errors.NewInternalError("Failed to get refresh token by hash", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
		token.Revoked = true
	}

	return &token, nil
}

// GetByTokenHashAndSessionID retrieves a refresh token by token hash and session ID
func (r *TokenRepository) GetByTokenHashAndSessionID(ctx context.Context, tokenHash string, sessionID valueobjects.SessionID) (*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND session_id = $2
	`

	var token entities.RefreshToken
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tokenHash, sessionID.String()).Scan(
		&token.ID,
		&token.UserID,
		&token.SessionID,
		&token.TokenHash,
		&token.CreatedAt,
		&token.ExpiresAt,
		&revokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Refresh token not found", nil)
		}
		return nil, errors.NewInternalError("Failed to get refresh token by hash and session ID", err)
	}

	if revokedAt.Valid {
		token.RevokedAt = &revokedAt.Time
		token.Revoked = true
	}

	return &token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *TokenRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.NewInternalError("Failed to get refresh tokens by user ID", err)
	}
	defer rows.Close()

	var tokens []*entities.RefreshToken
	for rows.Next() {
		var token entities.RefreshToken
		var revokedAt sql.NullTime

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.SessionID,
			&token.TokenHash,
			&token.CreatedAt,
			&token.ExpiresAt,
			&revokedAt,
		)
		if err != nil {
			return nil, errors.NewInternalError("Failed to scan refresh token", err)
		}

		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
			token.Revoked = true
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewInternalError("Error iterating refresh tokens", err)
	}

	return tokens, nil
}

// GetActiveByUserID retrieves all active (non-revoked, non-expired) refresh tokens for a user
func (r *TokenRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.NewInternalError("Failed to get active refresh tokens by user ID", err)
	}
	defer rows.Close()

	var tokens []*entities.RefreshToken
	for rows.Next() {
		var token entities.RefreshToken
		var revokedAt sql.NullTime

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.SessionID,
			&token.TokenHash,
			&token.CreatedAt,
			&token.ExpiresAt,
			&revokedAt,
		)
		if err != nil {
			return nil, errors.NewInternalError("Failed to scan active refresh token", err)
		}

		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
			token.Revoked = true
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewInternalError("Error iterating active refresh tokens", err)
	}

	return tokens, nil
}

// GetBySessionID retrieves all refresh tokens for a session
func (r *TokenRepository) GetBySessionID(ctx context.Context, sessionID string) ([]*entities.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, created_at, expires_at, revoked_at
		FROM refresh_tokens
		WHERE session_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, errors.NewInternalError("Failed to get refresh tokens by session ID", err)
	}
	defer rows.Close()

	var tokens []*entities.RefreshToken
	for rows.Next() {
		var token entities.RefreshToken
		var revokedAt sql.NullTime

		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.SessionID,
			&token.TokenHash,
			&token.CreatedAt,
			&token.ExpiresAt,
			&revokedAt,
		)
		if err != nil {
			return nil, errors.NewInternalError("Failed to scan refresh token", err)
		}

		if revokedAt.Valid {
			token.RevokedAt = &revokedAt.Time
			token.Revoked = true
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewInternalError("Error iterating refresh tokens", err)
	}

	return tokens, nil
}

// Create creates a new refresh token
func (r *TokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, session_id, token_hash, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.SessionID,
		token.TokenHash,
		token.CreatedAt,
		token.ExpiresAt,
	)

	if err != nil {
		return errors.NewInternalError("Failed to create refresh token", err)
	}

	return nil
}

// Update updates an existing refresh token
func (r *TokenRepository) Update(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET user_id = $2, session_id = $3, token_hash = $4, expires_at = $5, revoked_at = $6
		WHERE id = $1
	`

	var revokedAt interface{}
	if token.RevokedAt != nil {
		revokedAt = *token.RevokedAt
	}

	result, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.SessionID,
		token.TokenHash,
		token.ExpiresAt,
		revokedAt,
	)

	if err != nil {
		return errors.NewInternalError("Failed to update refresh token", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("Refresh token not found", nil)
	}

	return nil
}

// Delete deletes a refresh token
func (r *TokenRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewInternalError("Failed to delete refresh token", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("Refresh token not found", nil)
	}

	return nil
}

// Revoke revokes a refresh token
func (r *TokenRepository) Revoke(ctx context.Context, id string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to revoke refresh token", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("Refresh token not found or already revoked", nil)
	}

	return nil
}

// RevokeByUserID revokes all refresh tokens for a user
func (r *TokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE user_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, userID, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to revoke refresh tokens by user ID", err)
	}

	return nil
}

// RevokeBySessionID revokes all refresh tokens for a session
func (r *TokenRepository) RevokeBySessionID(ctx context.Context, sessionID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE session_id = $1 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, sessionID, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to revoke refresh tokens by session ID", err)
	}

	return nil
}

// RevokeByTokenHash revokes a refresh token by token hash with a reason
func (r *TokenRepository) RevokeByTokenHash(ctx context.Context, tokenHash string, reason string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	// For now, we ignore the reason parameter as our current schema doesn't store revocation reasons
	// In a production system, you might want to add a revocation_reason column
	result, err := r.db.ExecContext(ctx, query, tokenHash, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to revoke refresh token by hash", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("Refresh token not found or already revoked", nil)
	}

	return nil
}

// RevokeByTokenHashAndSessionID revokes a refresh token by token hash and session ID with a reason
func (r *TokenRepository) RevokeByTokenHashAndSessionID(ctx context.Context, tokenHash string, sessionID valueobjects.SessionID, reason string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE token_hash = $1 AND session_id = $3 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, tokenHash, time.Now(), sessionID.String())
	if err != nil {
		return errors.NewInternalError("Failed to revoke refresh token by hash and session ID", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("Refresh token not found or already revoked", nil)
	}

	return nil
}

// RevokeAllByUserID revokes all refresh tokens for a user with a reason
func (r *TokenRepository) RevokeAllByUserID(ctx context.Context, userID string, reason string) error {
	// For now, we ignore the reason parameter as our current schema doesn't store revocation reasons
	// In a production system, you might want to add a revocation_reason column
	return r.RevokeByUserID(ctx, userID)
}

// CleanupExpired removes expired refresh tokens
func (r *TokenRepository) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return errors.NewInternalError("Failed to cleanup expired refresh tokens", err)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens (alias for CleanupExpired)
func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	return r.CleanupExpired(ctx)
}
