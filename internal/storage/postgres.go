package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/4Noyis/url-shortener/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

// NewURLRepository creates a new URLRepository instance with the provided database connection pool
func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{db: db}
}

// GetLastURLID retrieves the highest ID from the urls table for generating the next unique ID
func (r *URLRepository) GetLastURLID() (int, error) {
	var lastID int

	query := "SELECT id FROM urls ORDER BY id DESC LIMIT 1"
	err := r.db.QueryRow(context.Background(), query).Scan(&lastID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get last URL ID: %w", err)
	}

	return lastID, nil
}

// CreateURL inserts a new URL mapping into the database without expiration
func (r *URLRepository) CreateURL(shortURL, longURL string) (*models.URL, error) {
	insertQuery := `INSERT INTO urls (short_url, long_url) VALUES ($1, $2) RETURNING id, created_at`

	var insertedID int
	var createdAt time.Time

	err := r.db.QueryRow(context.Background(), insertQuery, shortURL, longURL).Scan(&insertedID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert URL: %w", err)
	}

	return &models.URL{
		ID:        insertedID,
		ShortURL:  shortURL,
		LongURL:   longURL,
		CreatedAt: createdAt,
		Clicks:    0,
		IsActive:  true,
	}, nil
}

// GenerateNextID generates the next unique ID using a multiplication strategy to avoid collisions
func (r *URLRepository) GenerateNextID(lastID int) int {
	if lastID > 1000 {
		return (lastID + 1) * 1111
	}
	return (lastID + 1) * 11111
}

// URLExists checks if a long URL already exists in the database
func (r *URLRepository) URLExists(longURL string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM urls WHERE long_url = $1)"
	var exists bool
	err := r.db.QueryRow(context.Background(), query, longURL).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if URL exists: %w", err)
	}
	return exists, nil
}

// GetAllLongURLs retrieves all long URLs from the database for bloom filter initialization
func (r *URLRepository) GetAllLongURLs() ([]string, error) {
	query := "SELECT long_url FROM urls"
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all long URLs: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var longURL string
		if err := rows.Scan(&longURL); err != nil {
			return nil, fmt.Errorf("failed to scan long URL: %w", err)
		}
		urls = append(urls, longURL)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return urls, nil
}

// GetByShortURL retrieves a URL record by its short URL and validates it's not expired
func (r *URLRepository) GetByShortURL(shortURL string) (*models.URL, error) {
	query := `SELECT id, short_url, long_url, created_at, expires_at, clicks, user_id, is_active 
			  FROM urls WHERE short_url = $1 AND is_active = true`

	var url models.URL
	var expiresAt *time.Time
	var userIDStr *string

	err := r.db.QueryRow(context.Background(), query, shortURL).Scan(
		&url.ID,
		&url.ShortURL,
		&url.LongURL,
		&url.CreatedAt,
		&expiresAt,
		&url.Clicks,
		&userIDStr,
		&url.IsActive,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("short URL not found")
		}
		return nil, fmt.Errorf("failed to get URL by short URL: %w", err)
	}

	url.ExpiresAt = expiresAt
	url.UserID = userIDStr

	if expiresAt != nil && time.Now().After(*expiresAt) {
		return nil, fmt.Errorf("short URL has expired")
	}

	return &url, nil
}

// IncrementClicks increments the click counter for a specific short URL
func (r *URLRepository) IncrementClicks(shortURL string) error {
	query := "UPDATE urls SET clicks = clicks + 1 WHERE short_url = $1"
	_, err := r.db.Exec(context.Background(), query, shortURL)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}
	return nil
}

// CreateURLWithTTL inserts a new URL mapping with optional time-to-live expiration
func (r *URLRepository) CreateURLWithTTL(shortURL, longURL string, ttlSeconds *int) (*models.URL, error) {
	var expiresAt *time.Time
	if ttlSeconds != nil {
		expTime := time.Now().Add(time.Duration(*ttlSeconds) * time.Second)
		expiresAt = &expTime
	}

	insertQuery := `INSERT INTO urls (short_url, long_url, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at`

	var insertedID int
	var createdAt time.Time

	err := r.db.QueryRow(context.Background(), insertQuery, shortURL, longURL, expiresAt).Scan(&insertedID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert URL: %w", err)
	}

	return &models.URL{
		ID:        insertedID,
		ShortURL:  shortURL,
		LongURL:   longURL,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Clicks:    0,
		IsActive:  true,
	}, nil
}

// DeleteExpiredURLs removes all expired URLs from the database and returns the count of deleted records
func (r *URLRepository) DeleteExpiredURLs() (int, error) {
	query := "DELETE FROM urls WHERE expires_at IS NOT NULL AND expires_at < NOW()"
	result, err := r.db.Exec(context.Background(), query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired URLs: %w", err)
	}
	
	rowsAffected := result.RowsAffected()
	return int(rowsAffected), nil
}

// GetExpiredURLsCount returns the number of expired URLs in the database
func (r *URLRepository) GetExpiredURLsCount() (int, error) {
	query := "SELECT COUNT(*) FROM urls WHERE expires_at IS NOT NULL AND expires_at < NOW()"
	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count expired URLs: %w", err)
	}
	return count, nil
}
