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

func NewURLService(repo *storage.URLRepository, bloomFilter *filter.BloomFilter) *URLService {
	return &URLService{
		repo:        repo,
		bloomFilter: bloomFilter,
	}
}

func (s *URLService) ShortenURL(longURL string) (*models.URL, error) {
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

	url, err := s.repo.CreateURL(shortURL, longURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	s.bloomFilter.Add(longURL)

	return url, nil
}

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