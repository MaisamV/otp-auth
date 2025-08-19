package postgres

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/otp-auth/pkg/errors"
)

// Config holds PostgreSQL connection configuration
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns a default PostgreSQL configuration
func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		Database:        "otp_auth",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// NewConnection creates a new PostgreSQL database connection
func NewConnection(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.NewInternalError("Failed to open database connection", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, errors.NewInternalError("Failed to ping database", err)
	}

	return db, nil
}

// loadMigrationFiles loads migration files from the migrations directory
func loadMigrationFiles() ([]struct {
	version string
	query   string
}, error) {
	// Get the migrations directory path
	migrationsDir := "internal/infrastructure/persistence/postgres/migrations"
	
	// Read all files in the migrations directory
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)
	
	// Load migration content
	var migrations []struct {
		version string
		query   string
	}
	
	for _, filename := range sqlFiles {
		// Extract version from filename (remove .sql extension)
		version := strings.TrimSuffix(filename, ".sql")
		
		// Read file content
		filePath := filepath.Join(migrationsDir, filename)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}
		
		migrations = append(migrations, struct {
			version string
			query   string
		}{
			version: version,
			query:   string(content),
		})
	}
	
	return migrations, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`

	if _, err := db.Exec(createTableQuery); err != nil {
		return errors.NewInternalError("Failed to create migrations table", err)
	}

	// Load migration files from filesystem
	migrations, err := loadMigrationFiles()
	if err != nil {
		return errors.NewInternalError("Failed to load migration files", err)
	}

	// Run each migration
	for _, migration := range migrations {
		// Check if migration has already been applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", migration.version).Scan(&count)
		if err != nil {
			return errors.NewInternalError(fmt.Sprintf("Failed to check migration %s", migration.version), err)
		}

		if count > 0 {
			// Migration already applied, skip
			continue
		}

		// Apply migration
		if _, err := db.Exec(migration.query); err != nil {
			return errors.NewInternalError(fmt.Sprintf("Failed to apply migration %s", migration.version), err)
		}

		// Record migration as applied
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.version); err != nil {
			return errors.NewInternalError(fmt.Sprintf("Failed to record migration %s", migration.version), err)
		}
	}

	return nil
}