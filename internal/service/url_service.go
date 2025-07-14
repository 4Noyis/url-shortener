package service

import (
	"fmt"

	"github.com/4Noyis/url-shortener/internal/encoding"
	"github.com/4Noyis/url-shortener/internal/filter"
	"github.com/4Noyis/url-shortener/internal/models"
	"github.com/4Noyis/url-shortener/internal/storage"
)

type URLService struct {
	repo        *storage.URLRepository
	bloomFilter *filter.BloomFilter
}

// NewURLService creates a new URLService instance with repository and bloom filter dependencies
func NewURLService(repo *storage.URLRepository, bloomFilter *filter.BloomFilter) *URLService {
	return &URLService{
		repo:        repo,
		bloomFilter: bloomFilter,
	}
}

// ShortenURL creates a shortened URL for the given long URL without TTL
func (s *URLService) ShortenURL(longURL string) (*models.URL, error) {
	return s.ShortenURLWithTTL(longURL, nil)
}

// ShortenURLWithTTL creates a shortened URL with optional time-to-live expiration
func (s *URLService) ShortenURLWithTTL(longURL string, ttlSeconds *int) (*models.URL, error) {
	if s.bloomFilter.Test(longURL) {
		exists, err := s.repo.URLExists(longURL)
		if err != nil {
			return nil, fmt.Errorf("failed to check if URL exists: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("URL already exists")
		}
	}

	lastID, err := s.repo.GetLastURLID()
	if err != nil {
		return nil, fmt.Errorf("failed to get last URL ID: %w", err)
	}

	nextID := s.repo.GenerateNextID(lastID)
	shortURL := encoding.EncodeIntToBase62(int64(nextID))

	url, err := s.repo.CreateURLWithTTL(shortURL, longURL, ttlSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	s.bloomFilter.Add(longURL)

	return url, nil
}

// RedirectURL retrieves the original long URL for a given short URL and increments click count
func (s *URLService) RedirectURL(shortURL string) (string, error) {
	url, err := s.repo.GetByShortURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}

	if err := s.repo.IncrementClicks(shortURL); err != nil {
		return "", fmt.Errorf("failed to increment clicks: %w", err)
	}

	return url.LongURL, nil
}

// CleanupExpiredURLs removes expired URLs from the database and returns the count of deleted records
func (s *URLService) CleanupExpiredURLs() (int, error) {
	count, err := s.repo.DeleteExpiredURLs()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired URLs: %w", err)
	}
	return count, nil
}