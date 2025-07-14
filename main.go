package main

import (
	"log"
	"time"

	"github.com/4Noyis/url-shortener/config"
	"github.com/4Noyis/url-shortener/internal/cleanup"
	"github.com/4Noyis/url-shortener/internal/filter"
	"github.com/4Noyis/url-shortener/internal/handlers"
	"github.com/4Noyis/url-shortener/internal/server"
	"github.com/4Noyis/url-shortener/internal/service"
	"github.com/4Noyis/url-shortener/internal/storage"
)

// main initializes the URL shortener service with database connection, bloom filter, and HTTP server
func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	db, err := config.GetDBConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	urlRepo := storage.NewURLRepository(db)

	bloomFilter := filter.NewBloomFilter(1000000, 0.01)

	existingURLs, err := urlRepo.GetAllLongURLs()
	if err != nil {
		log.Printf("Warning: Failed to load existing URLs into bloom filter: %v", err)
	} else {
		bloomFilter.AddAll(existingURLs)
		log.Printf("Loaded %d existing URLs into bloom filter", len(existingURLs))
	}

	urlService := service.NewURLService(urlRepo, bloomFilter)
	urlHandler := handlers.NewURLHandler(urlService)

	// Start cleanup scheduler (runs every hour)
	cleanupScheduler := cleanup.NewScheduler(urlService, time.Hour)
	go cleanupScheduler.Start()

	router := server.SetupRoutes(urlHandler)
	server.StartServer(router, "8080")
}
